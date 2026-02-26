package postgres

import (
	"context"
	"log/slog"
	"messenger-auth/internal/models"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	tx     pgx.Tx
	logger *slog.Logger
}

func (r *UserRepository) Register(user *models.User) error {
	err := r.tx.QueryRow(context.Background(), "INSERT INTO users (username, password_hash) VALUES ($1, crypt($2, gen_salt('bf', 10))) RETURNING id", user.Username, user.Password).Scan(&user.Id)
	if err != nil {
		r.logger.Error("failed save user", "error", err)
	}
	return err
}

func (r *UserRepository) Login(user *models.User) error {
	err := r.tx.QueryRow(context.Background(), "SELECT id, created_at, updated_at FROM users WHERE username = $1 AND password_hash = crypt($2, password_hash)", user.Username, user.Password).Scan(&user.Id, &user.Created_at, &user.Updated_at)
	if err != nil {
		r.logger.Error("failed check user", "error", err)
	}
	return err
}
