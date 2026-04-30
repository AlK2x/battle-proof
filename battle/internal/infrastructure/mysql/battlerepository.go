package mysql

import (
	"battle/internal/domain"
	"context"
	"database/sql"
)

func NewMysqlBattleRepository(db *sql.DB) *MySqlBattleRepository {
	return &MySqlBattleRepository{
		db: db,
	}
}

type MySqlBattleRepository struct {
	db *sql.DB
}

func (r *MySqlBattleRepository) FindBattle(ctx context.Context, battleID string) (*domain.Battle, error) {
	query := "SELECT id, event_id, dancer1, dancer2, winner_id, status FROM battle WHERE battle_id = ?"
	rows, err := r.db.QueryContext(ctx, query, battleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		ID       string
		eventID  string
		dancer1  string
		dancer2  string
		winnerID *string
		status   domain.BattleStatus
	)
	for rows.Next() {
		err := rows.Scan(&ID, &eventID, &dancer1, &dancer2, winnerID, &status)
		if err != nil {
			return nil, err
		}
		result := domain.Battle{
			ID:        ID,
			EventID:   eventID,
			Dancer1:   dancer1,
			Dancer2:   dancer2,
			WinnnerID: winnerID,
			Status:    status,
		}
		return &result, nil
	}

	return nil, domain.ErrNoRows
}

func (r *MySqlBattleRepository) Store(ctx context.Context, battle domain.Battle) error {
	query := `
		INSERT INTO battle (id, event_id, dancer1, dancer2, winner_id, status)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			event_id = VALUES(event_id),
			dancer1 = VALUES(dancer1),
			dancer2 = VALUES(dancer2),
			winner_id = VALUES(winner_id),
			status = VALUES(status)
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, battle.ID, battle.EventID, battle.Dancer1, battle.Dancer2, battle.WinnnerID, battle.GetStatus())
	if err != nil {
		return err
	}
	return nil
}
