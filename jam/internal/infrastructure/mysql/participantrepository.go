package mysql

import (
	"database/sql"
	"jam/internal/domain/jam"
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
func (m *MySQLJamRepository) FindAll() ([]*jam.Jam, error) {
	panic("unimplemented")
}

// FindByID implements [jam.JamRepository].
func (m *MySQLJamRepository) FindByID(id jam.JamID) (*jam.Jam, error) {
	panic("unimplemented")
}

// Store implements [jam.JamRepository].
func (m *MySQLJamRepository) Store(p *jam.Jam) error {
	panic("unimplemented")
}

type MySQLParticipantRepository struct {
}

// DeleteByJamID implements [jam.ParticipantRepository].
func (m *MySQLParticipantRepository) DeleteByJamID(jamID jam.JamID) error {
	panic("unimplemented")
}

// FindByJamID implements [jam.ParticipantRepository].
func (m *MySQLParticipantRepository) FindByJamID(id jam.JamID) ([]*jam.Participant, error) {
	panic("unimplemented")
}

// Store implements [jam.ParticipantRepository].
func (m *MySQLParticipantRepository) Store(p *jam.Participant) error {
	panic("unimplemented")
}
