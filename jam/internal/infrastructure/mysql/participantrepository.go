package mysql

import (
	"context"
	"database/sql"
	"errors"
	"jam/internal/domain/jam"
	"time"
)

var _ jam.JamRepository = (*MySQLJamRepository)(nil)
var _ jam.ParticipantRepository = (*MySQLParticipantRepository)(nil)

type MySQLJamRepository struct {
	db *sql.DB
}

func NewMySQLJamRepository(db *sql.DB) jam.JamRepository {
	return &MySQLJamRepository{
		db: db,
	}
}

// FindAll implements [jam.JamRepository].
func (m *MySQLJamRepository) FindAll(ctx context.Context) ([]*jam.Jam, error) {
	query := `SELECT id, name, location, date, created_by FROM jam`
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var (
		ID        string
		Name      string
		Location  string
		Date      time.Time
		CreatedBy string
	)
	result := make([]*jam.Jam, 0)
	for rows.Next() {
		err := rows.Scan(&ID, &Name, &Location, &Date, &CreatedBy)
		if err != nil {
			return result, err
		}
		result = append(result, &jam.Jam{
			ID:        jam.JamID(ID),
			Name:      Name,
			Location:  Location,
			Date:      Date,
			CreatedBy: jam.UserID(CreatedBy),
		})
	}
	return result, rows.Err()
}

// FindByID implements [jam.JamRepository].
func (m *MySQLJamRepository) FindByID(ctx context.Context, id jam.JamID) (*jam.Jam, error) {
	query := `SELECT id, name, location, date, created_by FROM jam WHERE id = ?`
	row := m.db.QueryRowContext(ctx, query, id)
	var (
		ID        string
		Name      string
		Location  string
		Date      time.Time
		CreatedBy string
	)
	err := row.Scan(&ID, &Name, &Location, &Date, &CreatedBy)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &jam.Jam{
		ID:        jam.JamID(ID),
		Name:      Name,
		Location:  Location,
		Date:      Date,
		CreatedBy: jam.UserID(CreatedBy),
	}, nil
}

// Store implements [jam.JamRepository].
func (m *MySQLJamRepository) Store(ctx context.Context, p *jam.Jam) error {
	query := `INSERT INTO jam (id, name, location, date, created_by)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			location = VALUES(location),
			date = VALUES(date),
			created_by = VALUES(created_by)
	`
	_, err := m.db.ExecContext(ctx, query, p.ID, p.Name, p.Location, p.Date, p.CreatedBy)
	if err != nil {
		return err
	}
	return nil
}

func NewMySQLParticipantRepository(db *sql.DB) *MySQLParticipantRepository {
	return &MySQLParticipantRepository{
		db: db,
	}
}

type MySQLParticipantRepository struct {
	db *sql.DB
}

// DeleteByJamID implements [jam.ParticipantRepository].
func (m *MySQLParticipantRepository) DeleteByJamID(ctx context.Context, jamID jam.JamID) error {
	query := `DELETE FROM participant WHERE jam_id = ?`
	_, err := m.db.ExecContext(ctx, query, jamID)
	if err != nil {
		return err
	}
	return nil
}

// FindByJamID implements [jam.ParticipantRepository].
func (m *MySQLParticipantRepository) FindByJamID(ctx context.Context, ID jam.JamID) ([]*jam.Participant, error) {
	query := `SELECT jam_id, user_id, role FROM participant WHERE jam_id = ?`
	rows, err := m.db.QueryContext(ctx, query, ID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	defer func() { _ = rows.Close() }()
	var (
		jamID  string
		userID string
		role   string
	)
	result := make([]*jam.Participant, 0)
	for rows.Next() {
		err := rows.Scan(&jamID, &userID, &role)
		if err != nil {
			return result, err
		}
		result = append(result, &jam.Participant{
			JamID:  jam.JamID(jamID),
			UserID: jam.UserID(userID),
			Role:   jam.JamRole(role),
		})
	}
	return result, rows.Err()
}

// Store implements [jam.ParticipantRepository].
func (m *MySQLParticipantRepository) Store(ctx context.Context, p *jam.Participant) error {
	query := `INSERT UPDATE INTO participant (jam_id, user_id, role)
	          VALUES (?, ?, ?)
			  ON DUPLICATE KEY UPDATE
			  role = VALUES(role)
	`
	_, err := m.db.ExecContext(ctx, query, p.JamID, p.UserID, p.Role)
	if err != nil {
		return err
	}
	return nil
}
