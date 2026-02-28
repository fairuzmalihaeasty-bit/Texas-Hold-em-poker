package usecase

import (
    "github.com/example/texasholdem/internal/domain"
)

type EvalResult struct {
    Category int   `json:"category"`
    Ranks    []int `json:"ranks"`
}

func EvaluateHand(hole []domain.Card, community []domain.Card) EvalResult {
    cards := make([]domain.Card, 0, len(hole)+len(community))
    cards = append(cards, hole...)
    cards = append(cards, community...)
    hr := domain.Evaluate(cards)
    return EvalResult{Category: hr.Category, Ranks: hr.Ranks}
}
