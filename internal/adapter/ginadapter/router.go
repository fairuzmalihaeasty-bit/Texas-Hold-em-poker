package ginadapter

import (
	"net/http"

	"github.com/example/texasholdem/internal/adapter/poker"
	"github.com/example/texasholdem/internal/domain"
	"github.com/example/texasholdem/internal/usecase"
	"github.com/gin-gonic/gin"
)

// request/response models
type EvalReq struct {
	Hole      []string `json:"hole" binding:"required,len=2"`
	Community []string `json:"community"`
}

type CompareReq struct {
	Hand1 EvalReq `json:"hand1" binding:"required"`
	Hand2 EvalReq `json:"hand2" binding:"required"`
}

type SimReq struct {
	Hero        EvalReq  `json:"hero" binding:"required"`
	Community   []string `json:"community_known"`
	NumPlayers  int      `json:"num_players"`
	Iterations  int      `json:"iterations"`
	Concurrency int      `json:"concurrency"`
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/evaluate", evaluateHandler)
	r.POST("/compare", compareHandler)
	r.POST("/simulate", simulateHandler)
	r.GET("/health", func(c *gin.Context){ c.String(http.StatusOK, "OK") })
	// Serve static frontend files under /web to avoid wildcard conflicts with other routes
	r.Static("/web", "./web")
	// Serve index at site root
	r.GET("/", func(c *gin.Context) { c.File("./web/index.html") })
	// For any unknown route, return the index (SPA fallback)
	r.NoRoute(func(c *gin.Context) { c.File("./web/index.html") })
	return r
}

func parseCardSlice(in []string) ([]domain.Card, error) {
	out := make([]domain.Card, 0, len(in))
	for _, s := range in {
		c, err := poker.ParseCard(s)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func evaluateHandler(c *gin.Context) {
	var req EvalReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hole, err := parseCardSlice(req.Hole)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	community, err := parseCardSlice(req.Community)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := usecase.EvaluateHand(hole, community)
	c.JSON(http.StatusOK, gin.H{"rank": res})
}

func compareHandler(c *gin.Context) {
	var req CompareReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h1h, err := parseCardSlice(req.Hand1.Hole)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h1c, err := parseCardSlice(req.Hand1.Community)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h2h, err := parseCardSlice(req.Hand2.Hole)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h2c, err := parseCardSlice(req.Hand2.Community)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cr := usecase.CompareHandsDetailed(domain.Hand{Hole: h1h, Community: h1c}, domain.Hand{Hole: h2h, Community: h2c})
	winner := "tie"
	if cr.Winner == 1 {
		winner = "hand1"
	}
	if cr.Winner == -1 {
		winner = "hand2"
	}

	// convert best cards to strings for readability
	bestToStr := func(cards []domain.Card) []string {
		out := make([]string, 0, len(cards))
		for _, cc := range cards {
			if cc.Raw != "" {
				out = append(out, cc.Raw)
			} else {
				out = append(out, cc.String())
			}
		}
		return out
	}

	c.JSON(http.StatusOK, gin.H{
		"winner": winner,
		"hand1":  gin.H{"best": bestToStr(cr.Best1), "rank": cr.Rank1},
		"hand2":  gin.H{"best": bestToStr(cr.Best2), "rank": cr.Rank2},
	})
}

func simulateHandler(c *gin.Context) {
	var req SimReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hero, err := parseCardSlice(req.Hero.Hole)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	community, err := parseCardSlice(req.Community)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.NumPlayers < 2 {
		req.NumPlayers = 2
	}
	if req.Iterations <= 0 {
		req.Iterations = 10000
	}
	if req.Concurrency <= 0 {
		req.Concurrency = 4
	}
	res, err := usecase.RunMonteCarlo(hero, community, req.NumPlayers, req.Iterations, req.Concurrency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
