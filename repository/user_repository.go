package repository

import (
	"database/sql"
	"forum/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		"SELECT id, username, email, password, avatar_url, gender FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.AvatarURL, &user.Gender)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		"SELECT id, username, email, password, avatar_url, gender FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.AvatarURL, &user.Gender)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Create(username, email, hashedPassword string) (int, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (username, email, password, avatar_url, gender) VALUES (?, ?, ?, ?, ?)",
		username, email, hashedPassword, "", "",
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *UserRepo) EmailExists(email string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(1) FROM users WHERE email = ?", email).Scan(&count)
	return count > 0, err
}

func (r *UserRepo) UsernameExists(username string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(1) FROM users WHERE username = ?", username).Scan(&count)
	return count > 0, err
}

func (r *UserRepo) UpdateProfile(id int, username, email, avatarURL, gender string) error {
	_, err := r.db.Exec(
		"UPDATE users SET username = ?, email = ?, avatar_url = ?, gender = ? WHERE id = ?",
		username, email, avatarURL, gender, id,
	)
	return err
}

func (r *UserRepo) UpdatePassword(id int, hashedPassword string) error {
	_, err := r.db.Exec(
		"UPDATE users SET password = ? WHERE id = ?",
		hashedPassword, id,
	)
	return err
}
