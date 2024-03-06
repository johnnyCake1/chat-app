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

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	query := `SELECT id, nickname, email, password_hash, avatar_url FROM users WHERE id = $1;`
	row := r.db.QueryRow(query, id)

	err := row.Scan(&user.ID, &user.Nickname, &user.Email, &user.PasswordHash, &user.AvatarURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

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

func (r *UserRepository) FindAll() ([]model.User, error) {
	var users []model.User
	rows, err := r.db.Query("SELECT * FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Nickname, &user.Email, &user.PasswordHash, &user.AvatarURL)
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

func (r *UserRepository) AddNewUser(user model.User) (uint, error) {
	query := `
		INSERT INTO users (nickname, email, password_hash, avatar_url) 
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	var userID uint
	err := r.db.QueryRow(query, user.Nickname, user.Email, user.PasswordHash, user.AvatarURL).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to add a new user: %v", err)
	}

	return userID, nil
}
