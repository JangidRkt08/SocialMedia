package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "Server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad Request error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not Found error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "Not Found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusConflict, err.Error())
}
