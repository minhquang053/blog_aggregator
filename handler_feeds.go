package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/minhquang053/blog_aggregator/internal/database"
)

type feedResponse struct {
	Id        string
	Create_At time.Time
	Update_At time.Time
	Name      string
	Url       string
	User_Id   string
}

func databaseFeedToFeedResponse(feed database.Feed) feedResponse {
	return feedResponse{
		Id:        feed.ID.String(),
		Create_At: feed.CreatedAt,
		Update_At: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		User_Id:   feed.UserID.String(),
	}
}

func (cfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		log.Println("handlerFeedsCreate: " + err.Error())
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Feed does not exist")
		return
	}

	respondWithJSON(w, 200, struct {
		Feed        feedResponse
		Feed_Follow feedFollowResponse
	}{
		databaseFeedToFeedResponse(feed),
		databaseFollowToFollowResponse(feedFollow),
	})
}

func (cfg *apiConfig) handlerFeedsRead(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch feeds")
		return
	}
	respondWithJSON(w, 200, feeds)
}
