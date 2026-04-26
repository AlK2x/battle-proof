package domain

type SumScoreWinnerStrategy struct {
}

func (s *SumScoreWinnerStrategy) ChooseWinner(battle Battle, scores ...JudgeScore) BattleResult {
	var (
		dancer1Score int
		dancer2Score int
		winner       *string
	)
	for _, s := range scores {
		dancer1Score += s.Dancer1Score
		dancer2Score += s.Dancer2Score
	}
	if dancer1Score == dancer2Score {
		winner = nil
	} else if dancer1Score > dancer2Score {
		winner = &battle.Dancer1
	} else {
		winner = &battle.Dancer2
	}
	return BattleResult{
		BattleID:          battle.ID,
		WinnerID:          winner,
		Dancer1TotalScore: dancer1Score,
		Dander2TotalScore: dancer2Score,
	}
}
