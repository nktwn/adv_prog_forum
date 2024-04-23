package handlers

import (
	"forum/pkg/cookie"
	"net/http"
	"strconv"
)

func (h *handler) adminDashboard(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.app.ClientError(w, http.StatusForbidden)
		return
	}

	users, err := h.service.GetAllUsers()
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	data := h.app.NewTemplateData(r)
	data.Users = users
	h.app.Render(w, http.StatusOK, "admin_dashboard.html", data)
}

func (h *handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		h.app.ClientError(w, http.StatusForbidden)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.app.ClientError(w, http.StatusBadRequest)
		return
	}

	err = h.service.DeleteUser(userID)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

//	func (h *handler) isAdmin(r *http.Request) bool {
//		// логика определения администратора на будуще ежже
//		cookie := cookie.GetSessionCookie(r)
//		user, err := h.service.GetUserByToken(cookie.Value)
//		if err != nil {
//			return false
//		}
//
//		return user.IsAdmin
//	}

func (h *handler) isAdmin(r *http.Request) bool {
	cookie := cookie.GetSessionCookie(r)
	user, err := h.service.GetUserByToken(cookie.Value)
	if err != nil {
		return false
	}
	if user.Email == "prostok06@gmail.com" {
		return true
	}

	return false
}
