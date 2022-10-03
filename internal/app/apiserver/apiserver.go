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

func (s *APIServer) Start() error {
	if err := s.configureLogLevel(); err != nil {
		return err
	}
	s.configureRouter()
	if err := s.configureStore(); err != nil {
		return err
	}
	s.logger.Info("Start API server at " + time.Now().Format(time.RFC850))
	err := http.ListenAndServe(s.config.BindAddr, s.router)
	return err
}

func (s *APIServer) configureLogLevel() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	s.logger.Info("Set new log level: " + s.config.LogLevel)
	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/hello", s.HandleHello())

}

func (s *APIServer) HandleHello() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if _, err := io.WriteString(writer, "Hello!"); err != nil {
			log.Fatal(err)
		}
	}
}

func (s *APIServer) configureStore() error {
	st := store.NewStore(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}
