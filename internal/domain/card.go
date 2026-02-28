package domain

import (
    "fmt"
    "strings"
    "errors"
)

type Suit int
type Rank int

const (
    Spades Suit = iota
    Hearts
    Diamonds
    Clubs
)

var suitMap = map[byte]Suit{
    's': Spades,
    'S': Spades,
    'h': Hearts,
    'H': Hearts,
    'd': Diamonds,
    'D': Diamonds,
    'c': Clubs,
    'C': Clubs,
}

var rankMap = map[byte]Rank{
    '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7,
    '8': 8, '9': 9, 'T': 10, 't': 10, 'J': 11, 'j': 11,
    'Q': 12, 'q': 12, 'K': 13, 'k': 13, 'A': 14, 'a': 14,
}

type Card struct {
    Rank Rank `json:"rank"`
    Suit Suit `json:"suit"`
    Raw  string `json:"raw"`
}

func ParseCard(s string) (Card, error) {
    s = strings.TrimSpace(s)
    if len(s) < 2 {
        return Card{}, errors.New("card string too short")
    }
    // Accept both formats: rank+suite (As, Td) and suit+rank (HA, S7, CT)
    // Detect if first char is a suit
    first := s[0]
    last := s[len(s)-1]

    // helper to parse rank token
    parseRankToken := func(token string) (Rank, bool) {
        if len(token) == 1 {
            if v, ok := rankMap[token[0]]; ok {
                return v, true
            }
            return 0, false
        }
        // support "10" or "T"
        if token == "10" {
            return 10, true
        }
        if v, ok := rankMap[token[0]]; ok {
            return v, true
        }
        return 0, false
    }

    // if first char is a suit, parse suit-first
    if _, ok := suitMap[first]; ok {
        suit := suitMap[first]
        rankToken := s[1:]
        if rankToken == "" {
            return Card{}, fmt.Errorf("invalid rank in %q", s)
        }
        rank, ok := parseRankToken(rankToken)
        if !ok {
            return Card{}, fmt.Errorf("invalid rank %q", rankToken)
        }
        return Card{Rank: rank, Suit: suit, Raw: s}, nil
    }

    // otherwise assume rank-first (e.g., "As" or "10h")
    rankToken := s[:len(s)-1]
    suitChar := last
    suit, ok := suitMap[suitChar]
    if !ok {
        return Card{}, fmt.Errorf("invalid suit %q", string(suitChar))
    }
    rank, ok := parseRankToken(rankToken)
    if !ok {
        return Card{}, fmt.Errorf("invalid rank %q", rankToken)
    }
    return Card{Rank: rank, Suit: suit, Raw: s}, nil
}

func (c Card) String() string {
    return c.Raw
}
