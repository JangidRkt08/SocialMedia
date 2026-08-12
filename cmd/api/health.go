package main

import "net/http"

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("OK"))

}
