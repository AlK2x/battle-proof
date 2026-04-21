package domain

type ScoreRepository interface {
	FindAllForBattle(battleID string) ([]JudgeScore, error)
}

type BattleRepository interface {
	FindBattle(battleID string) (Battle, error)
}
