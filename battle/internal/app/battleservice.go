package app

import (
	"battle/internal/domain"
	"context"
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

func (bs *BattleService) CreateBattle(ctx context.Context, params CreateBattleData) error {
	command := bs.commandFactory.CreateStartBattleCommand(
		ctx,
		params.EventID,
		params.Dancer1ID,
		params.Dancer2ID,
	)
	command.Execute()
	return command.Err
}

func (bs *BattleService) SubmitScore() error {
	return nil
}

func (bs *BattleService) FinishBattle(battleID string) (*domain.BattleResult, error) {
	command := bs.commandFactory.CreateFinishBattleCommand(battleID)
	command.Execute()

	return command.Result, command.Err
}
