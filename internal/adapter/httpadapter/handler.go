package httpadapter

import (
    "encoding/json"
    "net/http"

    "github.com/example/texasholdem/internal/adapter/poker"
    "github.com/example/texasholdem/internal/usecase"
    "github.com/example/texasholdem/internal/domain"
)

type EvalReq struct {
    Hole      []string `json:"hole"`
    Community []string `json:"community"`
}

type CompareReq struct {
    Hand1 EvalReq `json:"hand1"`
    Hand2 EvalReq `json:"hand2"`
}

type SimReq struct {
    Hero          EvalReq `json:"hero"`
    Villain       *EvalReq `json:"villain"`
    CommunityKnown []string `json:"community_known"`
    Iterations    int     `json:"iterations"`
    Concurrency   int     `json:"concurrency"`
}

func parseReqToCards(arr []string) ([]domain.Card, error) {
    out := make([]domain.Card, 0, len(arr))
    for _, s := range arr {
        c, err := poker.ParseCard(s)
        if err != nil {
            return nil, err
        }
        out = append(out, c)
    }
    return out, nil
}

func EvaluateHandler(w http.ResponseWriter, r *http.Request) {
    var req EvalReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    hole, err := parseReqToCards(req.Hole)
    if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    community, err := parseReqToCards(req.Community)
    if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    res := usecase.EvaluateHand(hole, community)
    json.NewEncoder(w).Encode(res)
}

func CompareHandler(w http.ResponseWriter, r *http.Request) {
    var req CompareReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    h1h, err := parseReqToCards(req.Hand1.Hole); if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    h1c, err := parseReqToCards(req.Hand1.Community); if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    h2h, err := parseReqToCards(req.Hand2.Hole); if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    h2c, err := parseReqToCards(req.Hand2.Community); if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    cr := usecase.CompareHandsDetailed(domain.Hand{Hole: h1h, Community: h1c}, domain.Hand{Hole: h2h, Community: h2c})
    winnerStr := "tie"
    if cr.Winner == 1 { winnerStr = "hand1" }
    if cr.Winner == -1 { winnerStr = "hand2" }

    out := map[string]interface{}{
        "winner": winnerStr,
        "hand1": map[string]interface{}{"best": cr.Best1, "rank": cr.Rank1},
        "hand2": map[string]interface{}{"best": cr.Best2, "rank": cr.Rank2},
    }
    json.NewEncoder(w).Encode(out)
}

func SimulateHandler(w http.ResponseWriter, r *http.Request) {
    var req SimReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    heroHole, err := parseReqToCards(req.Hero.Hole); if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    var villainHole []domain.Card
    if req.Villain != nil {
        villainHole, err = parseReqToCards(req.Villain.Hole); if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    }
    communityKnown, err := parseReqToCards(req.CommunityKnown); if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    if req.Iterations <= 0 { req.Iterations = 10000 }
    if req.Concurrency <= 0 { req.Concurrency = 4 }
    res := usecase.RunMonteCarloSingleOpponent(heroHole, villainHole, communityKnown, req.Iterations, req.Concurrency)
    json.NewEncoder(w).Encode(res)
}
