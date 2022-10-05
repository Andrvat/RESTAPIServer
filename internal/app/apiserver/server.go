package apiserver

import (
	"awesomeProject/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

func NewServer(store store.Store) *Server {
	s := &Server{
		store:  &store,
		router: mux.NewRouter(),
		logger: logrus.New(),
	}
	s.configureRouter()

	return s
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/users", s.HandleUsersCreate()).Methods("POST")
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) HandleUsersCreate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
