package service

import (
	"forum/internal/repo"
	"forum/models"
)

type service struct {
	repo repo.RepoI
}

type ServiceI interface {
	UserServiceI
	CategoryServiceI
	PostServiceI
	GetAllUsers() ([]models.User, error)
	DeleteUser(int) error
	ActivateUser(token string) error
}

type UserServiceI interface {
	GetUser(int) *models.User
	CreateUser(*models.User) error
	Authenticate(string, string) (*models.Session, error)
	DeleteSession(string) error
	UpdateUserPassword(token string, newPassword string) error
	GetUserByToken(token string) (*models.User, error)
}

type PostServiceI interface {
	CreatePost(string, string, string, []int) (int, error)
	GetPostByID(int) (*models.Post, error)
	GetAllPostPaginated(int, int) (*[]models.Post, error)
	GetPageNumber(int) (int, error)
	GetAllPostByCategories(categories []int) (*[]models.Post, error)
	GetAllPostByUser(token string) (*[]models.Post, error)
}

type CategoryServiceI interface {
	GetAllCategory() ([]string, error)
}

func New(r repo.RepoI) ServiceI {
	return &service{
		r,
	}
}
