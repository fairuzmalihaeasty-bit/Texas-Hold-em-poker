package domain

import (
    "sort"
)

// Category values increase with strength
const (
    HighCard = iota
    OnePair
    TwoPair
    Trips
    Straight
    Flush
    FullHouse
    Quads
    StraightFlush
)

type HandRank struct {
    Category int   `json:"category"`
    Ranks    []int `json:"ranks"` // tiebreaker ranks descending
}

// Evaluate takes the 7 (or fewer) cards and returns the best HandRank
func Evaluate(cards []Card) HandRank {
    // normalize ranks (2..14)
    if len(cards) == 0 {
        return HandRank{Category: HighCard, Ranks: []int{}}
    }

    rankCount := make(map[int]int)
    suits := make(map[Suit][]int)
    ranksSet := make(map[int]bool)
    for _, c := range cards {
        r := int(c.Rank)
        rankCount[r]++
        suits[c.Suit] = append(suits[c.Suit], r)
        ranksSet[r] = true
    }

    // build sorted unique ranks (desc)
    uniqueRanks := make([]int, 0, len(ranksSet))
    for r := range ranksSet {
        uniqueRanks = append(uniqueRanks, r)
    }
    sort.Sort(sort.Reverse(sort.IntSlice(uniqueRanks)))

    // check for straight and straight with ace-low
    straightTop := findStraightTop(uniqueRanks)

    // check flush
    flushSuit := Suit(-1)
    for s, rs := range suits {
        if len(rs) >= 5 {
            flushSuit = s
            break
        }
    }

    // check straight flush
    if flushSuit >= 0 {
        rs := suits[flushSuit]
        uniq := uniqueInts(rs)
        sort.Sort(sort.Reverse(sort.IntSlice(uniq)))
        sfTop := findStraightTop(uniq)
        if sfTop > 0 {
            return HandRank{Category: StraightFlush, Ranks: []int{sfTop}}
        }
    }

    // quads
    quadRank := 0
    for r, cnt := range rankCount {
        if cnt == 4 && r > quadRank {
            quadRank = r
        }
    }
    if quadRank > 0 {
        kicker := highestExcluding(uniqueRanks, []int{quadRank})
        return HandRank{Category: Quads, Ranks: []int{quadRank, kicker}}
    }

    // full house (trips + pair)
    trips := []int{}
    pairs := []int{}
    for r, cnt := range rankCount {
        if cnt >= 3 {
            trips = append(trips, r)
        } else if cnt == 2 {
            pairs = append(pairs, r)
        }
    }
    sort.Sort(sort.Reverse(sort.IntSlice(trips)))
    sort.Sort(sort.Reverse(sort.IntSlice(pairs)))
    if len(trips) > 0 {
        if len(trips) > 1 {
            // use top trips as full house (top trips + second trips as pair)
            return HandRank{Category: FullHouse, Ranks: []int{trips[0], trips[1]}}
        }
        if len(pairs) > 0 {
            return HandRank{Category: FullHouse, Ranks: []int{trips[0], pairs[0]}}
        }
    }

    // flush
    if flushSuit >= 0 {
        rs := suits[flushSuit]
        uniq := uniqueInts(rs)
        sort.Sort(sort.Reverse(sort.IntSlice(uniq)))
        top5 := takeTop(uniq, 5)
        return HandRank{Category: Flush, Ranks: top5}
    }

    // straight
    if straightTop > 0 {
        return HandRank{Category: Straight, Ranks: []int{straightTop}}
    }

    // trips
    if len(trips) > 0 {
        kicker := highestExcludingMultiple(uniqueRanks, []int{trips[0]}, 2)
        return HandRank{Category: Trips, Ranks: append([]int{trips[0]}, kicker...)}
    }

    // two pair
    if len(pairs) >= 2 {
        top2 := []int{pairs[0], pairs[1]}
        kicker := highestExcluding(uniqueRanks, top2)
        return HandRank{Category: TwoPair, Ranks: []int{top2[0], top2[1], kicker}}
    }

    // one pair
    if len(pairs) == 1 {
        kickers := highestExcludingMultiple(uniqueRanks, []int{pairs[0]}, 3)
        return HandRank{Category: OnePair, Ranks: append([]int{pairs[0]}, kickers...)}
    }

    // high card
    top5 := takeTop(uniqueRanks, 5)
    return HandRank{Category: HighCard, Ranks: top5}
}

func findStraightTop(sortedUnique []int) int {
    if len(sortedUnique) < 5 {
        return 0
    }
    // include ace-low possibility
    seen := map[int]bool{}
    for _, r := range sortedUnique {
        seen[r] = true
    }
    // add 1 if Ace exists
    if seen[14] && !seen[1] {
        sortedUnique = append(sortedUnique, 1)
        sort.Sort(sort.Reverse(sort.IntSlice(sortedUnique)))
    }
    // scan for 5-run
    for i := 0; i <= len(sortedUnique)-5; i++ {
        ok := true
        top := sortedUnique[i]
        for j := 1; j < 5; j++ {
            if !contains(sortedUnique, top-j) {
                ok = false
                break
            }
        }
        if ok {
            return top
        }
    }
    return 0
}

func contains(arr []int, v int) bool {
    for _, x := range arr {
        if x == v {
            return true
        }
    }
    return false
}

func uniqueInts(a []int) []int {
    m := map[int]bool{}
    out := []int{}
    for _, v := range a {
        if !m[v] {
            m[v] = true
            out = append(out, v)
        }
    }
    return out
}

func takeTop(a []int, n int) []int {
    if len(a) <= n {
        return append([]int{}, a...)
    }
    return append([]int{}, a[:n]...)
}

func highestExcluding(sortedDesc []int, exclude []int) int {
    ex := map[int]bool{}
    for _, e := range exclude {
        ex[e] = true
    }
    for _, r := range sortedDesc {
        if !ex[r] {
            return r
        }
    }
    return 0
}

func highestExcludingMultiple(sortedDesc []int, exclude []int, take int) []int {
    ex := map[int]bool{}
    for _, e := range exclude {
        ex[e] = true
    }
    out := []int{}
    for _, r := range sortedDesc {
        if !ex[r] {
            out = append(out, r)
            if len(out) == take {
                break
            }
        }
    }
    return out
}
