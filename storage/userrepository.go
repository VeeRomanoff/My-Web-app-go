package storage

import (
	"fmt"
	"github.com/VeeRomanoff/mywebapp/internal/mywebapp/models"
	"log"
)

var (
	usersTable string = "users"
)

// UserRepository - Instance of User repository (model interface)
type UserRepository struct {
	storage *Storage
}

func (ur *UserRepository) SelectAll() ([]*models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s", usersTable)
	rows, err := ur.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*models.User, 0)

	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.Login, u.Password)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, &u)
	}
	return users, nil
}

func (ur *UserRepository) Create(u *models.User) (*models.User, error) {
	if ur.storage.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	query := fmt.Sprintf("INSERT INTO %s (login, password) VALUES ($1, $2) RETURNING id", usersTable)
	err := ur.storage.db.QueryRow(query, u.Login, u.Password).Scan(&u.ID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *UserRepository) FindUserById(id int) (*models.User, bool, error) {
	users, err := ur.SelectAll()
	var found bool
	if err != nil {
		return nil, false, err
	}
	var userFound *models.User
	for _, u := range users {
		if u.ID == id {
			userFound = u
			found = true
			break
		}
	}
	return userFound, found, nil
}

func (ur *UserRepository) FindUserByLogin(login string) (*models.User, bool, error) {
	users, err := ur.SelectAll()
	var found bool
	if err != nil {
		return nil, false, err
	}
	var userFound *models.User
	for _, u := range users {
		if u.Login == login {
			userFound = u
			found = true
			break
		}
	}
	return userFound, found, nil
}

func (ur *UserRepository) DeleteUserById(id int) (*models.User, error) {
	user, exist, err := ur.FindUserById(id)
	if err != nil {
		return nil, err
	}
	if exist {
		query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", usersTable)
		_, err := ur.storage.db.Exec(query, id)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (ur *UserRepository) UpdateUserById(id int, u *models.User) (*models.User, error) {
	_, exist, err := ur.FindUserById(id)
	if err != nil {
		return nil, err
	}
	if exist {
		query := fmt.Sprintf("UPDATE %s SET login = $1, password = $2 WHERE id = $3", usersTable)
		_, err = ur.storage.db.Exec(query, u.Login, u.Password, u.ID)
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}
