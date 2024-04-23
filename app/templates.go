package app

import (
	"bytes"
	"fmt"
	"forum/models"
	"forum/pkg/cookie"
	"forum/ui"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"
)

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func sequence(start, end int) []int {
	var seq []int
	for i := start; i <= end; i++ {
		seq = append(seq, i)
	}
	return seq
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"sequence": sequence,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/*.layout.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}

func (app *Application) Render(w http.ResponseWriter, status int, page string, data *models.TemplateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.ServerError(w, err)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *Application) NewTemplateData(r *http.Request) *models.TemplateData {
	return &models.TemplateData{
		CurrentYear: time.Now().Year(),
		//Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		//CSRFToken:       nosurf.Token(r),
	}
}

func (app *Application) isAuthenticated(r *http.Request) bool {
	cookie := cookie.GetSessionCookie(r)
	return cookie != nil && cookie.Value != ""
}
