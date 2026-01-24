package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"messenger-auth/internal/jwt"
	"messenger-auth/internal/models"
	"messenger-auth/internal/storage"
	"net/http"

	"github.com/go-playground/validator"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(logger *slog.Logger, storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("login request received")
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("failed to decode login request", "error", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		u := &models.User{
			Username: req.Username,
			Password: req.Password,
		}
		validate := validator.New()
		if err := validate.Struct(u); err != nil {
			logger.Error("validation failed for login request", "error", err)
			http.Error(w, "Invalid username or password", http.StatusBadRequest)
			return
		}
		// authenticate user
		uow, err := storage.CreateUnitOfWork()
		if err != nil {
			logger.Error("failed to create unit of work", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer uow.Rollback()
		if err = uow.UserRepository().Login(u); err != nil {
			logger.Error("failed to login user", "error", err)
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		// create token
		token, err := jwt.CreateToken(u.Id)
		if err != nil {
			logger.Error("failed to create token", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write([]byte(fmt.Sprintf("{\"token\": \"%s\"}", token))); err != nil {
			logger.Error("failed to write token to response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		logger.Info("user logged in successfully", "user_id", u.Id)
	}
}
