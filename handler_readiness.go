package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	type readiness struct {
		Status string `json:"status"`
	}

	respondWithJSON(w, 200, readiness{
		Status: "ok",
	})
}

func handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}
