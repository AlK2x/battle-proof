package domain

import (
	"context"
	"time"
)

type DancerLevel string

const (
	LevelBeginner DancerLevel = "beginner"
	LevelAmateur  DancerLevel = "amateur"
	LevelPro      DancerLevel = "pro"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Level     *DancerLevel
	CreatedAt time.Time
}

type UserRepository interface {
	Store(ctx context.Context, user User) error
	Find(ctx context.Context, userID string) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
}
