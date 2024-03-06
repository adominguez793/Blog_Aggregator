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

	feedFollows, err := cfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed follows")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Feed       Feed
		FeedFollow FeedFollows
	}{
		databaseFeedToFeed(feed),
		databaseFeedFollowToFeedFollow(feedFollows),
	})
}

func (cfg *apiConfig) handlerFeedsGet(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve all feeds")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	}
}

func databaseFeedsToFeeds(feeds []database.Feed) []Feed {
	cleanFeeds := make([]Feed, len(feeds))
	for i, feed := range feeds {
		cleanFeeds[i] = databaseFeedToFeed(feed)
	}
	return cleanFeeds
}
