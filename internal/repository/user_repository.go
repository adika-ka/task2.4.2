package repository

import (
	"context"
	"database/sql"
	"repository/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) (int, error)
	GetByID(ctx context.Context, id int) (model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user model.User) (int, error) {
	query := `INSERT INTO users (email, password_hash, first_name, last_name, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.FirstName, user.LastName).Scan(&id)
	return id, err
}

func (r *userRepo) GetByID(ctx context.Context, id int) (model.User, error) {
	query := `SELECT id, email, password_hash, first_name, last_name, created_at, updated_at
	FROM users WHERE id = $1`

	var user model.User
	err := r.db.GetContext(ctx, &user, query, id)
	return user, err
}

func (r *userRepo) Update(ctx context.Context, user model.User) error {
	query := `UPDATE users SET email=$1, first_name=$2, last_name=$3, updated_at=NOW() WHERE id=$4`
	_, err := r.db.ExecContext(ctx, query, user.Email, user.FirstName, user.LastName, user.ID)
	return err
}

func (r *userRepo) Delete(ctx context.Context, id int) error {
	query := `UPDATE users SET deleted_at=NOW() WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepo) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	query := `SELECT id, email, first_name, last_name, created_at, updated_at, deleted_at FROM users
	ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	var users []model.User
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	return users, err
}

func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT 1 FROM users WHERE email = $1 LIMIT 1`

	var dummy int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
