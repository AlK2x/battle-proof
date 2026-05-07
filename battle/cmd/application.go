package main

import (
	"battle/internal/app"
	"battle/internal/domain"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func CreateBattleService(battleRepo domain.BattleRepository, scoreRepo domain.ScoreRepository) *app.BattleService {
	winStrategy := &domain.SumScoreWinnerStrategy{}
	commandFactory := domain.NewCommandFactory(battleRepo, scoreRepo, winStrategy)
	return app.NewBattleService(commandFactory)
}
