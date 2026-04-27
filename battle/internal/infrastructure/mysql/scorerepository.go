package mysql

import "battle/internal/domain"

type MySqlScoreRepository struct {
}

func (r *MySqlScoreRepository) FindAllForBattle(battleID string) ([]domain.JudgeScore, error) {
	panic("not implemented") // TODO: Implement
}
