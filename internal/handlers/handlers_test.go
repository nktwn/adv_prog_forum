package handlers

import (
	"forum/app"
	"forum/internal/repo"
	"forum/internal/service"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func newTestServer(t *testing.T) *httptest.Server {
	infoLog := log.New(os.Stdout, "\u001b[32mINFO\t\u001b[0m", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "\u001b[31mERROR\t\u001b[0m", log.Ldate|log.Ltime|log.Lshortfile)

	// cfg := config.MustLoad()

	tc, err := app.NewTemplateCache()

	if err != nil {
		errLog.Fatal(err)
	}

	app := app.New(infoLog, errLog, tc)

	r, err := repo.New("../../data/storage.db")
	if err != nil {
		log.Fatal(err)
	}

	s := service.New(r)

	h := New(s, app)

	srv := httptest.NewServer(h.Routes())
	return srv
}

func TestPostView(t *testing.T) {

	ts := newTestServer(t)

	defer ts.Close()

	tests := []struct {
		name     string
		url      string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			url:      "/post/1",
			wantCode: http.StatusOK,
		},
		{
			name:     "Non-existent ID",
			url:      "/post/100",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			url:      "/post/-1",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Decimal ID",
			url:      "/post/1.77",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "String ID",
			url:      "/post/bruh",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Empty ID",
			url:      "/sdfadsf",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tt.url) // Use ts.URL
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("got %d; want %d", resp.StatusCode, tt.wantCode)
			}

			if tt.wantBody != "" {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(string(body), tt.wantBody) {
					t.Errorf("unexpected body: got %s; want %s", body, tt.wantBody)
				}
			}
		})
	}
}
func TestHomeGet(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	tests := []struct {
		name       string
		url        string
		wantCode   int
		wantInBody string
	}{

		{
			name:     "No page",
			url:      "/fdsafads",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "pagination Page",
			url:      "/?page=4",
			wantCode: http.StatusOK,
		},
		{
			name:     "Wrong pagination Page",
			url:      "/?page=fdsafdas",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Wrong pagination Page",
			url:      "/?page=0.32",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Wrong pagination Page",
			url:      "/?page=-43",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tc.url)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantCode {
				t.Errorf("got %d; want %d", resp.StatusCode, tc.wantCode)
			}

			if tc.wantInBody != "" {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(string(body), tc.wantInBody) {
					t.Errorf("unexpected body: got %s; want %s", body, tc.wantInBody)
				}
			}
		})
	}
}
func TestHomePost(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	infoLog := log.New(os.Stdout, "\u001b[32mINFO\t\u001b[0m", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "\u001b[31mERROR\t\u001b[0m", log.Ldate|log.Ltime|log.Lshortfile)

	// cfg := config.MustLoad()

	tc, err := app.NewTemplateCache()

	if err != nil {
		errLog.Fatal(err)
	}

	app := app.New(infoLog, errLog, tc)

	r, err := repo.New("../../data/storage.db")
	if err != nil {
		log.Fatal(err)
	}

	s := service.New(r)

	h := New(s, app)

	tests := []struct {
		name       string
		categories []string
		wantCode   int
	}{
		{
			name:       "Valid Submission",
			categories: []string{"1", "2"},
			wantCode:   http.StatusOK,
		},
		{
			name:       "Invalid Submission",
			categories: []string{"1000"},
			wantCode:   http.StatusBadRequest,
		},
		{
			name:       "Invalid Submission",
			categories: []string{"dsafasdf"},
			wantCode:   http.StatusBadRequest,
		},
		{
			name:       "Invalid Submission",
			categories: []string{},
			wantCode:   http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{}
			for _, cat := range tc.categories {
				form.Add("categories", cat)
			}

			req, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()

			h.homePost(rr, req)

			if rr.Code != tc.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tc.wantCode)
			}

		})
	}
}
