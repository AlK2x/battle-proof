package domain

import "github.com/google/uuid"

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

func (cf *CommandFactory) CreateStartBattleCommand(eventID, dancer1ID, dancer2ID string) *StartBattleCommand {
	return &StartBattleCommand{
		eventID:          eventID,
		dancer1ID:        dancer1ID,
		dancer2ID:        dancer2ID,
		battleRepository: cf.battleRepository,
	}
}

type StartBattleCommand struct {
	eventID          string
	dancer1ID        string
	dancer2ID        string
	battleRepository BattleRepository
	Err              error
}

func (sb *StartBattleCommand) Execute() {
	battle := Battle{
		ID:      uuid.NewString(),
		status:  StatusPending,
		EventID: sb.eventID,
		Dancer1: sb.dancer1ID,
		Dancer2: sb.dancer2ID,
	}
	sb.Err = sb.battleRepository.Store(battle)
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

	Err    error
	Result *BattleResult
}

func (fb *FinishBattleCommand) Execute() {
	battle, err := fb.battleRepository.FindBattle(fb.battleID)
	if err != nil {
		fb.Err = err
		return
	}
	scores, err := fb.scoreRepository.FindAllForBattle(fb.battleID)
	if err != nil {
		fb.Err = err
		return
	}
	result := fb.winnerStrategy.ChooseWinner(battle, scores...)
	fb.Result = &result
}
