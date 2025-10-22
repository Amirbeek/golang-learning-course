package main

import (
	"net/http"
)

type api struct {
	addr string
}

func (s *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Hello from the server"))
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/":
			w.Write([]byte("index page"))
			return
		case "/users":
			w.Write([]byte("users page"))
			return
		default:
			http.NotFound(w, r)
		}
	}
}

func main() {
	s := &api{addr: ":8081"}

	srv := &http.Server{Addr: s.addr, Handler: s}

	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}
