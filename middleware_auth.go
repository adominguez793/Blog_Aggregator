package main

import (
	"net/http"

	"github.com/adominguez793/Blog_Aggregator/internal/auth"
	"github.com/adominguez793/Blog_Aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't find api key")
			return
		}

		user, err := cfg.DB.APIKeyGetUser(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Failed to get user via api key")
			return
		}

		handler(w, r, user)
	}
}
