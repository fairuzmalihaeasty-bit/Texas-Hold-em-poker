package usecase

import (
    "testing"

    "github.com/example/texasholdem/internal/domain"
)

// Deterministic test: provide full 5-card community and both hero+villain holes.
// There is no randomness in the outcome, so Monte Carlo should return deterministic wins/losses.
func TestRunMonteCarloSingleOpponent_DeterministicFullBoard(t *testing.T) {
    hero := []domain.Card{mustParseCard(t, "HA"), mustParseCard(t, "DA")}
    villain := []domain.Card{mustParseCard(t, "HK"), mustParseCard(t, "DK")}
    community := []domain.Card{mustParseCard(t, "KD"), mustParseCard(t, "QD"), mustParseCard(t, "JD"), mustParseCard(t, "2S"), mustParseCard(t, "3C")}

    // compute expected outcome by direct compare
    resDirect := CompareHandsDetailed(domain.Hand{Hole: hero, Community: community}, domain.Hand{Hole: villain, Community: community})

    // run Monte Carlo with provided villain and full community
    iterations := 20
    sim := RunMonteCarloSingleOpponent(hero, villain, community, iterations, 2)

    // ensure sum matches and deterministic result aligns
    if sim.Wins+sim.Ties+sim.Losses != iterations {
        t.Fatalf("sum mismatch: wins+ties+losses=%d, iterations=%d", sim.Wins+sim.Ties+sim.Losses, iterations)
    }

    // expected winner
    if resDirect.Winner == 1 {
        if sim.Wins != iterations { t.Fatalf("expected all wins, got %+v", sim) }
    } else if resDirect.Winner == -1 {
        if sim.Losses != iterations { t.Fatalf("expected all losses, got %+v", sim) }
    } else {
        if sim.Ties != iterations { t.Fatalf("expected all ties, got %+v", sim) }
    }
}

func TestRunMonteCarlo_BasicProperties(t *testing.T) {
    // basic property test: outputs sum to iterations and percentages are within [0,1]
    hero := []domain.Card{mustParseCard(t, "AS"), mustParseCard(t, "KD")}
    iterations := 50
    res, err := RunMonteCarlo(hero, nil, 2, iterations, 2)
    if err != nil { t.Fatalf("RunMonteCarlo error: %v", err) }
    if res.Wins+res.Ties+res.Losses != iterations { t.Fatalf("sum mismatch: %+v", res) }
    if res.WinPct < 0 || res.WinPct > 1 { t.Fatalf("winPct out of range: %v", res.WinPct) }
    if res.TiePct < 0 || res.TiePct > 1 { t.Fatalf("tiePct out of range: %v", res.TiePct) }
    if res.LossPct < 0 || res.LossPct > 1 { t.Fatalf("lossPct out of range: %v", res.LossPct) }
}
