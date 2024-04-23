package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
)

func (s *Sqlite) GetUserIDByToken(token string) (int, error) {
	op := "sqlite.GetUserIDByToken"
	stmt := `SELECT user_id FROM sessions WHERE token = ?`
	var userID int

	err := s.db.QueryRow(stmt, token).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, models.ErrNoRecord
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (s *Sqlite) CreateSession(session *models.Session) error {
	op := "sqlite.CreateSession"
	stmt := `INSERT INTO sessions(user_id, token, exp_time) VALUES(?, ?, ?)`
	_, err := s.db.Exec(stmt, session.UserID, session.Token, session.ExpTime)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Sqlite) DeleteSessionByUserID(userID int) error {
	op := "sqlite.DeleteSessionByUserID"
	stmt := `DELETE FROM sessions WHERE user_id = ?`
	if _, err := s.db.Exec(stmt, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}

func (s *Sqlite) DeleteSessionByToken(token string) error {
	op := "sqlite.DeleteSessionByToken"
	stmt := `DELETE FROM sessions WHERE token = ?`
	if _, err := s.db.Exec(stmt, token); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
