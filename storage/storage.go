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
func (s *Storage) Article() *ArticleRepository {
	if s.articleRepository != nil {
		return s.articleRepository
	}
	s.articleRepository = &ArticleRepository{
		storage: s,
	}
	return nil
}

// Public Repo for User
func (s *Storage) User() *UserRepository {
	if s.articleRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{
		storage: s,
	}
	return nil
}
