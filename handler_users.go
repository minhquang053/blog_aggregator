package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/minhquang053/blog_aggregator/internal/database"
)

type userResponse struct {
	Id        string
	Create_At time.Time
	Update_At time.Time
	Name      string
	Api_Key   string
}

func databaseUserToUserResponse(user database.User) userResponse {
	return userResponse{
		Id:        user.ID.String(),
		Create_At: user.CreatedAt,
		Update_At: user.UpdatedAt,
		Name:      user.Name,
		Api_Key:   user.ApiKey,
	}
}

// Post /v1/users
func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		log.Println("handleLogin: " + err.Error())
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, 200, databaseUserToUserResponse(user))
}

// Get v1/users
func (apiCfg *apiConfig) handlerUsersRead(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUserResponse(user))
}
