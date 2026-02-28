package httpadapter

import "net/http"

func NewServer() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/evaluate", EvaluateHandler)
    mux.HandleFunc("/compare", CompareHandler)
    mux.HandleFunc("/simulate", SimulateHandler)
    return mux
}
