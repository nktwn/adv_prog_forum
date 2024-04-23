package sqlite

import (
	"encoding/json"
	"errors"
	"forum/models"
	"math/rand"
	"net/http"
	"time"
)

const urlToParse = "https://newsapi.org/v2/everything?q=tesla&from=2024-02-01&sortBy=publishedAt&apiKey=fee3805d20864aac832bfa378f590396"

func init() {
	rand.Seed(time.Now().UnixNano())
}

type PostParse struct {
	Author     string `json"Omar Sohail"`
	Title      string `json:"title"`
	Content    string `json:"description"`
	ImgURL     string `json:"urlToImage"`
	Categories []int  `-`
}

type Form struct {
	Article []PostParse `json:"articles"`
}

var client *http.Client

func init() {
	client = &http.Client{Timeout: 10 * time.Second}
}

func GetJson(url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func ParseNews() (*[]PostParse, error) {
	var form Form
	err := GetJson(urlToParse, &form)
	if err != nil {
		return nil, err
	}

	return &form.Article, nil
}

func (s *Sqlite) CreateParsePosts() error {
	posts, err := ParseNews()

	if err != nil {
		return err
	}
	for _, post := range *posts {
		UserID := 1
		if post.Author != "" {
			UserID, err = s.createParseUser(post.Author)
			if err != nil {
				if errors.Is(err, models.ErrDuplicateEmail) {
					continue
				}
				return err
			}
		}

		postID, err := s.CreatePost(UserID, post.Title, post.Content, post.ImgURL)
		if err != nil {
			return err
		}
		categories := getRandomCategory(2)
		s.AddCategoryToPost(postID, categories)
	}
	return nil
}

func (s *Sqlite) createParseUser(name string) (int, error) {
	password := "nothingToSay"
	email := getEmail(name)
	user := models.UserSignupForm{
		Name:     name,
		Email:    email,
		Password: password,
	}
	id, err := s.CreateUserAndReturnID(user.FormToUser())
	return int(id), err
}

func getEmail(name string) string {
	for i, ch := range name {
		if ch == ' ' {
			return name[:i]
		}
	}
	return ""
}

func getRandomCategory(len int) []int {
	arr := make([]int, len)
	for i := range arr {
		arr[i] = rand.Intn(4) + 1
	}
	return arr
}
