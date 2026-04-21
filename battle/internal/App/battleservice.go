package app

import "battle/internal/domain"

func NewBattleService(
	winnerStrategy domain.WinnerStrategy,
	battleRepository domain.BattleRepository,
	scoreRepository domain.ScoreRepository,
) *BattleService {
	return &BattleService{
		winnerStrategy: winnerStrategy,
	}
}

type BattleService struct {
	winnerStrategy   domain.WinnerStrategy
	battleRepository domain.BattleRepository
	scoreRepository  domain.ScoreRepository
}

func (bs *BattleService) FinishBattle(battleID string) (*domain.BattleResult, error) {
	battle, err := bs.battleRepository.FindBattle(battleID)
	if err != nil {
		return nil, err
	}
	scores, err := bs.scoreRepository.FindAllForBattle(battleID)
	if err != nil {
		return nil, err
	}
	result := bs.winnerStrategy.ChooseWinner(battle, scores...)
	return &result, nil
}
