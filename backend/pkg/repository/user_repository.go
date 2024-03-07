package repository

import (
	"backend/pkg/model"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository with the given database connection.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByID finds a user by their ID and returns the user details. Returns nil if user is not found.
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	query := `SELECT id, nickname, email, password_hash, avatar_url FROM users WHERE id = $1;`
	row := r.db.QueryRow(query, id)

	err := row.Scan(&user.ID, &user.Nickname, &user.Email, &user.PasswordHash, &user.AvatarURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// FindByEmail finds a user by their email address and returns the user details. Returns nil if user is not found.
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, nickname, email, password_hash, avatar_url FROM users WHERE email = $1;`
	row := r.db.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Nickname, &user.Email, &user.PasswordHash, &user.AvatarURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// FindAll fetches all users from the database and returns a slice of user models.
func (r *UserRepository) FindAll() ([]model.User, error) {
	var users []model.User
	rows, err := r.db.Query("SELECT id, nickname, email, password_hash, avatar_url FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Nickname, &user.Email, &user.AvatarURL)
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

// AddNewUser adds a new user to the database and returns the newly created user.
func (r *UserRepository) AddNewUser(user model.User) (model.User, error) {
	query := `
		INSERT INTO users (nickname, email, password_hash, avatar_url) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, nickname, email, password_hash, avatar_url, created_at;
	`
	var newUser model.User
	err := r.db.QueryRow(query, user.Nickname, user.Email, user.PasswordHash, user.AvatarURL).
		Scan(&newUser.ID, &newUser.Nickname, &newUser.Email, &newUser.PasswordHash, &newUser.AvatarURL, &newUser.CreatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to add a new user: %v", err)
	}

	return newUser, nil
}
