package domain

type Battle struct {
	ID        string
	EventID   string
	Dancer1   string
	Dancer2   string
	WinnnerID *string
	status    BattleStatus

	sm *BattleStateMachine
}

func (b *Battle) GetStatus() BattleStatus {
	return b.status
}

func (b *Battle) Start() error {
	return b.sm.Change(StatusProgress)
}

func (b *Battle) Finish() error {
	return b.sm.Change(StatusFinished)
}

type JudgeScore struct {
	BattleID     string
	JudgeID      string
	Dancer1Score int
	Dancer2Score int
}
