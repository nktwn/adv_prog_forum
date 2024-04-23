package models

import (
	"forum/pkg/validator"
	"strconv"
	"time"
)

type Post struct {
	PostID     int
	UserID     int
	UserName   string
	Title      string
	Content    string
	ImageName  string
	Created    time.Time
	Like       int
	Dislike    int
	Comment    *[]Comment
	Categories map[int]string
}

type Comment struct {
	CommentId      int
	PostID         int
	CreatedUserID  int
	Content        string
	CreatedTime    time.Time
	LikeCounter    string
	DislikeCounter string
}

type PostForm struct {
	Title               string   `form:"title"`
	Content             string   `form:"content"`
	Categories          []int    `form:"category"`
	CategoriesString    []string `form:"category"`
	validator.Validator `form:"-"`
}

func (f *PostForm) ConverCategories() error {
	for _, str := range f.CategoriesString {
		nb, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		f.Categories = append(f.Categories, nb)
	}
	return nil
}
