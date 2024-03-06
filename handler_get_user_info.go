package main

import (
	"net/http"
	"time"

	"github.com/adominguez793/Blog_Aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersGetInfo(w http.ResponseWriter, r *http.Request, user database.User) {
	// authHeader := r.Header.Get("Authorization")
	// if authHeader == "" {
	// 	respondWithError(w, http.StatusUnauthorized, "Client lacks API Key in Header")
	// 	return
	// }

	// splitHeader := strings.Split(authHeader, " ")
	// if len(splitHeader) < 2 {
	// 	respondWithError(w, http.StatusUnauthorized, "Improper header. Access denied.")
	// 	return
	// }
	// APIKey := splitHeader[1]

	// user, err := cfg.DB.APIKeyGetUser(r.Context(), APIKey)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "Failed to get user via api key")
	// 	return
	// }

	type returnVals struct {
		ID         uuid.UUID `json:"id"`
		Created_At time.Time `json:"created_at"`
		Updated_At time.Time `json:"updated_at"`
		Name       string    `json:"name"`
		APIKey     string    `json:"api_key"`
	}
	respondWithJSON(w, http.StatusOK, returnVals{
		ID:         user.ID,
		Created_At: user.CreatedAt,
		Updated_At: user.UpdatedAt,
		Name:       user.Name,
		APIKey:     user.ApiKey,
	})

}
