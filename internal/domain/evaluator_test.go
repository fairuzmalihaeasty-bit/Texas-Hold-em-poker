package domain

import (
    "testing"
)

func mustParse(t *testing.T, s string) Card {
    t.Helper()
    c, err := ParseCard(s)
    if err != nil {
        t.Fatalf("parse %q: %v", s, err)
    }
    return c
}

func TestParseCardFormats(t *testing.T) {
    c := mustParse(t, "HA")
    if c.Rank != 14 || c.Suit != Hearts {
        t.Fatalf("expected HA -> Rank=14 Hearts, got %+v", c)
    }
    c2 := mustParse(t, "As")
    if c2.Rank != 14 {
        t.Fatalf("expected As -> rank 14, got %+v", c2)
    }
    c3 := mustParse(t, "C10")
    if c3.Rank != 10 || c3.Suit != Clubs {
        t.Fatalf("expected C10 -> 10 Clubs, got %+v", c3)
    }
}

func TestBestHand_RoyalFlush(t *testing.T) {
    cards := []Card{
        mustParse(t, "ST"), mustParse(t, "SJ"), mustParse(t, "SQ"), mustParse(t, "SK"), mustParse(t, "SA"),
        mustParse(t, "D2"), mustParse(t, "C3"),
    }
    hr, best := BestHand(cards)
    if hr.Category != StraightFlush {
        t.Fatalf("expected StraightFlush, got %d", hr.Category)
    }
    if len(best) != 5 {
        t.Fatalf("expected 5 cards, got %d", len(best))
    }
}

func TestBestHand_FullHouse(t *testing.T) {
    cards := []Card{
        mustParse(t, "H4"), mustParse(t, "D4"), mustParse(t, "C4"),
        mustParse(t, "S9"), mustParse(t, "D9"), mustParse(t, "S2"), mustParse(t, "H3"),
    }
    hr, _ := BestHand(cards)
    if hr.Category != FullHouse {
        t.Fatalf("expected FullHouse, got %d", hr.Category)
    }
}

func TestTieBreaking_PairKicker(t *testing.T) {
    // hand A: pair of Aces, kickers K,Q,J
    cardsA := []Card{
        mustParse(t, "HA"), mustParse(t, "DA"), mustParse(t, "SK"), mustParse(t, "SQ"), mustParse(t, "SJ"), mustParse(t, "C2"), mustParse(t, "D3"),
    }
    // hand B: pair of Aces, kickers K,Q,9 (worse than A)
    cardsB := []Card{
        mustParse(t, "HA"), mustParse(t, "DA"), mustParse(t, "SK"), mustParse(t, "SQ"), mustParse(t, "S9"), mustParse(t, "C2"), mustParse(t, "D3"),
    }
    hrA, _ := BestHand(cardsA)
    hrB, _ := BestHand(cardsB)
    if compareHandRank(hrA, hrB) <= 0 {
        t.Fatalf("expected hand A to beat hand B by kicker, got A=%+v B=%+v", hrA, hrB)
    }
}
