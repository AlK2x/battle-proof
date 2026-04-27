package mysql

import "battle/internal/domain"

type MySqlBattleRepository struct {
}

func (r *MySqlBattleRepository) FindBattle(battleID string) (domain.Battle, error) {
	panic("not implemented") // TODO: Implement
}

func (r *MySqlBattleRepository) Store(battle domain.Battle) error {
	panic("not implemented") // TODO: Implement
}
