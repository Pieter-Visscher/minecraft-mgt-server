package main

import (
	"fmt"
	"log"
	"net/http"
	//"context"
//	"flag"

	"minecraft-mgt-server/minecraft"
	"minecraft-mgt-server/k8s"
)

type Server struct {
	Manager *minecraft.Manager
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func serverCreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "server create handler")
}

func (s *Server) testpage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "webserver is running\n")
    err := s.Manager.CreateServer(r.Context(), "test-deploy")
    if err != nil {
        log.Printf("Error creating server: %v", err)
    }
}

func main() {
	kClient, err := k8s.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to Kubernetes: %v", err)
	}

	app := &Server{
		Manager: &minecraft.Manager{K8s: kClient},
	}

	log.Println("starting server on :8080....")
	mux := http.NewServeMux()

	mux.HandleFunc("POST /servers", serverCreateHandler)
	//mux.HandleFunc("GET /servers", serverListHandler)
	//mux.HandleFunc("GET /players", playersHandler)
	mux.HandleFunc("/test", app.testpage)

	wrappedMux := loggingMiddleware(mux)

	log.Fatal(http.ListenAndServe(":8080", wrappedMux))
}
