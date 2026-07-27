package jam

import (
	"context"
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
	Store(ctx context.Context, p *Jam) error
	FindByID(ctx context.Context, id JamID) (*Jam, error)
	FindAll(ctx context.Context) ([]*Jam, error)
}

type ParticipantRepository interface {
	Store(ctx context.Context, p *Participant) error
	DeleteByJamID(ctx context.Context, jamID JamID) error
	FindByJamID(ctx context.Context, id JamID) ([]*Participant, error)
}
