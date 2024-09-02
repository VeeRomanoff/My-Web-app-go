package mywebapp

import (
	"github.com/VeeRomanoff/mywebapp/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
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

func (a *MyWebApp) Start() error {
	if err := a.configureLoggerField(); err != nil {
		return err
	}

	a.configureRouterField()

	return http.ListenAndServe(a.config.Port, a.router)
}

func (a *MyWebApp) configureLoggerField() error {
	level, err := logrus.ParseLevel(a.config.Logger)
	if err != nil {
		return err
	}
	a.logger.SetLevel(level)
	return nil
}

func (a *MyWebApp) configureRouterField() {
	a.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
}

func (a *MyWebApp) configureStorageField() error {
	storage := a.storage.New(a.config.StorageConfig)
	a.logger.Info("trying to open db storage")
	err := storage.Open()
	if err != nil {
		return err
	}
	a.storage = storage
	return nil
}
