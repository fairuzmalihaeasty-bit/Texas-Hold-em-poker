package domain

import "testing"

func TestHandCategories(t *testing.T) {
    cases := []struct{
        name string
        cards []string
        want int
    }{
        {"StraightFlush", []string{"S5","S6","S7","S8","S9","H2","D3"}, StraightFlush},
        {"Quads", []string{"H9","D9","S9","C9","HA","D2","C3"}, Quads},
        {"FullHouse", []string{"H4","D4","C4","S9","D9","H2","C3"}, FullHouse},
        {"Flush", []string{"HA","HK","HQ","HJ","H9","D2","C3"}, Flush},
        {"StraightWheel", []string{"SA","H2","D3","C4","S5","D9","H7"}, Straight},
        {"Trips", []string{"HK","DK","SK","C2","D3","H4","S5"}, Trips},
        {"TwoPair", []string{"H5","D5","S9","C9","HA","D2","C3"}, TwoPair},
        {"OnePair", []string{"HA","DA","S9","C3","D4","H7","S8"}, OnePair},
        {"HighCard", []string{"HA","KD","C9","D3","S5","H7","C2"}, HighCard},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            cards := make([]Card, 0, len(tc.cards))
            for _, s := range tc.cards { cards = append(cards, mustParse(t, s)) }
            hr, best := BestHand(cards)
            if hr.Category != tc.want {
                t.Fatalf("%s: got category %d, want %d; best=%v", tc.name, hr.Category, tc.want, best)
            }
            if len(best) != 5 {
                t.Fatalf("%s: best hand len = %d, want 5", tc.name, len(best))
            }
        })
    }
}
