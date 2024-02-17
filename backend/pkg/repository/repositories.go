package repository

import (
	"database/sql"
)

// Repositories contains all the repo structs
type Repositories struct {
	UserRepo *UserRepository
}

// InitRepositories should be called in main.go
func InitRepositories(db *sql.DB) *Repositories {
	userRepo := NewUserRepository(db)
	return &Repositories{UserRepo: userRepo}
}
