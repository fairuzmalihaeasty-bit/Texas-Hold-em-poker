package usecase

import (
    "math/rand"
    "sync"
    "time"
    "fmt"

    "github.com/example/texasholdem/internal/domain"
)

type SimResult struct {
    Wins       int     `json:"wins"`
    Ties       int     `json:"ties"`
    Losses     int     `json:"losses"`
    Iterations int     `json:"iterations"`
    WinPct     float64 `json:"winPct"`
    TiePct     float64 `json:"tiePct"`
    LossPct    float64 `json:"lossPct"`
}

// RunMonteCarlo runs Monte Carlo simulations for hero vs (numPlayers-1) opponents.
// hero must contain 2 cards. communityKnown may contain 0,3,4,or 5 cards.
// iterations is number of Monte Carlo iterations; concurrency sets worker goroutines.
func RunMonteCarlo(hero []domain.Card, communityKnown []domain.Card, numPlayers int, iterations int, concurrency int) (SimResult, error) {
    if len(hero) != 2 {
        return SimResult{}, fmt.Errorf("hero must have exactly 2 hole cards")
    }
    if numPlayers < 2 {
        return SimResult{}, fmt.Errorf("numPlayers must be >= 2")
    }
    if iterations <= 0 {
        iterations = 10000
    }
    if concurrency <= 0 {
        concurrency = 4
    }

    deck := fullDeck()
    // remove known hero and known community
    known := append([]domain.Card{}, hero...)
    known = append(known, communityKnown...)
    removeCards(&deck, known...)

    jobs := make(chan int, iterations)
    var mu sync.Mutex
    wins := 0
    ties := 0
    losses := 0

    compareRanks := func(a, b domain.HandRank) int {
        if a.Category > b.Category { return 1 }
        if a.Category < b.Category { return -1 }
        for i := 0; i < len(a.Ranks) && i < len(b.Ranks); i++ {
            if a.Ranks[i] > b.Ranks[i] { return 1 }
            if a.Ranks[i] < b.Ranks[i] { return -1 }
        }
        if len(a.Ranks) > len(b.Ranks) { return 1 }
        if len(a.Ranks) < len(b.Ranks) { return -1 }
        return 0
    }

    worker := func() {
        r := rand.New(rand.NewSource(time.Now().UnixNano()))
        for range jobs {
            // copy and shuffle deck
            d := make([]domain.Card, len(deck))
            copy(d, deck)
            r.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
            idx := 0

            // deal opponents
            opponentHoles := make([][]domain.Card, numPlayers-1)
            for p := 0; p < numPlayers-1; p++ {
                opponentHoles[p] = []domain.Card{d[idx], d[idx+1]}
                idx += 2
            }

            // build full community
            community := make([]domain.Card, len(communityKnown))
            copy(community, communityKnown)
            for len(community) < 5 {
                community = append(community, d[idx])
                idx++
            }

            // evaluate hero
            hCards := append([]domain.Card{}, hero...)
            hCards = append(hCards, community...)
            hRank, _ := domain.BestHand(hCards)

            // evaluate opponents
            topComp := hRank
            heroIsTop := true
            tiedCount := 1
            for _, oh := range opponentHoles {
                oc := append([]domain.Card{}, oh...)
                oc = append(oc, community...)
                orank, _ := domain.BestHand(oc)
                cmp := compareRanks(orank, topComp)
                if cmp > 0 {
                    topComp = orank
                    heroIsTop = false
                    tiedCount = 0
                }
                if cmp == 0 {
                    // if opponent equals current top, increment tiedCount appropriately
                    if heroIsTop {
                        tiedCount++
                    } else {
                        // tied among opponents; hero not top
                    }
                }
                if cmp > 0 {
                    // new top is opponent, reset tied count
                    // if the new top equals hero? handled above
                }
            }

            // Determine result for hero
            // Re-evaluate hero comparisons against topComp
            finalCmp := compareRanks(hRank, topComp)
            mu.Lock()
            if finalCmp > 0 {
                wins++
            } else if finalCmp < 0 {
                losses++
            } else {
                // tie: count how many players tie for top including hero
                // to compute tie share we consider it a tie
                ties++
            }
            mu.Unlock()
        }
    }

    var wg sync.WaitGroup
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() { defer wg.Done(); worker() }()
    }

    for i := 0; i < iterations; i++ {
        jobs <- i
    }
    close(jobs)
    wg.Wait()

    total := iterations
    res := SimResult{
        Wins:       wins,
        Ties:       ties,
        Losses:     losses,
        Iterations: total,
    }
    if total > 0 {
        res.WinPct = float64(wins) / float64(total)
        res.TiePct = float64(ties) / float64(total)
        res.LossPct = float64(losses) / float64(total)
    }
    return res, nil
}

