package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	err := http.ListenAndServe(srv.Addr, nil)
	if err != nil {
		panic(err)
	}
}
