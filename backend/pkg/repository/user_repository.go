package repository

import (
	"backend/pkg/model"
	"database/sql"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	return nil, nil
}

func (r *UserRepository) FindAll() (*model.User, error) {
	return nil, nil
}
