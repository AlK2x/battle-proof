package main

import (
	app "battle/internal/app"
	"battle/internal/domain"
)

func CreateBattleService(battleRepo domain.BattleRepository, scoreRepo domain.ScoreRepository) *app.BattleService {
	winStrategy := &domain.SumScoreWinnerStrategy{}
	commandFactory := domain.NewCommandFactory(battleRepo, scoreRepo, winStrategy)
	return app.NewBattleService(commandFactory)
}
