package mysql

import (
	"context"
	"database/sql"
	"time"
	"user/internal/domain"
)

func NewMysqlUserRepository(db *sql.DB) *MysqlUserRepository {
	return &MysqlUserRepository{
		db: db,
	}
}

type MysqlUserRepository struct {
	db *sql.DB
}

func (r *MysqlUserRepository) Store(ctx context.Context, user domain.User) error {
	query := `
		INSERT INTO user (id, name, email, level, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			level = VALUES(level),
			email = VALUES(email)
	`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Level, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *MysqlUserRepository) Find(ctx context.Context, userID string) (*domain.User, error) {
	query := `SELECT id, name, email, level, created_at FROM user WHERE id = ?`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	var (
		ID        string
		name      string
		email     string
		level     *domain.DancerLevel
		createdAt time.Time
	)
	defer rows.Close()
	err = rows.Scan(&ID, &name, &email, &level, &createdAt)
	if err != nil {
		return nil, err
	}
	result := domain.User{
		ID:        ID,
		Name:      name,
		Email:     email,
		Level:     level,
		CreatedAt: createdAt,
	}
	return &result, rows.Err()
}

func (r *MysqlUserRepository) ListUsers(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, name, email, level, created_at FROM user`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var (
		ID        string
		name      string
		email     string
		level     *domain.DancerLevel
		createdAt time.Time
	)
	defer rows.Close()
	result := make([]*domain.User, 0)
	for rows.Next() {
		err := rows.Scan(&ID, &name, &email, &level, &createdAt)
		if err != nil {
			return result, err
		}
		result = append(result, &domain.User{
			ID:        ID,
			Name:      name,
			Email:     email,
			Level:     level,
			CreatedAt: createdAt,
		})
	}
	return result, rows.Err()
}
