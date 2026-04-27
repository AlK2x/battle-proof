package main

import (
	app "battle/internal/App"
	"battle/internal/domain"
	"battle/internal/infrastructure/mysql"
)

func CreateBattleService() *app.BattleService {
	winStrategy := &domain.SumScoreWinnerStrategy{}
	battleRepo := &mysql.MySqlBattleRepository{}
	scoreRepo := &mysql.MySqlScoreRepository{}
	commandFactory := domain.NewCommandFactory(battleRepo, scoreRepo, winStrategy)
	return app.NewBattleService(commandFactory)
}
