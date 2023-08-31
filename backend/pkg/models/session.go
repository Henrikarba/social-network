package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Session struct {
	UserID    int    `db:"user_id"`
	UUID      string `db:"uuid"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func ValidateSession(db *sqlx.DB, id int, uuid string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM tokens WHERE user_id = ? AND uuid = ?`
	err := db.Get(&count, query, id, uuid)
	if err != nil {
		return false, fmt.Errorf("error checking UUID and UserID match: %v", err)
	}
	return count > 0, nil
}

func NewSession(db *sqlx.DB, s Session) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Printf("Transaction rolled back due to an error: %v", r)
		}
	}()

	_, err = tx.NamedExec(`INSERT INTO tokens (user_id, uuid) VALUES (:user_id, :uuid)`, s)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert token: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
