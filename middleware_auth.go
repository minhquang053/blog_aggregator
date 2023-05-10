package main

import (
	"net/http"

	"github.com/minhquang053/blog_aggregator/internal/auth"
	"github.com/minhquang053/blog_aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 401, "Couldn't find API key")
			return
		}
		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 404, "User not found")
			return
		}

		handler(w, r, user)
	}
}
