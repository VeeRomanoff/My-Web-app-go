package mywebapp

import (
	"github.com/VeeRomanoff/mywebapp/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	prefix string = "/api/v1"
)

type MyWebApp struct {
	config  *Config
	router  *mux.Router
	logger  *logrus.Logger
	storage *storage.Storage
}

func NewMyWebApp(config *Config) *MyWebApp {
	return &MyWebApp{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (app *MyWebApp) Start() error {
	if err := app.configureLoggerField(); err != nil {
		return err
	}

	app.configureRouterField()

	if err := app.configureStorageField(); err != nil {
		return err
	}

	return http.ListenAndServe(app.config.Port, app.router)
}

func (app *MyWebApp) configureLoggerField() error {
	level, err := logrus.ParseLevel(app.config.Logger)
	if err != nil {
		return err
	}
	app.logger.SetLevel(level)
	return nil
}

func (app *MyWebApp) configureRouterField() {
	app.router.HandleFunc(prefix+"/articles", app.GetAllArticles).Methods("GET")
	//app.router.HandleFunc(prefix+"/articles/{id}", app.GetArticleById).Methods("GET")
	app.router.HandleFunc(prefix+"/articles", app.CreateArticle).Methods("POST")
	app.router.HandleFunc(prefix+"/user", app.CreateUser).Methods("POST")
	//app.router.HandleFunc(prefix+"/articles/{id}", app.UpdateArticleById).Methods("PUT")
	//app.router.HandleFunc(prefix+"/articles/{id}", app.DeleteArticleById).Methods("DELETE")
}

func (app *MyWebApp) configureStorageField() error {
	storage := app.storage.New(app.config.StorageConfig)
	app.logger.Info("trying to open db storage")
	err := storage.Open()
	if err != nil {
		app.logger.Info("couldnt open db storage")
		return err
	}
	app.logger.Info("opened db storage")
	app.storage = storage
	return nil
}
