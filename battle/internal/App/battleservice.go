package app

import (
	"battle/internal/domain"
	"errors"
)

var (
	ErrEventNotFound       = errors.New("event not found")
	ErrDanceMultipleBattle = errors.New("dancer can participate only in one battle")
	ErrDancerNoDancerRole  = errors.New("dancer don't have dancer role")
	ErrNoSubmittedBattle   = errors.New("no submitted or in progress battle")
)

func NewBattleService(
	commandFactory domain.CommandFactory,
) *BattleService {
	return &BattleService{
		commandFactory: commandFactory,
	}
}

type BattleService struct {
	commandFactory domain.CommandFactory
}

func (bs *BattleService) CreateBattle() error {
	return nil
}

func (bs *BattleService) SubmitScore() error {
	return nil
}

func (bs *BattleService) FinishBattle(battleID string) (*domain.BattleResult, error) {
	command := bs.commandFactory.CreateFinishBattleCommand(battleID)
	command.Execute()

	return command.Result, command.Err
}
