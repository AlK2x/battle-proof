package domain

import "time"

type BattleResult struct {
	BattleID          string
	WinnerID          *string
	Dancer1TotalScore int
	Dander2TotalScore int
	FinishedAt        *time.Time
	Status            BattleStatus
}

type WinnerStrategy interface {
	ChooseWinner(battle Battle, scores ...JudgeScore) BattleResult
}
