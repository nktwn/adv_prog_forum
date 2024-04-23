package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
	"strings"
)

func (s *Sqlite) CreatePost(userID int, title, content, imageName string) (int, error) {
	op := "sqlite.CreatePost"
	const query = `INSERT INTO posts (user_id, title, content, image_name) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, userID, title, content, imageName)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return int(postID), nil
}

func (s *Sqlite) GetPostByID(postID int) (*models.Post, error) {
	op := "sqlite.GetPostByID"
	stmt := `SELECT id, user_id, title, content, created, like, dislike, image_name FROM posts WHERE id = ? ORDER BY created asc`
	post := models.Post{}

	err := s.db.QueryRow(stmt, postID).Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Like, &post.Dislike, &post.ImageName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &post, nil
}

func (s *Sqlite) GetAllPost() ([]models.Post, error) {
	const query = `SELECT post_id, user_id, title, content, created, like, dislike, image_name FROM Post`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Like, &post.Dislike, &post.ImageName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *Sqlite) AddLikeAndDislike(isLike bool, userID, postID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Check if the user has already liked/disliked the post
	var isExists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM Post_User_Like WHERE user_id = ? AND post_id = ?)`
	err = tx.QueryRow(checkQuery, userID, postID).Scan(&isExists)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !isExists {
		// Insert like/dislike
		insertQuery := `INSERT INTO Post_User_Like (user_id, post_id, is_like) VALUES (?, ?, ?)`
		_, err := tx.Exec(insertQuery, userID, postID, isLike)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Update Post like/dislike count
		updateQuery := ""
		if isLike {
			updateQuery = `UPDATE Post SET like = like + 1 WHERE post_id = ?`
		} else {
			updateQuery = `UPDATE Post SET dislike = dislike + 1 WHERE post_id = ?`
		}
		_, err = tx.Exec(updateQuery, postID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *Sqlite) DeleteLikeAndDislike(userID, postID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	//is user liked or disliked
	var isLike bool
	checkQuery := `SELECT is_like FROM Post_User_Like WHERE user_id = ? AND post_id = ?`
	err = tx.QueryRow(checkQuery, userID, postID).Scan(&isLike)
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete the like/dislike
	deleteQuery := `DELETE FROM Post_User_Like WHERE user_id = ? AND post_id = ?`
	_, err = tx.Exec(deleteQuery, userID, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// decrement the like or dislike
	updateQuery := ""
	if isLike {
		updateQuery = `UPDATE Post SET like = like - 1 WHERE post_id = ? AND like > 0`
	} else {
		updateQuery = `UPDATE Post SET dislike = dislike - 1 WHERE post_id = ? AND dislike > 0`
	}
	_, err = tx.Exec(updateQuery, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Sqlite) GetAllPostByUserID(userID int) (*[]models.Post, error) {
	const query = `SELECT id, user_id, title, content, created, like, dislike, image_name FROM posts WHERE user_id=?`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Like, &post.Dislike, &post.ImageName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (s *Sqlite) GetAllPostByCategories(categoryIDs []int) (*[]models.Post, error) {
	query := `SELECT p.id, p.user_id, p.title, p.content, p.created, p.like, p.dislike, p.image_name, u.name
              FROM posts AS p
              INNER JOIN post_category AS pc ON p.id = pc.post_id
              JOIN users u ON p.user_id = u.id
              WHERE pc.category_id IN (?` + strings.Repeat(",?", len(categoryIDs)-1) + `)
              GROUP BY p.id`

	args := make([]interface{}, len(categoryIDs))
	for i, id := range categoryIDs {
		args[i] = id
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Like, &post.Dislike, &post.ImageName, &post.UserName); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (s *Sqlite) GetAllPostPaginated(page int, pageSize int) (*[]models.Post, error) {
	op := "sqlite.GetAllPostPaginated"
	offset := (page - 1) * pageSize
	stmt := `SELECT p.id, p.user_id, p.title, p.content, p.created, p.like, p.dislike, p.image_name, u.name 
	FROM posts p 
	JOIN users u ON p.user_id = u.id 
	LIMIT ? OFFSET ?`

	rows, err := s.db.Query(stmt, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Created, &post.Like, &post.Dislike, &post.ImageName, &post.UserName); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, post)
	}
	return &posts, nil
}

func (s *Sqlite) GetPageNumber(pageSize int) (int, error) {
	op := "sqlite.GetPageNumber"
	stmt := `SELECT COUNT(*) FROM posts`

	var totalPosts int
	err := s.db.QueryRow(stmt).Scan(&totalPosts)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	totalPages := (totalPosts + pageSize - 1) / pageSize

	return totalPages, nil
}
