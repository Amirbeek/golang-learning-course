package main

import (
	"log"
	"net/http"
)

func (app *application) internalServeError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internalServeError: %s path %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusInternalServerError, "the server encountered an error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("badRequesetError: %s path %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("notFoundResponse: %s path %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusNotFound, "not found")
}
