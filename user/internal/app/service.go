package app

import (
	"context"
	"time"
	"user/internal/domain"

	"github.com/google/uuid"
)

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type UserService struct {
	userRepo domain.UserRepository
}

type CreateUserParams struct {
	Name  string
	Email string
	Level *domain.DancerLevel
}

func (us *UserService) CreateUser(ctx context.Context, params CreateUserParams) (*domain.User, error) {
	uuid := uuid.NewString()
	user := domain.User{
		ID:        uuid,
		Name:      params.Name,
		Email:     params.Email,
		Level:     params.Level,
		CreatedAt: time.Now(),
	}
	err := us.userRepo.Store(ctx, user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) FindUser(ctx context.Context, userID string) (*domain.User, error) {
	return us.userRepo.Find(ctx, userID)
}

func (us *UserService) ListUsres(ctx context.Context) ([]*domain.User, error) {
	return us.userRepo.ListUsers(ctx)
}
