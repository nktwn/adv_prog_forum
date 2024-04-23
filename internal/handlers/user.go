package handlers

import (
	"errors"
	"forum/models"
	"forum/pkg/cookie"
	"forum/pkg/validator"
	"net/http"
)

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	methodResolver(w, r, h.loginGet, h.loginPost)
}

func (h *handler) loginGet(w http.ResponseWriter, r *http.Request) {
	data := h.app.NewTemplateData(r)
	data.Form = models.UserLoginForm{}
	h.app.Render(w, http.StatusOK, "login.html", data)
}

func (h *handler) loginPost(w http.ResponseWriter, r *http.Request) {
	form := models.UserLoginForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := h.app.NewTemplateData(r)
		data.Form = form
		h.app.Render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	session, err := h.service.Authenticate(form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) || errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("general", "Invalid email or password")
		} else if errors.Is(err, models.ErrNotActivated) {
			form.AddFieldError("general", "Please confirm your registration")
		} else {
			h.app.ServerError(w, err)
			return
		}

		data := h.app.NewTemplateData(r)
		data.Form = form
		h.app.Render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	cookie.SetSessionCookie(w, session.Token, session.ExpTime)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) signup(w http.ResponseWriter, r *http.Request) {
	methodResolver(w, r, h.signupGet, h.signupPost)
}

func (h *handler) signupGet(w http.ResponseWriter, r *http.Request) {
	data := h.app.NewTemplateData(r)
	data.Form = models.UserSignupForm{}
	h.app.Render(w, http.StatusOK, "signup.html", data)

}

func (h *handler) signupPost(w http.ResponseWriter, r *http.Request) {
	form := models.UserSignupForm{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := h.app.NewTemplateData(r)
		data.Form = form
		h.app.Render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}
	//
	user := form.FormToUser()
	err := h.service.CreateUser(&user)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := h.app.NewTemplateData(r)
			data.Form = form
			h.app.Render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			h.app.ServerError(w, err)
		}
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func (h *handler) logoutPost(w http.ResponseWriter, r *http.Request) {
	c := cookie.GetSessionCookie(r)
	if c != nil {
		h.service.DeleteSession(c.Value)
		cookie.ExpireSessionCookie(w)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) PostByUser(w http.ResponseWriter, r *http.Request) {
	c := cookie.GetSessionCookie(r)
	posts, err := h.service.GetAllPostByUser(c.Value)

	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	data := h.app.NewTemplateData(r)

	data.Posts = posts

	h.app.Render(w, http.StatusOK, "user_posts.html", data)

}

func (h *handler) userView(w http.ResponseWriter, r *http.Request) {
	data := h.app.NewTemplateData(r)
	c := cookie.GetSessionCookie(r)

	user, err := h.service.GetUserByToken(c.Value)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	data.User = user

	h.app.Render(w, http.StatusOK, "user.html", data)
}

func (h *handler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	methodResolver(w, r, h.UpdateUserPasswordGet, h.UpdateUserPasswordPost)
}

func (h *handler) UpdateUserPasswordGet(w http.ResponseWriter, r *http.Request) {
	data := h.app.NewTemplateData(r)
	data.Form = models.AccountPasswordUpdateForm{}
	h.app.Render(w, http.StatusOK, "password.html", data)
}

func (h *handler) UpdateUserPasswordPost(w http.ResponseWriter, r *http.Request) {
	form := models.AccountPasswordUpdateForm{
		CurrentPassword:         r.FormValue("currentPassword"),
		NewPassword:             r.FormValue("newPassword"),
		NewPasswordConfirmation: r.FormValue("newPasswordConfirmation"),
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwords do not match")
	if !form.Valid() {
		data := h.app.NewTemplateData(r)
		data.Form = form
		h.app.Render(w, http.StatusUnprocessableEntity, "password.html", data)
		return
	}

	c := cookie.GetSessionCookie(r)

	if err := h.service.UpdateUserPassword(c.Value, form.NewPassword); err != nil {
		h.app.ServerError(w, err)
	}

	http.Redirect(w, r, "/account", http.StatusSeeOther)
}

func (h *handler) activateAccount(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	err := h.service.ActivateUser(token)
	if err != nil {
		http.Error(w, "Unable to activate account", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login?activated=true", http.StatusSeeOther)
}
