package domain

import (
	"context"

	"github.com/google/uuid"
)

type Command interface {
	Execute()
}

func NewCommandFactory(
	battleRepo BattleRepository,
	scoreRepo ScoreRepository,
	winStrategy WinnerStrategy,
) CommandFactory {
	return CommandFactory{
		winnerStrategy:   winStrategy,
		battleRepository: battleRepo,
		scoreRepository:  scoreRepo,
	}
}

type CommandFactory struct {
	winnerStrategy   WinnerStrategy
	battleRepository BattleRepository
	scoreRepository  ScoreRepository
}

func (cf *CommandFactory) CreateFinishBattleCommand(ID string) *FinishBattleCommand {
	return &FinishBattleCommand{
		winnerStrategy:   cf.winnerStrategy,
		battleRepository: cf.battleRepository,
		scoreRepository:  cf.scoreRepository,
		battleID:         ID,
	}
}

func (cf *CommandFactory) CreateStartBattleCommand(ctx context.Context, eventID, dancer1ID, dancer2ID string) *StartBattleCommand {
	return &StartBattleCommand{
		eventID:          eventID,
		dancer1ID:        dancer1ID,
		dancer2ID:        dancer2ID,
		battleRepository: cf.battleRepository,

		ctx: ctx,
	}
}

type StartBattleCommand struct {
	eventID          string
	dancer1ID        string
	dancer2ID        string
	battleRepository BattleRepository
	Err              error

	ctx context.Context
}

func (sb *StartBattleCommand) Execute() {
	battle := Battle{
		ID:      uuid.NewString(),
		Status:  StatusPending,
		EventID: sb.eventID,
		Dancer1: sb.dancer1ID,
		Dancer2: sb.dancer2ID,
	}
	sb.Err = sb.battleRepository.Store(sb.ctx, battle)
}

type SubmitScoreCommand struct {
}

func (ss *SubmitScoreCommand) Execute() {
	panic("not implemented") // TODO: Implement
}

type FinishBattleCommand struct {
	battleID         string
	winnerStrategy   WinnerStrategy
	battleRepository BattleRepository
	scoreRepository  ScoreRepository

	ctx    context.Context
	Err    error
	Result *BattleResult
}

func (fb *FinishBattleCommand) Execute() {
	battle, err := fb.battleRepository.FindBattle(fb.ctx, fb.battleID)
	if err != nil {
		fb.Err = err
		return
	}
	scores, err := fb.scoreRepository.FindAllForBattle(fb.ctx, fb.battleID)
	if err != nil {
		fb.Err = err
		return
	}
	result := fb.winnerStrategy.ChooseWinner(battle, scores...)
	fb.Result = &result
}
