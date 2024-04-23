package repo

import (
	"forum/internal/repo/sqlite"
	"forum/models"
)

type UserRepo interface {
	CreateUser(*models.User) error
	GetUserByID(int) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	// UpdateUserByID(string) (*models.User, error)
	Authenticate(email, password string) (int, error)
	UpdateUserPassword(id int, password string) error
	UpdateUserEmail(id int, email string) error
	UpdateUserName(id int, name string) error
	GetAllUsers() ([]*models.User, error)
	DeleteUser(int) error
}

type SessionRepo interface {
	GetUserIDByToken(string) (int, error)
	CreateSession(*models.Session) error
	DeleteSessionByUserID(int) error
	DeleteSessionByToken(string) error
}

type PostRepo interface {
	CreatePost(userID int, title, content, imageName string) (int, error)
	GetPostByID(int) (*models.Post, error)
	GetCategoriesByPostID(int) (map[int]string, error)
	// GetAllPost() (*models.Post, error)
	// UpdatePost(string, *models.Post) error
	//AddLikeAndDislike(bool, string, string) error
	// DeleteLikeAndDislike(int, int) error
	GetAllPostByUserID(int) (*[]models.Post, error)
	GetAllPostByCategories(categories []int) (*[]models.Post, error)
	GetPageNumber(pageSize int) (int, error)
	GetAllPostPaginated(page int, pageSize int) (*[]models.Post, error)
}

type CategoryRepo interface {
	CreateParsePosts() error
	AddCategoryToPost(int, []int) error
	GetALLCategory() ([]string, error)
	// CreateCategory(string) error
}

// type CommentRepo interface {
// 	CreateComment(*models.Comment) error
// 	GetAllCommentByPostID(string) (*[]models.Post, error)
// 	GetAllCommentByUserID(string) (*[]models.Post, error)
// 	AddLikeAndDislike(bool, string, string) error
// }

type RepoI interface {
	UserRepo
	SessionRepo
	PostRepo
	CategoryRepo
	GetUserByActivationToken(activationToken string) (*models.User, error)
	ActivateUser(userID int) error
	// CommentRepo
}

func New(storagePath string) (RepoI, error) {
	sqliteDB, err := sqlite.NewDB(storagePath)
	if err != nil {
		return nil, err
	}
	return sqliteDB, nil
}
