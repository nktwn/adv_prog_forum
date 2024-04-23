package service

import (
	"fmt"
	"forum/models"
	"net/smtp"
)

func (s *service) GetUser(id int) *models.User {
	return nil
}

func (s *service) DeleteSession(token string) error {
	if err := s.repo.DeleteSessionByToken(token); err != nil {
		return err
	}
	return nil
}

func (s *service) Authenticate(email string, password string) (*models.Session, error) {
	userID, err := s.repo.Authenticate(email, password)
	if err != nil {
		return nil, err
	}
	session := models.NewSession(userID)

	if err = s.repo.DeleteSessionByUserID(userID); err != nil {
		return nil, err
	}

	if err = s.repo.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}
func (s *service) CreateUser(user *models.User) error {
	err := s.repo.CreateUser(user) // Разыменование указателя

	if err != nil {
		return err
	}

	err = sendActivationEmail(user.Email, user.ActivationToken)
	if err != nil {
		return err // Обработка ошибки отправки email
	}
	return nil
}

func (s *service) GetUserByToken(token string) (*models.User, error) {
	userID, err := s.repo.GetUserIDByToken(token)

	if err != nil {
		return nil, err
	}

	return s.repo.GetUserByID(userID)
}

func (s *service) UpdateUserPassword(token string, newPassword string) error {

	userID, err := s.repo.GetUserIDByToken(token)
	if err != nil {
		return err
	}

	return s.repo.UpdateUserPassword(userID, newPassword)
}
func (s *service) GetAllUsers() ([]models.User, error) {
	userPointers, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	users := make([]models.User, len(userPointers))
	for i, userPointer := range userPointers {
		users[i] = *userPointer
	}

	return users, nil
}

func (s *service) DeleteUser(userID int) error {
	return s.repo.DeleteUser(userID)
}

const (
	smtpServer   = "smtp.mail.ru"
	smtpPort     = "587"
	smtpUser     = "nurkhat.sergaziev@mail.ru"
	smtpPassword = "itFCzPK0wcRVDxzTeDaD"
)

func sendActivationEmail(email, token string) error {
	from := smtpUser
	to := []string{email}
	url := "http://localhost:8080/activate?token=" + token
	body := fmt.Sprintf("To activate your account, please click on the following link: %s", url)
	msg := "From: " + from + "\n" +
		"To: " + email + "\n" +
		"Subject: Activate Your Account\n\n" +
		body

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, to, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ActivateUser(token string) error {
	user, err := s.repo.GetUserByActivationToken(token)
	if err != nil {
		return err // Пользователь не найден или другая ошибка
	}

	if user.IsActivated {
		return fmt.Errorf("user already activated")
	}

	return s.repo.ActivateUser(user.ID)
}
