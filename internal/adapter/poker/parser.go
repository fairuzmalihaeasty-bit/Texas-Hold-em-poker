package poker

import (
    "github.com/example/texasholdem/internal/domain"
    "strings"
    "fmt"
)

type Card = domain.Card

func ParseCard(s string) (Card, error) {
    s = strings.TrimSpace(s)
    // support formats like "As" or "Ah"
    if len(s) < 2 {
        return Card{}, fmt.Errorf("invalid card %q", s)
    }
    return domain.ParseCard(s)
}
