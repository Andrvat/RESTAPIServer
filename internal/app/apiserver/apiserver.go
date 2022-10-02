package apiserver

import (
	"awesomeProject/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"time"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

func NewServer(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (receiver APIServer) Start() error {
	if err := receiver.configureLogLevel(); err != nil {
		return err
	}
	receiver.configureRouter()
	if err := receiver.configureStore(); err != nil {
		return err
	}
	receiver.logger.Info("Start API server at " + time.Now().Format(time.RFC850))
	err := http.ListenAndServe(receiver.config.BindAddr, receiver.router)
	return err
}

func (receiver APIServer) configureLogLevel() error {
	level, err := logrus.ParseLevel(receiver.config.LogLevel)
	if err != nil {
		return err
	}
	receiver.logger.SetLevel(level)
	receiver.logger.Info("Set new log level: " + receiver.config.LogLevel)
	return nil
}

func (receiver APIServer) configureRouter() {
	receiver.router.HandleFunc("/hello", receiver.HandleHello())

}

func (receiver APIServer) HandleHello() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if _, err := io.WriteString(writer, "Hello!"); err != nil {
			log.Fatal(err)
		}
	}
}

func (receiver APIServer) configureStore() error {
	st := store.NewStore(receiver.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	receiver.store = st
	return nil
}
