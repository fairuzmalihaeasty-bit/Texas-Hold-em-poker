package domain

// BestHand brute-forces all 5-card combinations out of given cards (up to 7)
// and returns the best HandRank and the corresponding 5 cards.
func BestHand(cards []Card) (HandRank, []Card) {
    n := len(cards)
    if n < 5 {
        // evaluate all available cards
        hr := Evaluate(cards)
        return hr, append([]Card{}, cards...)
    }

    bestRank := HandRank{Category: -1}
    var bestCombo []Card

    // generate all 5-combinations (n choose 5)
    idx := []int{0,1,2,3,4}
    limit := n
    for {
        combo := []Card{cards[idx[0]], cards[idx[1]], cards[idx[2]], cards[idx[3]], cards[idx[4]]}
        r := Evaluate(combo)
        if compareHandRank(r, bestRank) > 0 {
            bestRank = r
            bestCombo = append([]Card{}, combo...)
        }

        // increment indices
        i := 4
        for i >= 0 {
            idx[i]++
            if idx[i] < limit-(4-i) {
                for j := i+1; j < 5; j++ {
                    idx[j] = idx[j-1] + 1
                }
                break
            }
            i--
        }
        if i < 0 {
            break
        }
    }
    return bestRank, bestCombo
}

// compareHandRank returns 1 if a>b, 0 if equal, -1 if a<b
func compareHandRank(a, b HandRank) int {
    if b.Category == -1 && a.Category >= 0 {
        return 1
    }
    if a.Category < b.Category {
        return -1
    }
    if a.Category > b.Category {
        return 1
    }
    // same category: lexicographic compare ranks
    for i := 0; i < len(a.Ranks) && i < len(b.Ranks); i++ {
        if a.Ranks[i] < b.Ranks[i] {
            return -1
        }
        if a.Ranks[i] > b.Ranks[i] {
            return 1
        }
    }
    if len(a.Ranks) < len(b.Ranks) {
        return -1
    }
    if len(a.Ranks) > len(b.Ranks) {
        return 1
    }
    return 0
}
