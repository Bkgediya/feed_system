package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Bkgediya/feed_system/auth-service/internal/model"
)

type UserRepository interface {
	Create(u *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByID(id int64) (*model.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(u *model.User) error {
	q := `INSERT INTO users (username, email, password_hash, created_at, updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`
	err := r.db.QueryRow(q, u.Username, u.Email, u.PasswordHash, time.Now(), time.Now()).Scan(&u.ID)
	return err
}

func (r *userRepo) GetByEmail(email string) (*model.User, error) {
	q := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	var u model.User
	err := r.db.QueryRow(q, email).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByID(id int64) (*model.User, error) {
	q := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
	var u model.User
	err := r.db.QueryRow(q, id).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
