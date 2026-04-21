package domain

type BattleResult struct {
	battleID          string
	winnerID          *string
	dancer1TotalScore int
	dander2TotalScore int
}

type WinnerStrategy interface {
	ChooseWinner(battle Battle, scores ...JudgeScore) BattleResult
}
