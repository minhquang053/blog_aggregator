package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/minhquang053/blog_aggregator/internal/database"
)

type feedFollowResponse struct {
	Id        uuid.UUID
	Feed_Id   uuid.UUID
	User_Id   uuid.UUID
	Create_At time.Time
	Update_At time.Time
}

func databaseFollowToFollowResponse(feedFollow database.FeedFollow) feedFollowResponse {
	return feedFollowResponse{
		Id:        feedFollow.ID,
		Feed_Id:   feedFollow.FeedID,
		User_Id:   feedFollow.UserID,
		Create_At: feedFollow.CreatedAt,
		Update_At: feedFollow.UpdatedAt,
	}
}

func (cfg *apiConfig) handlerFeedFollowsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Feed_Id uuid.UUID
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		log.Println("handlerFeedsCreate: " + err.Error())
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    params.Feed_Id,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Feed does not exist")
		return
	}

	respondWithJSON(w, 200, databaseFollowToFollowResponse(feedFollow))
}

func (cfg *apiConfig) handlerFeedFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	ffid := chi.URLParam(r, "feedFollowID")
	id, err := uuid.Parse(ffid)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     id,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete feed_follow")
		return
	}

	respondWithJSON(w, 200, struct{}{})
}

func (cfg *apiConfig) handlerFeedFollowsRead(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch feeds for user")
		return
	}
	respondWithJSON(w, 200, feedFollows)
}
