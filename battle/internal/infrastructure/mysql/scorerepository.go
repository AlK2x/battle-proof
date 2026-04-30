package mysql

import (
	"battle/internal/domain"
	"context"
	"database/sql"
)

func NewMysqlScoreRepository(db *sql.DB) *MySqlScoreRepository {
	return &MySqlScoreRepository{
		db: db,
	}
}

type MySqlScoreRepository struct {
	db *sql.DB
}

func (r *MySqlScoreRepository) FindAllForBattle(ctx context.Context, battleID string) ([]domain.JudgeScore, error) {
	query := "SELECT battle_id, judge_id, dancer1_score, dancer2_score FROM judge_score WHERE battle_id = ?"

	rows, err := r.db.QueryContext(ctx, query, battleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		battlID      string
		judgeID      string
		dancer1Score int
		dancer2Score int
	)
	result := make([]domain.JudgeScore, 0)
	for rows.Next() {
		err := rows.Scan(&battlID, &judgeID, &dancer1Score, &dancer2Score)
		if err != nil {
			return nil, err
		}
		result = append(result, domain.JudgeScore{
			BattleID:     battlID,
			JudgeID:      judgeID,
			Dancer1Score: dancer1Score,
			Dancer2Score: dancer2Score,
		})
	}

	return result, rows.Err()
}
