package jam

import (
	"time"

	"github.com/google/uuid"
)

type UserID string
type JamID string
type JamRole string

const (
	JamRoleDancer JamRole = "dancer"
	JamRoleJudge  JamRole = "judge"
	JamRoleMedia  JamRole = "media"
)

type Jam struct {
	ID        JamID
	Name      string
	Location  string
	Date      time.Time
	CreatedBy UserID
}

func NewJam(name string, location string, date time.Time, createdBy UserID) Jam {
	id := JamID(uuid.NewString())
	return Jam{
		ID:        id,
		Name:      name,
		Location:  location,
		Date:      date,
		CreatedBy: createdBy,
	}
}

type Participant struct {
	JamID  JamID
	UserID UserID
	Role   JamRole
}

type JamRepository interface {
	Store(p *Jam) error
	FindByID(id JamID) (*Jam, error)
	FindAll() ([]*Jam, error)
}

type ParticipantRepository interface {
	Store(p *Participant) error
	DeleteByJamID(jamID JamID) error
	FindByJamID(id JamID) ([]*Participant, error)
}
