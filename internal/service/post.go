package service

import (
	"forum/models"
)

func (s *service) CreatePost(title, content, token string, categories []int) (int, error) {
	userID, err := s.repo.GetUserIDByToken(token)
	if err != nil {
		return 0, err
	}

	postID, err := s.repo.CreatePost(userID, title, content, "Nan")
	if err != nil {
		return 0, err
	}

	if err = s.repo.AddCategoryToPost(postID, AddCategory(categories)); err != nil {
		return 0, err
	}
	return postID, err
}

func (s *service) GetPostByID(id int) (*models.Post, error) {
	post, err := s.repo.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.GetCategoriesByPostID(id)
	if err != nil {
		return nil, err
	}
	post.Categories = categories
	return post, nil
}

func (s *service) GetAllPostPaginated(curentPage, pageSize int) (*[]models.Post, error) {
	posts, err := s.repo.GetAllPostPaginated(curentPage, pageSize)
	if err != nil {
		return nil, err
	}
	if err = s.getCategoryToPost(posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *service) GetPageNumber(pageSize int) (int, error) {
	return s.repo.GetPageNumber(pageSize)
}

func (s *service) GetAllPostByCategories(categories []int) (*[]models.Post, error) {
	posts, err := s.repo.GetAllPostByCategories(AddCategory(categories))
	if err != nil {
		return nil, err
	}

	if err = s.getCategoryToPost(posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *service) GetAllPostByUser(token string) (*[]models.Post, error) {
	userID, err := s.repo.GetUserIDByToken(token)
	if err != nil {
		return nil, err
	}
	posts, err := s.repo.GetAllPostByUserID(userID)

	if err != nil {
		return nil, err
	}
	if err = s.getCategoryToPost(posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *service) getCategoryToPost(posts *[]models.Post) error {
	for i := range *posts {
		categories, err := s.repo.GetCategoriesByPostID((*posts)[i].PostID)
		if err != nil {
			return err

		}
		(*posts)[i].Categories = categories
	}
	return nil
}
