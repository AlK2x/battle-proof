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
		battleID:          battle.ID,
		winnerID:          winner,
		dancer1TotalScore: dancer1Score,
		dander2TotalScore: dancer2Score,
	}
}
