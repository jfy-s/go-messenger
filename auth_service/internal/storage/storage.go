package storage

import "messenger-auth/internal/models"

type Storage interface {
	CreateUnitOfWork() (UnitOfWork, error)
	Close()
}

type UnitOfWork interface {
	UserRepository() UserRepository
	Commit() error
	Rollback() error
}

type UserRepository interface {
	Register(User *models.User) error
	Login(User *models.User) error
}
