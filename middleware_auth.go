package main

import (
	"net/http"
	"strings"

	"github.com/adominguez793/Blog_Aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Client lacks API Key in Header")
			return
		}

		splitHeader := strings.Split(authHeader, " ")
		if len(splitHeader) < 2 {
			respondWithError(w, http.StatusUnauthorized, "Improper header. Access denied.")
			return
		}
		APIKey := splitHeader[1]

		user, err := cfg.DB.APIKeyGetUser(r.Context(), APIKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Failed to get user via api key")
			return
		}

		handler(w, r, user)
	}
}
