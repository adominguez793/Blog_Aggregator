package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/adominguez793/Blog_Aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type FeedFollows struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func (cfg *apiConfig) handlerFeedFollowsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode json")
		return
	}

	feedFollows, err := cfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		msg := fmt.Sprintf("Failed to create feed follows: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}
	respondWithJSON(w, http.StatusOK, FeedFollows{
		ID:        feedFollows.ID,
		CreatedAt: feedFollows.CreatedAt,
		UpdatedAt: feedFollows.UpdatedAt,
		UserID:    feedFollows.UserID,
		FeedID:    feedFollows.FeedID,
	})
}

func (cfg *apiConfig) handlerFeedFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowsIDString := chi.URLParam(r, "feedFollowID")
	feedFollowsID, err := uuid.Parse(feedFollowsIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to convert Feed Follows ID string into an UUID")
		return
	}

	err = cfg.DB.DeleteSpecificFeedFollows(r.Context(), database.DeleteSpecificFeedFollowsParams{
		ID:     feedFollowsID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong Feed ID or User ID")
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})

}

func (cfg *apiConfig) handlerFeedFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "False User ID")
		// NOTE: It's not guaranteed that a false user ID is the potential cause of an error here
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))
}

// These two functions below are extremely unclean. They look clean. They're not.
func databaseFeedFollowToFeedFollow(feedFollow database.FeedFollow) FeedFollows {
	return FeedFollows{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		UserID:    feedFollow.UserID,
		FeedID:    feedFollow.FeedID,
	}
}

func databaseFeedFollowsToFeedFollows(feedFollows []database.FeedFollow) []FeedFollows {
	cleanFeedFollows := make([]FeedFollows, len(feedFollows))
	for i, feedFollow := range feedFollows {
		cleanFeedFollows[i] = databaseFeedFollowToFeedFollow(feedFollow)
	}
	return cleanFeedFollows
}
