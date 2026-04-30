package domain

import (
	"context"
	"errors"
)

var ErrNoRows = errors.New("err_no_rows")

type ScoreRepository interface {
	FindAllForBattle(ctx context.Context, battleID string) ([]JudgeScore, error)
}

type BattleRepository interface {
	FindBattle(ctx context.Context, battleID string) (*Battle, error)
	Store(ctx context.Context, battle Battle) error
}
