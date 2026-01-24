package handlers

import (
	"log/slog"
	"messenger-auth/internal/jwt"
	"net/http"
)

func UpdateToken(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("update token request received")
		token := r.Header.Get("Authorization")
		if token == "" {
			logger.Error("no authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token = token[len("Bearer "):]
		id, err := jwt.ParseToken(token)
		if err != nil {
			logger.Error("failed to parse token", "error", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		newToken, err := jwt.CreateToken(id)
		if err != nil {
			logger.Error("failed to create new token", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Authorization", "Bearer "+newToken)
		w.WriteHeader(http.StatusCreated)
		logger.Info("token created", "user_id", id)
	}
}
