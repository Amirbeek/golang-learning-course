package main

import "net/http"

type api struct {
	addr string
}

func (a api) getUsersHandler(writer http.ResponseWriter, request *http.Request) {

}

func (a api) postUsersHandler(writer http.ResponseWriter, request *http.Request) {

}

func main() {
	api := &api{addr: ":8081"}
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    api.addr,
		Handler: mux,
	}

	mux.HandleFunc("GET /users", api.getUsersHandler)
	mux.HandleFunc("POST /users", api.postUsersHandler)

	srv.ListenAndServe()
}
