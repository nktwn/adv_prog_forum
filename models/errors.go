package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")

	ErrNotActivated = errors.New("models: user not activated")
)
