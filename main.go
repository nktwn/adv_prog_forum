package main

import (
	"fmt"
	"forum/internal/repo/sqlite"
	"os"
)

func main() {
	posts, err := sqlite.ParseNews()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, post := range *posts {
		fmt.Println(post.Author)
	}
}
