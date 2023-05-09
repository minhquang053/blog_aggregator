package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	type readinessResponse struct {
		Status string `json:"status"`
	}

	respondWithJSON(w, 200, readinessResponse{
		Status: "ok",
	})
}
