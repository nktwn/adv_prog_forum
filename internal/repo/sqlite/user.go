package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

func (s *Sqlite) GetUserByEmail(email string) (*models.User, error) {
	op := "sqlite.GetUserByEmail"
	var u models.User
	stmt := `SELECT id, name, email, created FROM users WHERE id=?`
	err := s.db.QueryRow(stmt, email).Scan(&u.ID, &u.Name, &u.Email, &u.Created)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &u, nil

}
func (s *Sqlite) CreateUserAndReturnID(u models.User) (int64, error) {
	op := "sqlite.CreateUser"
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES (?, ?, ?, CURRENT_TIMESTAMP)`
	result, err := s.db.Exec(stmt, u.Name, u.Email, string(u.HashedPassword))
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return 0, models.ErrDuplicateEmail
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: unable to get last insert ID: %w", op, err)
	}
	return id, nil
}

func (s *Sqlite) CreateUser(u *models.User) error {
	u.ActivationToken = generateActivationToken()
	op := "sqlite.CreateUser"
	stmt := `INSERT INTO users (name, email, hashed_password, is_activated, activation_token, created) VALUES(?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`
	_, err := s.db.Exec(stmt, u.Name, u.Email, string(u.HashedPassword), false, u.ActivationToken)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return models.ErrDuplicateEmail
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Sqlite) GetUserByID(id int) (*models.User, error) {
	op := "sqlite.GetUserByID"
	var u models.User
	stmt := `SELECT id, name, email, created FROM users WHERE id=?`
	err := s.db.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &u, nil
}

func (s *Sqlite) Authenticate(email, password string) (int, error) {
	op := "sqlite.Authenticate"
	var id int
	var hashedPassword []byte
	var isActivated bool

	stmt := `SELECT id, hashed_password, is_activated FROM users WHERE email=?`
	err := s.db.QueryRow(stmt, email).Scan(&id, &hashedPassword, &isActivated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if !isActivated {
		return 0, models.ErrNotActivated
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Sqlite) UpdateUserPassword(id int, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("sqlite.UpdateUserPassword: could not hash password: %w", err)
	}

	stmt := `UPDATE users SET hashed_password = ? WHERE id = ?`
	_, err = s.db.Exec(stmt, hashedPassword, id)
	if err != nil {
		return fmt.Errorf("sqlite.UpdateUserPassword: %w", err)
	}
	return nil
}

func (s *Sqlite) UpdateUserEmail(id int, email string) error {
	stmt := `UPDATE users SET email = ? WHERE id = ?`
	_, err := s.db.Exec(stmt, email, id)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return models.ErrDuplicateEmail
		}
		return fmt.Errorf("sqlite.UpdateUserEmail: %w", err)
	}
	return nil
}

func (s *Sqlite) UpdateUserName(id int, name string) error {
	stmt := `UPDATE users SET name = ? WHERE id = ?`
	_, err := s.db.Exec(stmt, name, id)
	if err != nil {
		return fmt.Errorf("sqlite.UpdateUserName: %w", err)
	}
	return nil
}
func generateActivationToken() string {
	return uuid.New().String()
}

func (s *Sqlite) GetUserByActivationToken(token string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, is_activated FROM users WHERE activation_token = ?`
	err := s.db.QueryRow(query, token).Scan(&user.ID, &user.Name, &user.Email, &user.IsActivated)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Sqlite) ActivateUser(userID int) error {
	query := `UPDATE users SET is_activated = 1, activation_token = '' WHERE id = ?`
	_, err := s.db.Exec(query, userID)
	return err
}
