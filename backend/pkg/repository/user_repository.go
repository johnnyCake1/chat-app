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

func (r *UserRepository) FindAll() ([]model.User, error) {
	var users []model.User
	rows, err := r.db.Query("SELECT * FROM users;") // Ensure your table name is correct
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.AvatarURL)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
