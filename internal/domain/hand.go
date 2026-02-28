package domain

import "fmt"

type Hand struct {
    Hole      []Card `json:"hole"`
    Community []Card `json:"community"`
}

func (h Hand) All() []Card {
    out := make([]Card, 0, len(h.Hole)+len(h.Community))
    out = append(out, h.Hole...)
    out = append(out, h.Community...)
    return out
}

func (h Hand) Validate() error {
    if len(h.Hole) != 2 {
        return fmt.Errorf("hole must have 2 cards")
    }
    if len(h.Community) > 5 {
        return fmt.Errorf("community can be up to 5 cards")
    }
    return nil
}
