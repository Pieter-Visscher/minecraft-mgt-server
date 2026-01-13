package main

import (
	"fmt"
	"log"
	"net/http"

)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func serverCreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "server create handler")
}

func testpage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "webserver is running")
}

func main() {
	mux := http.NewServeMux()

	log.Println("starting server on :8080....")

	mux.HandleFunc("POST /servers", serverCreateHandler)
	//mux.HandleFunc("GET /servers", serverListHandler)
	//mux.HandleFunc("GET /players", playersHandler)
	mux.HandleFunc("/test", testpage)

	wrappedMux := loggingMiddleware(mux)

	log.Fatal(http.ListenAndServe(":8080", wrappedMux))
}
