package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/adominguez793/Blog_Aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to decode parameters")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		msg := fmt.Sprintf("Failed to create feed: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserID    uuid.UUID `json:"user_id"`
	}
	respondWithJSON(w, http.StatusOK, returnVals{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	})
}
