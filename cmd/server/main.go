package main

import (
    "log"
    "os"

    "github.com/example/texasholdem/internal/adapter/ginadapter"
)

func main() {
    r := ginadapter.NewRouter()
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    addr := ":" + port
    log.Printf("starting Gin server on %s", addr)
    if err := r.Run(addr); err != nil {
        log.Fatal(err)
    }
}
