package service

import (
	"context"
	"fmt"
	"repository/internal/model"
	"repository/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GetUserByID(ctx context.Context, id int) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, limit, offset int) ([]model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user model.User) (int, error) {
	if user.Email == "" {
		return 0, fmt.Errorf("email cannot be empty")
	}
	ex, err := s.repo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return 0, err
	}
	if ex {
		return 0, fmt.Errorf("email already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("hash error: %v", err)
	}

	user.PasswordHash = string(hash)

	return s.repo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id int) (model.User, error) {
	if id <= 0 {
		return model.User{}, fmt.Errorf("invalid user id")
	}

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return u, nil
}

func (s *userService) UpdateUser(ctx context.Context, user model.User) error {
	if user.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid user id")
	}
	return s.repo.Delete(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, limit, offset int) ([]model.User, error) {
	if limit < 1 {
		limit = 10
	}

	return s.repo.List(ctx, limit, offset)
}
