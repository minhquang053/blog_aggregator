package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/minhquang053/blog_aggregator/internal/database"
)

func (cfg *apiConfig) handlerPostsRead(w http.ResponseWriter, r *http.Request, user database.User) {
	limit := chi.URLParam(r, "limit")
	numPosts, err := strconv.Atoi(limit)
	if err != nil {
		numPosts = 10
	}

	defaultLimit := 10
	if numPosts <= 0 {
		numPosts = defaultLimit
	}

	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(numPosts),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch posts")
		return
	}

	respondWithJSON(w, 200, posts)
}
