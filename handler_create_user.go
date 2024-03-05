package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adominguez793/Blog_Aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode parameters")
		return
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate new UUID")
		return
	}

	// encode(sha256(random()::text::bytea), 'hex')

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        userID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	type returnVals struct {
		ID         uuid.UUID `json:"id"`
		Created_At time.Time `json:"created_at"`
		Updated_At time.Time `json:"updated_at"`
		Name       string    `json:"name"`
	}
	respondWithJSON(w, http.StatusOK, returnVals{
		ID:         user.ID,
		Created_At: user.CreatedAt,
		Updated_At: user.UpdatedAt,
		Name:       user.Name,
	})
}
