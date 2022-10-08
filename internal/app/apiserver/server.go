package apiserver

import (
	_ "awesomeProject/docs"
	"awesomeProject/internal/app/model"
	"awesomeProject/internal/app/store"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	"net/http"
	"time"
)

const (
	SessionName      = "xxx"
	UserIdSessionKey = "user_id"

	contextKeyUser contextKey = iota
	contextKeyRequestId
)

type contextKey int8

type SignRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
	s.router.Use(s.SetRequestId)
	s.router.Use(s.LogRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	privateSubRouter := s.router.PathPrefix("/private").Subrouter()
	privateSubRouter.Use(s.AuthenticateUser)
	privateSubRouter.HandleFunc("/whoami", s.handleWhoAmI()).Methods("GET")

	s.router.PathPrefix("/documentation").Handler(httpSwagger.WrapHandler)
}

func (s *Server) SetRequestId(nextFunc http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestId := uuid.New().String()
		writer.Header().Set("X-Request-ID", requestId)
		newContext := context.WithValue(request.Context(), contextKeyRequestId, requestId)
		nextFunc.ServeHTTP(writer, request.WithContext(newContext))
	})
}

func (s *Server) AuthenticateUser(nextFunc http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session, err := (*s.sessions).Get(request, SessionName)
		if err != nil {
			s.handleError(writer, request, http.StatusInternalServerError, err)
			return
		}

		id, exist := session.Values[UserIdSessionKey]
		if !exist {
			s.handleError(writer, request, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		user, err := (*s.store).UserRepository().FindById(id.(int))
		if err != nil {
			s.handleError(writer, request, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		newContext := context.WithValue(request.Context(), contextKeyUser, user)
		nextFunc.ServeHTTP(writer, request.WithContext(newContext))
	})
}

func (s *Server) LogRequest(nextFunc http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		localLogger := s.logger.WithFields(logrus.Fields{
			"remote_addr": request.RemoteAddr,
			"request_id":  request.Context().Value(contextKeyRequestId),
		})
		localLogger.Infof("Started %s %s", request.Method, request.RequestURI)

		startTime := time.Now()

		responseWriter := &ResponseWriter{writer, http.StatusOK}
		nextFunc.ServeHTTP(responseWriter, request)

		status := fmt.Sprintf("status %s %v", http.StatusText(responseWriter.statusCode), responseWriter.statusCode)
		localLogger.Infof("Completed in %v wtih %s", time.Now().Sub(startTime), status)

	})
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleWhoAmI() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		user := request.Context().Value(contextKeyUser).(*model.User)
		s.respond(writer, request, http.StatusOK, user)
	}
}

// @Summary CreateUser
// @Tags authentication
// @Description Create new user and store in database
// @ID user-create
// @Accept json
// @Produce json
// @Param input body SignRequest true "Info about email and password"
// @Success 201 {integer} 1
// @Failure 400 {object} error
// @Failure 422 {object} error
// @Router /users [post]
func (s *Server) handleUsersCreate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userMeta := &SignRequest{}
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
	return func(writer http.ResponseWriter, request *http.Request) {
		userMeta := &SignRequest{}
		if err := json.NewDecoder(request.Body).Decode(userMeta); err != nil {
			s.handleError(writer, request, http.StatusBadRequest, err)
			return
		}
		user, err := (*s.store).UserRepository().FindByEmail(userMeta.Email)
		if err != nil || !user.HasSamePawword(userMeta.Password) {
			s.handleError(writer, request, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := (*s.sessions).Get(request, SessionName)
		if err != nil {
			s.handleError(writer, request, http.StatusInternalServerError, errIncorrectEmailOrPassword)
			return
		}

		session.Values[UserIdSessionKey] = user.Id
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
