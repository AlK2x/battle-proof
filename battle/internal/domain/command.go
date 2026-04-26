package domain

type Command interface {
	Execute()
}

type CommandFactory struct {
	winnerStrategy   WinnerStrategy
	battleRepository BattleRepository
	scoreRepository  ScoreRepository
}

func (cf *CommandFactory) CreateFinishBattleCommand(ID string) *FinishBattleCommand {
	return &FinishBattleCommand{
		battleID: ID,
	}
}

type StartBattleCommand struct {
}

func (sb *StartBattleCommand) Execute() {
	panic("not implemented") // TODO: Implement
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
