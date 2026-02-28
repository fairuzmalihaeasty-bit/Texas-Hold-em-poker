package usecase

import (
    "testing"

    "github.com/example/texasholdem/internal/domain"
)

func mustParseCard(t *testing.T, s string) domain.Card {
    t.Helper()
    c, err := domain.ParseCard(s)
    if err != nil { t.Fatalf("parse %q: %v", s, err) }
    return c
}

func TestCompareHandsDetailed_WinnerAndBest(t *testing.T) {
    h1 := domain.Hand{Hole: []domain.Card{mustParseCard(t, "HA"), mustParseCard(t, "DA")}, Community: []domain.Card{mustParseCard(t, "2D"), mustParseCard(t, "7C"), mustParseCard(t, "JH"), mustParseCard(t, "TD"), mustParseCard(t, "3S")}}
    h2 := domain.Hand{Hole: []domain.Card{mustParseCard(t, "KS"), mustParseCard(t, "KH")}, Community: h1.Community}
    res := CompareHandsDetailed(h1, h2)
    if res.Winner != 1 {
        t.Fatalf("expected hand1 to win, got winner=%d", res.Winner)
    }
    if res.Rank1.Category == res.Rank2.Category && len(res.Rank1.Ranks) > 0 && len(res.Rank2.Ranks) > 0 && res.Rank1.Ranks[0] <= res.Rank2.Ranks[0] {
        t.Fatalf("unexpected rank ordering: %+v vs %+v", res.Rank1, res.Rank2)
    }
}

func TestCompareHandsDetailed_Tie(t *testing.T) {
    // identical hands -> tie
    h1 := domain.Hand{Hole: []domain.Card{mustParseCard(t, "HA"), mustParseCard(t, "DA")}, Community: []domain.Card{mustParseCard(t, "KD"), mustParseCard(t, "QD"), mustParseCard(t, "JD"), mustParseCard(t, "2S"), mustParseCard(t, "3C")}}
    h2 := h1
    res := CompareHandsDetailed(h1, h2)
    if res.Winner != 0 {
        t.Fatalf("expected tie winner=0, got %d", res.Winner)
    }
}
