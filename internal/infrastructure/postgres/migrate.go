package postgres

import (
	"database/sql"
	"fmt"
	"os"
)

func RunMigration(db *sql.DB, dir string) error {
	var sqlBytes []byte
	var err error

	switch dir {
	case "up":
		sqlBytes, err = os.ReadFile("migrations/model_up.sql")
	case "down":
		sqlBytes, err = os.ReadFile("migrations/model_down.sql")
	default:
		return fmt.Errorf("unknown migration direction: %s (must be 'up' or 'down')", dir)
	}

	if err != nil {
		return fmt.Errorf("failed to read: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	if _, err := db.Exec(string(sqlBytes)); err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return fmt.Errorf("failed to rollback: %w", rollBackErr)
		}
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	fmt.Printf("All %s migrations applied successfully!\n", dir)
	return nil
}
