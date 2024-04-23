package handlers

import (
	"fmt"
	"net/http"
)

// func decorator(){

// }
func methodResolver(w http.ResponseWriter, r *http.Request, get, post func(w http.ResponseWriter, r *http.Request)) {
	fmt.Println(r.URL)
	switch r.Method {
	case http.MethodGet:
		get(w, r)
	case http.MethodPost:
		post(w, r)
	default:
		//error
	}
}
