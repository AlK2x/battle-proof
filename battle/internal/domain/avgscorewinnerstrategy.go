package domain

type AvgScoreWinnerStrategy struct {
}

func (s *AvgScoreWinnerStrategy) ChooseWinner(battle Battle, scores ...JudgeScore) BattleResult {
	var (
		dancer1Score int
		dancer2Score int
		dancer1Avg   float64
		dancer2Avg   float64
		winner       *string
	)
	cnt := float64(len(scores))
	for _, s := range scores {
		dancer1Score += s.Dancer1Score
		dancer2Score += s.Dancer2Score
	}
	dancer1Avg = float64(dancer1Score) / cnt
	dancer2Avg = float64(dancer2Score) / cnt
	if dancer1Avg == dancer2Avg {
		winner = nil
	} else if dancer1Avg > dancer2Avg {
		winner = &battle.Dancer1
	} else {
		winner = &battle.Dancer2
	}
	return BattleResult{
		BattleID:          battle.ID,
		WinnerID:          winner,
		Dancer1TotalScore: int(dancer1Avg),
		Dander2TotalScore: int(dancer2Avg),
	}
}
