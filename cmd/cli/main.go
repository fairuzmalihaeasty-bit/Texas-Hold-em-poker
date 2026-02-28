package main

import (
    "flag"
    "fmt"
    "log"
    "strings"

    "github.com/example/texasholdem/internal/usecase"
    "github.com/example/texasholdem/internal/adapter/poker"
)

func main() {
    hero := flag.String("hero", "As,Kd", "comma-separated hole cards")
    iter := flag.Int("iterations", 10000, "Monte Carlo iterations")
    flag.Parse()

    holeStr := strings.Split(*hero, ",")
    hole := make([]poker.Card, 0, len(holeStr))
    for _, s := range holeStr {
        c, err := poker.ParseCard(strings.TrimSpace(s))
        if err != nil {
            log.Fatalf("invalid card %q: %v", s, err)
        }
        hole = append(hole, c)
    }

    res := usecase.RunMonteCarloSingleOpponent(hole, nil, nil, *iter, 4)
    fmt.Printf("Result after %d iterations: win=%.4f tie=%.4f loss=%.4f\n", res.Iterations, res.WinPct, res.TiePct, res.LossPct)
}
