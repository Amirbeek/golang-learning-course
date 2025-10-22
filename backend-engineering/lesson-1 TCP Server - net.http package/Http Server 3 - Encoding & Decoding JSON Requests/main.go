package main

import (
	"log"
	"net/http"
)

func main() {
	apiSrv := &api{addr: ":8081"}

	mux := http.NewServeMux()

	// Go 1.22+ method-based patterns:
	mux.HandleFunc("GET /users", apiSrv.getUsersHandler)
	mux.HandleFunc("POST /user", apiSrv.createUsersHandler)

	//  < 1.22
	/*
		mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				apiSrv.getUsersHandler(w, r)
			case http.MethodPost:
				apiSrv.createUsersHandler(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})
	*/

	srv := &http.Server{
		Addr:    apiSrv.addr,
		Handler: mux,
	}

	log.Printf("HTTP server listening on %s", apiSrv.addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
