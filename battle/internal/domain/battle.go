package domain

type BattleStatus string

const (
	StatusPending  BattleStatus = "pending"
	StatusProgress BattleStatus = "in_progress"
	StatusFinished BattleStatus = "finished"
)

type Battle struct {
	ID        string
	EventID   string
	Dancer1   string
	Dancer2   string
	WinnnerID *string
	Status    BattleStatus
}

type JudgeScore struct {
	BattleID     string
	JudgeID      string
	Dancer1Score int
	Dancer2Score int
}
