package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/adominguez793/Blog_Aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPostsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10

	specifiedLimit, err := strconv.Atoi(limitStr)
	if err == nil {
		limit = specifiedLimit
	}

	posts, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to get posts for user")
		return
	}

	respondWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePostsToPosts(posts []database.Post) []Post {
	cleanPosts := make([]Post, len(posts))

	for i, post := range posts {
		cleanPosts[i] = databasePostToPost(post)
	}
	return cleanPosts
}

func databasePostToPost(post database.Post) Post {
	return Post{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Title:       post.Title,
		Url:         post.Url,
		Description: post.Description,
		PublishedAt: post.PublishedAt,
		FeedID:      post.FeedID,
	}
}
