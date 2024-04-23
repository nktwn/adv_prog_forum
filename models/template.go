package models

type TemplateData struct {
	CurrentYear     int
	Post            *Post
	Posts           *[]Post
	Categories      []string
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	User            *User
	NumberOfPage    int
	CurrentPage     int
	Users           []User
}
