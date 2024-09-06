package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage struct {
	config *Config
	db     *sql.DB

	userRepository    *UserRepository
	articleRepository *ArticleRepository
}

func (s *Storage) New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURI)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}

// Public Repo for Article
// (ex. s.Article().FindById())
func (s *Storage) Article() *ArticleRepository {
	if s.articleRepository != nil {
		return s.articleRepository
	}
	s.articleRepository = &ArticleRepository{
		storage: s,
	}
	return s.articleRepository
}

// Public Repo for User
// (ex. s.User().FindById())
func (s *Storage) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{
		storage: s,
	}
	return s.userRepository
}
