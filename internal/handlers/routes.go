package handlers

import (
	"forum/ui"
	"net/http"
	"path/filepath"
)

func (h *handler) Routes() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.FS(ui.Files)})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", fileServer)

	mux.HandleFunc("/", h.home)
	mux.HandleFunc("/post/", h.postView)
	mux.HandleFunc("/post/create", h.postCreate)
	mux.HandleFunc("/login", h.login)
	mux.HandleFunc("/signup", h.signup)
	mux.HandleFunc("/logout", h.logoutPost)

	mux.HandleFunc("/account/view", h.PostByUser)
	mux.HandleFunc("/account/password", h.UpdateUserPassword)
	mux.HandleFunc("/account", h.userView)
	mux.HandleFunc("/admin/dashboard", h.adminDashboard)
	mux.HandleFunc("/admin/delete", h.deleteUser)
	mux.HandleFunc("/activate", h.activateAccount)

	return mux
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}

	}

	return f, nil
}
