package usecase

import (
    "github.com/example/texasholdem/internal/domain"
)

type CompareResult struct {
    Best1  []domain.Card    `json:"best1"`
    Rank1  domain.HandRank  `json:"rank1"`
    Best2  []domain.Card    `json:"best2"`
    Rank2  domain.HandRank  `json:"rank2"`
    Winner int              `json:"winner"` // 1 = hand1, -1 = hand2, 0 = tie
}

// CompareHandsDetailed returns the best 5-card hand and HandRank for each player,
// plus the winner (1=hand1, -1=hand2, 0=tie).
func CompareHandsDetailed(h1, h2 domain.Hand) CompareResult {
    cards1 := append([]domain.Card{}, h1.Hole...)
    cards1 = append(cards1, h1.Community...)
    r1, best1 := domain.BestHand(cards1)

    cards2 := append([]domain.Card{}, h2.Hole...)
    cards2 = append(cards2, h2.Community...)
    r2, best2 := domain.BestHand(cards2)

    winner := 0
    if r1.Category > r2.Category {
        winner = 1
    } else if r1.Category < r2.Category {
        winner = -1
    } else {
        // same category: lexicographic compare
        for i := 0; i < len(r1.Ranks) && i < len(r2.Ranks); i++ {
            if r1.Ranks[i] > r2.Ranks[i] {
                winner = 1
                break
            }
            if r1.Ranks[i] < r2.Ranks[i] {
                winner = -1
                break
            }
        }
        if winner == 0 {
            if len(r1.Ranks) > len(r2.Ranks) {
                winner = 1
            } else if len(r1.Ranks) < len(r2.Ranks) {
                winner = -1
            } else {
                winner = 0
            }
        }
    }

    return CompareResult{Best1: best1, Rank1: r1, Best2: best2, Rank2: r2, Winner: winner}
}

// Keep backwards-compatible simple compare
func CompareHands(h1, h2 domain.Hand) int {
    return CompareHandsDetailed(h1, h2).Winner
}
