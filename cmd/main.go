package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "8080"

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}

func (cfg *apiConfig) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func GetHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	http.FileServer(http.Dir(".")).ServeHTTP(w, r)
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir(".")).ServeHTTP(w, r)
}

func GetHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}
func main() {
	mux := http.NewServeMux()
	cfg := apiConfig{}

	mux.Handle("GET /app/", cfg.middlewareMetricsInc(http.HandlerFunc(GetHtml)))
	mux.Handle("GET /app/assets/", cfg.middlewareMetricsInc(http.HandlerFunc(GetImage)))
	mux.Handle("GET /healthz", cfg.middlewareMetricsInc(http.HandlerFunc(GetHealthz)))
	mux.HandleFunc("/reset", cfg.Reset)
	mux.HandleFunc("/metrics", cfg.Metrics)

	// Создаем сервер с заданным портом и обработчиком
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
