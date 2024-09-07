package storage

import (
	"fmt"
	"github.com/VeeRomanoff/mywebapp/internal/mywebapp/models"
	"log"
)

// ArticleRepository - Instance of Article repository (model interface)
type ArticleRepository struct {
	storage *Storage
}

var (
	articleTable string = "articles"
)

func (ar *ArticleRepository) Create(a *models.Article) (*models.Article, error) {
	query := fmt.Sprintf("INSERT INTO %s (title, author, content) VALUES ($1, $2, $3) RETURNING id", articleTable)
	err := ar.storage.db.QueryRow(query, a.Title, a.Author, a.Content).Scan(&a.ID)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// FindArticleByTitle bool in case if resource doesn't exist but connection is opened (this case is not considered to be an error)
func (ar *ArticleRepository) FindArticleById(id int) (*models.Article, bool, error) {
	articles, err := ar.SelectAll()
	var found bool
	if err != nil {
		return nil, false, err
	}
	var articleFound *models.Article
	for _, a := range articles {
		if a.ID == id {
			articleFound = a
			found = true
			break
		}
	}

	return articleFound, found, nil
}

func (ar *ArticleRepository) SelectAll() ([]*models.Article, error) {
	query := fmt.Sprintf("SELECT * FROM %s", articleTable)
	rows, err := ar.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Prepare where to read articles
	articles := make([]*models.Article, 0)
	for rows.Next() {
		a := models.Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Author, &a.Content)
		if err != nil {
			log.Println(err)
			continue
		}
		articles = append(articles, &a)
	}

	return articles, nil
}

func (ar *ArticleRepository) UpdateArticleById(id int, a *models.Article) (*models.Article, error) {
	_, exist, err := ar.FindArticleById(id)
	if err != nil {
		return nil, err
	}
	if exist {
		query := fmt.Sprintf("UPDATE %s SET title = $1, content = $2 WHERE id = $3", articleTable)
		_, err = ar.storage.db.Exec(query, a.Title, a.Content, a.ID)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (ar *ArticleRepository) DeleteById(id int) (*models.Article, error) {
	article, ok, err := ar.FindArticleById(id)
	if err != nil {
		return nil, err
	}
	if ok {
		query := fmt.Sprintf("DELETE FROM %s WHERE ID = $1", articleTable)
		_, err := ar.storage.db.Exec(query, article.ID)
		if err != nil {
			return nil, err
		}
	}
	return article, nil
}
