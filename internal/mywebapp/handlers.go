package mywebapp

import (
	"encoding/json"
	"github.com/VeeRomanoff/mywebapp/internal/mywebapp/models"
	"net/http"
)

// Message это вспомогательная структура для сообщений о результате
type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

// GetAllArticles returns all articles from the database
func (app *MyWebApp) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	app.logger.Info("Get All Articles GET /api/v1/articles...")
	// articles - nil???
	articles, err := app.storage.Article().SelectAll()
	// The problem that can occur while accessing "SelectAll" is a database problem. Let's handle it
	if err != nil {
		app.logger.Info("Error getting all articles: ", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles accessing database. Try again later.",
			IsError:    true,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(articles)
}

func (app *MyWebApp) CreateArticle(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	app.logger.Info("Post Article POST /api/v1/article...")
	// @RequestBody
	var article models.Article
	// json from client might be invalid
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		app.logger.Info("Invalid json received from client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided JSON is invalid",
			IsError:    true,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}
	a, err := app.storage.Article().Create(&article)
	// The problem that can occur while accessing "Create" is a database problem. Let's handle it
	if err != nil {
		app.logger.Info("Error creating the article: ", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles accessing database. Try again later.",
			IsError:    true,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
}

func (app *MyWebApp) CreateUser(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	app.logger.Info("Post User POST /api/v1/user...")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.logger.Info("Invalid json received from client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided JSON is invalid",
			IsError:    true,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}
	u, err := app.storage.User().Create(&user)
	if err != nil {
		app.logger.Info("Error creating the user: ", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles accessing database. Try again later.",
			IsError:    true,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}