// Backwards-compatible single-opponent API. If villain has 2 cards provided,
// they are used every iteration; otherwise villain is randomized.
func RunMonteCarloSingleOpponent(hero []domain.Card, villain []domain.Card, communityKnown []domain.Card, iterations int, concurrency int) SimResult {
    // If villain provided, we'll run a specialized loop where villain is fixed.
    if len(villain) == 2 {
        // build deck and remove hero, villain, community
        deck := fullDeck()
        known := append([]domain.Card{}, hero...)
        known = append(known, villain...)
        known = append(known, communityKnown...)
        removeCards(&deck, known...)

        jobs := make(chan int, iterations)
        var mu sync.Mutex
        wins := 0
        ties := 0
        losses := 0

        compareRanks := func(a, b domain.HandRank) int {
            if a.Category > b.Category { return 1 }
            if a.Category < b.Category { return -1 }
            for i := 0; i < len(a.Ranks) && i < len(b.Ranks); i++ {
                if a.Ranks[i] > b.Ranks[i] { return 1 }
                if a.Ranks[i] < b.Ranks[i] { return -1 }
            }
            if len(a.Ranks) > len(b.Ranks) { return 1 }
            if len(a.Ranks) < len(b.Ranks) { return -1 }
            return 0
        }

        worker := func() {
            r := rand.New(rand.NewSource(time.Now().UnixNano()))
            for range jobs {
                d := make([]domain.Card, len(deck))
                copy(d, deck)
                r.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
                idx := 0

                community := make([]domain.Card, len(communityKnown))
                copy(community, communityKnown)
                for len(community) < 5 {
                    community = append(community, d[idx])
                    idx++
                }

                hCards := append([]domain.Card{}, hero...)
                hCards = append(hCards, community...)
                hRank, _ := domain.BestHand(hCards)

                vCards := append([]domain.Card{}, villain...)
                vCards = append(vCards, community...)
                vRank, _ := domain.BestHand(vCards)

                cmp := compareRanks(hRank, vRank)
                mu.Lock()
                if cmp > 0 { wins++ } else if cmp < 0 { losses++ } else { ties++ }
                mu.Unlock()
            }
        }

        var wg sync.WaitGroup
        if concurrency <= 0 { concurrency = 4 }
        for i := 0; i < concurrency; i++ { wg.Add(1); go func() { defer wg.Done(); worker() }() }
        for i := 0; i < iterations; i++ { jobs <- i }
        close(jobs)
        wg.Wait()

        total := iterations
        res := SimResult{Wins: wins, Ties: ties, Losses: losses, Iterations: total}
        if total > 0 {
            res.WinPct = float64(wins) / float64(total)
            res.TiePct = float64(ties) / float64(total)
            res.LossPct = float64(losses) / float64(total)
        }
        return res
    }

    // otherwise random villain: call general RunMonteCarlo with numPlayers=2
    r, _ := RunMonteCarlo(hero, communityKnown, 2, iterations, concurrency)
    return r
}

func fullDeck() []domain.Card {
    out := make([]domain.Card, 0, 52)
    ranks := []domain.Rank{2,3,4,5,6,7,8,9,10,11,12,13,14}
    for _, r := range ranks {
        for s := domain.Spades; s <= domain.Clubs; s++ {
            out = append(out, domain.Card{Rank: r, Suit: s})
        }
    }
    return out
}

func cardKey(c domain.Card) string {
    return fmt.Sprintf("%d:%d", int(c.Rank), int(c.Suit))
}

func removeCards(deck *[]domain.Card, remove ...domain.Card) {
    m := map[string]bool{}
    for _, c := range remove {
        m[cardKey(c)] = true
    }
    out := (*deck)[:0]
    for _, c := range *deck {
        if !m[cardKey(c)] {
            out = append(out, c)
        }
    }
    *deck = out
}
