package apiserver

import (
	"awesomeProject/internal/app/model"
	"awesomeProject/internal/app/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	sessionName = "xxx"
)

type Server struct {
	logger   *logrus.Logger
	router   *mux.Router
	store    *store.Store
	sessions *sessions.Store
}

func NewServer(store store.Store, sessions sessions.Store) *Server {
	s := &Server{
		store:    &store,
		router:   mux.NewRouter(),
		logger:   logrus.New(),
		sessions: &sessions,
	}
	s.configureRouter()

	return s
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleUsersCreate() http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		userMeta := &Request{}
		if err := json.NewDecoder(request.Body).Decode(userMeta); err != nil {
			s.handleError(writer, request, http.StatusBadRequest, err)
			return
		}
		user := &model.User{
			Email: userMeta.Email,
			Password: &model.Password{
				Original: userMeta.Password,
			},
		}
		err := (*s.store).UserRepository().Create(user)
		if err != nil {
			s.handleError(writer, request, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(writer, request, http.StatusCreated, model.Sanitized(user))
	}
}

func (s *Server) handleSessionsCreate() http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		userMeta := &Request{}
		if err := json.NewDecoder(request.Body).Decode(userMeta); err != nil {
			s.handleError(writer, request, http.StatusBadRequest, err)
			return
		}
		user, err := (*s.store).UserRepository().FindByEmail(userMeta.Email)
		if err != nil || !user.HasSamePawword(userMeta.Password) {
			s.handleError(writer, request, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := (*s.sessions).Get(request, sessionName)
		if err != nil {
			s.handleError(writer, request, http.StatusInternalServerError, errIncorrectEmailOrPassword)
			return
		}

		session.Values["user_id"] = user.Id
		err = (*s.sessions).Save(request, writer, session)
		if err != nil {
			s.handleError(writer, request, http.StatusInternalServerError, errIncorrectEmailOrPassword)
			return
		}

		s.respond(writer, request, http.StatusOK, nil)
	}
}

func (s *Server) handleError(writer http.ResponseWriter, request *http.Request, status int, err error) {
	s.respond(writer, request, status, map[string]string{"error": err.Error()})
}

func (s *Server) respond(writer http.ResponseWriter, request *http.Request, status int, data interface{}) {
	writer.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(writer).Encode(data)
		if err != nil {
			return
		}
	}
}
