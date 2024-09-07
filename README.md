# My-Web-app-go

## Работа с моделью

### Шаг 0. Откат миграции
Для выполнения отката ```migrate -path migrations -database "postgres://localhost:5432/restapi?sslmode=disable&user=postgres&password=postgres" down```


### Шаг 1. Новая миграция
Заходим в файл ```migrations/.....up.sql```
```
CREATE TABLE users (
    id bigserial not null primary key,
    login varchar not null unique,
    password varchar not null
);

CREATE TABLE articles (
    id bigserial not null primary key,
    title varchar not null unique,
    author varchar not null,
    content varchar not null
);
```

Выполним команду ```migrate -path migrations -database "postgres://localhost:5432/restapi?sslmode=disable&user=postgres&password=postgres" down```

### Шаг 2. Определим модели
Для того, чтобы определить модели ```internal/app/models/``` 2 модуля:
* user.go
* article.go

```
//user.go
package models

//User model defeniton
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

```

```
//article.go
package models

//Article model defenition
type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

```

### Шаг 3. Определение "репозиториев"
Работать с моделями будем через репозитории. Для этого инициализируем 2 файла:
* ```storage/userrepository.go```
* ```storage/articlerepository.go```

```
//articlerepository.go
package storage

//Instance of Article repository (model interface)
type ArticleRepository struct {
    storage *Storage
}

```

Аналогично для юзера.

### Шаг 4. Выделение публичного доступа к репозиторию
Хотим, чтобы наше приложение общалось с моделями через репозитории (которые будут содержать необходимый набор метод для взаимодействия с бд). Нам необходимо определить 2 метода у хранилища , которые будут предоставлять публичные репозитории:
```
//storage.go

//Instance of storage
type Storage struct {
	config *Config
	// DataBase FileDescriptor
	db *sql.DB
	//Subfield for repo interfacing (model user)
	userRepository *UserRepository
	//Subfield for repo interfaceing (model article)
	articleRepository *ArticleRepository
}

....

//Public Repo for Article
func (s *Storage) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{
		storage: s,
	}
	return nil
}

//Public Repo for User
func (s *Storage) Article() *ArticleRepository {
	if s.articleRepository != nil {
		return s.articleRepository
	}
	s.articleRepository = &ArticleRepository{
		storage: s,
	}
	return nil
}

```

### Шаг 5. Что будет уметь делать UserRepo?
* Сохранять нового пользователя в бд (INSERT user'a или Create)
* Для аутентификации нужен функционал поиска пользователя по ```login```
* Выдача всех пользователей из бд
```
package storage

import (
	"fmt"
	"log"

	"github.com/vlasove/go2/7.ServerAndDB2/internal/app/models"
)

//Instance of User repository (model interface)
type UserRepository struct {
	storage *Storage
}

var (
	tableUser string = "users"
)

//Create User in db
func (ur *UserRepository) Create(u *models.User) (*models.User, error) {
	query := fmt.Sprintf("INSERT INTO %s (login, password) VALUES ($1, $2) RETURNING id", tableUser)
	if err := ur.storage.db.QueryRow(query, u.Login, u.Password).Scan(&u.ID); err != nil {
		return nil, err
	}
	return u, nil
}

//Find user by login
func (ur *UserRepository) FindByLogin(login string) (*models.User, bool, error) {
	users, err := ur.SelectAll()
	var founded bool
	if err != nil {
		return nil, founded, err
	}
	var userFinded *models.User
	for _, u := range users {
		if u.Login == login {
			userFinded = u
			founded = true
			break
		}
	}
	return userFinded, founded, nil
}

//Select all users in db
func (ur *UserRepository) SelectAll() ([]*models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableUser)
	rows, err := ur.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//Подготовим, куда будем читать
	users := make([]*models.User, 0)
	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.ID, &u.Login, &u.Password)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, &u)
	}
	return users, nil
}

```

### Шаг 6. Что нужно от ArticleRepo?
* Уметь доавлять статью в бд
* Уметь удалять по id
* Получать все статьи
* Получать статью по id
* Обновлять (дома)
```
articlerepository.go
package storage

import (
	"fmt"
	"log"

	"github.com/vlasove/go2/7.ServerAndDB2/internal/app/models"
)

//Instance of Article repository (model interface)
type ArticleRepository struct {
	storage *Storage
}

var (
	tableArticle string = "articles"
)

//Добавить статью в бд
func (ar *ArticleRepository) Create(a *models.Article) (*models.Article, error) {
	query := fmt.Sprintf("INSERT INTO %s (title, author, content) VALUES ($1, $2, $3) RETURNING id", tableArticle)
	if err := ar.storage.db.QueryRow(query, a.Title, a.Author, a.Content).Scan(&a.ID); err != nil {
		return nil, err
	}

	return a, nil

}

//Удалять статью по id
func (ar *ArticleRepository) DeleteById(id int) (*models.Article, error) {
	article, ok, err := ar.FindArticleById(id)
	if err != nil {
		return nil, err
	}
	if ok {
		query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", tableArticle)
		_, err := ar.storage.db.Exec(query, id)
		if err != nil {
			return nil, err
		}
	}
	return article, nil
}

//Получаем статью по id
func (ar *ArticleRepository) FindArticleById(id int) (*models.Article, bool, error) {
	articles, err := ar.SelectAll()
	var founded bool
	if err != nil {
		return nil, founded, err
	}
	var articleFinded *models.Article
	for _, a := range articles {
		if a.ID == id {
			articleFinded = a
			founded = true
			break
		}
	}
	return articleFinded, founded, nil
}

//Получим все статьи в бд
func (ar *ArticleRepository) SelectAll() ([]*models.Article, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableArticle)
	rows, err := ar.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//Подготовим, куда будем читать
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

```

### Шаг 7. Описание маршрутизатора для данного проекта.
Зайдем в ```api```
```
//Пытаемся отконфигурировать маршрутизатор (а конкретнее поле router API)
func (a *API) configreRouterField() {
	a.router.HandleFunc(prefix+"/articles", a.GetAllArticles).Methods("GET")
	a.router.HandleFunc(prefix+"/articles/{id}", a.GetArticleById).Methods("GET")
	a.router.HandleFunc(prefix+"/articles/{id}", a.DeleteArticleById).Methods("DELETE")
	a.router.HandleFunc(prefix+"/articles", a.PostArticle).Methods("POST")
	a.router.HandleFunc(prefix+"/user/register", a.PostUserRegister).Methods("POST")

}
```

Создадим файл ```internal/app/api/handlers.go```
```
```

## Реализация обработчиков

Из-за того, что пока у ```users``` всего один обработчик, будет держать все handlers в одном месте :
```
internal/app/api/handlers.go
```

Внутри определим 2 сущности:
```
package api

import "net/http"

//Вспомогательная структура для формирования сообщений
type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

```

### Шаг 1. Реализация обработчика GetAllArticles
```
//Возвращает все статьи из бд на данный момент
func (api *API) GetAllArticles(writer http.ResponseWriter, req *http.Request) {
	//Инициализируем хедеры
	initHeaders(writer)
	//Логируем момент начало обработки запроса
	api.logger.Info("Get All Artiles GET /api/v1/articles")
	//Пытаемся что-то получить от бд
	articles, err := api.storage.Article().SelectAll()
	if err != nil {
		//Что делаем, если была ошибка на этапе подключения?
		api.logger.Info("Error while Articles.SelectAll : ", err)
		msg := Message{
			StatusCode: 501,
			Message:    "We have some troubles to accessing database. Try again later",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(articles)
}
```

### Шаг 2. Реализация PostArticle
```
`````