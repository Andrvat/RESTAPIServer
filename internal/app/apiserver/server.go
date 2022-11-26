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

	userContextKey contextKey = iota
	requestIdContextKey
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

	s.router.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)

	s.router.HandleFunc("/sign-up", s.handleUserCreate()).Methods("POST")
	s.router.HandleFunc("/sign-in", s.handleSessionCreate()).Methods("POST")

	privateSubRouter := s.router.PathPrefix("/authorized").Subrouter()
	privateSubRouter.Use(s.AuthenticateUser)
	privateSubRouter.HandleFunc("/whoami", s.handleWhoAmI()).Methods("GET")
	privateSubRouter.HandleFunc("/users", s.handleUsersGetAll()).Methods("GET")
	privateSubRouter.HandleFunc("/update", s.handleUserUpdate()).Methods("POST", "PUT")
	privateSubRouter.HandleFunc("/delete", s.handleUserDelete()).Methods("DELETE")
	privateSubRouter.HandleFunc("/logout", s.handleSessionLogout()).Methods("PUT")

}

func (s *Server) SetRequestId(nextFunc http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.New().String()
		w.Header().Set("X-Request-ID", requestId)
		newContext := context.WithValue(r.Context(), requestIdContextKey, requestId)
		nextFunc.ServeHTTP(w, r.WithContext(newContext))
	})
}

func (s *Server) AuthenticateUser(nextFunc http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := (*s.sessions).Get(r, SessionName)
		if err != nil {
			s.handleError(w, r, http.StatusInternalServerError, err)
			return
		}

		id, exist := session.Values[UserIdSessionKey]
		if !exist {
			s.handleError(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		user, err := (*s.store).UserRepository().FindById(id.(int))
		if err != nil {
			s.handleError(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		newContext := context.WithValue(r.Context(), userContextKey, user)
		nextFunc.ServeHTTP(w, r.WithContext(newContext))
	})
}

func (s *Server) LogRequest(nextFunc http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		localLogger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(requestIdContextKey),
		})
		localLogger.Infof("Started %s %s", r.Method, r.RequestURI)

		startTime := time.Now()

		responseWriter := &ResponseWriter{w, http.StatusOK}
		nextFunc.ServeHTTP(responseWriter, r)

		status := fmt.Sprintf("status %s %v", http.StatusText(responseWriter.statusCode), responseWriter.statusCode)
		localLogger.Infof("Completed in %v wtih %s", time.Now().Sub(startTime), status)

	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// @Summary WhoAmI
// @Tags common
// @Description Get general info about yourself after authorization
// @ID user-whoami
// @Accept json
// @Produce json
// @Success 200 {integer} 1
// @Failure 401 {object} error
// @Router /authorized/whoami [get]
func (s *Server) handleWhoAmI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		maybeUser := r.Context().Value(userContextKey)
		if maybeUser == nil {
			s.handleError(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		user := maybeUser.(*model.User)
		s.respond(w, r, http.StatusOK, user)
	}
}

// @Summary CreateUser
// @Tags registration
// @Description Create new user and store in database
// @ID user-create
// @Accept json
// @Produce json
// @Param input body SignRequest true "Info about email and password"
// @Success 201 {integer} 1
// @Failure 400 {object} error
// @Failure 422 {object} error
// @Router /sign-up [post]
func (s *Server) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userMeta := &SignRequest{}
		if err := json.NewDecoder(r.Body).Decode(userMeta); err != nil {
			s.handleError(w, r, http.StatusBadRequest, err)
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
			s.handleError(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusCreated, model.Sanitized(user))
	}
}

// @Summary CreateSession
// @Tags authentication
// @Description Create new session for existing user
// @ID session-create
// @Accept json
// @Produce json
// @Param input body SignRequest true "Info about email and password"
// @Success 200
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /sign-in [post]
func (s *Server) handleSessionCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userMeta := &SignRequest{}
		if err := json.NewDecoder(r.Body).Decode(userMeta); err != nil {
			s.handleError(w, r, http.StatusBadRequest, err)
			return
		}
		user, err := (*s.store).UserRepository().FindByEmail(userMeta.Email)
		if err != nil || !user.HasSamePassword(userMeta.Password) {
			s.handleError(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := (*s.sessions).Get(r, SessionName)
		if err != nil {
			s.handleError(w, r, http.StatusInternalServerError, errIncorrectEmailOrPassword)
			return
		}

		session.Values[UserIdSessionKey] = user.Id
		err = (*s.sessions).Save(r, w, session)
		if err != nil {
			s.handleError(w, r, http.StatusInternalServerError, errIncorrectEmailOrPassword)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

// @Summary SessionLogout
// @Tags common
// @Description Log out from current session after authorization
// @ID session-logout
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /authorized/logout [put]
func (s *Server) handleSessionLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := (*s.sessions).Get(r, SessionName)
		if err != nil {
			s.handleError(w, r, http.StatusInternalServerError, err)
			return
		}

		delete(session.Values, UserIdSessionKey)

		err = (*s.sessions).Save(r, w, session)
		if err != nil {
			s.handleError(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

// @Summary AllUsers
// @Tags common
// @Description Get all existing users
// @ID users-get-all
// @Accept json
// @Produce json
// @Success 200
// @Failure 500 {object} error
// @Router /authorized/users [get]
func (s *Server) handleUsersGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := (*s.store).UserRepository().AllUsers()
		if err != nil {
			s.handleError(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, users)
	}
}

// @Summary UpdateUser
// @Tags common
// @Description Update yourself after authorization
// @ID users-update
// @Accept json
// @Produce json
// @Param input body SignRequest true "New email or password"
// @Success 200
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Router /authorized/update [post]
// @Router /authorized/update [put]
func (s *Server) handleUserUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		maybeContextUser := r.Context().Value(userContextKey)
		if maybeContextUser == nil {
			s.handleError(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		contextUser := maybeContextUser.(*model.User)

		userMeta := &SignRequest{}
		if err := json.NewDecoder(r.Body).Decode(userMeta); err != nil {
			s.handleError(w, r, http.StatusBadRequest, err)
			return
		}

		if userMeta.Email == "" && userMeta.Password == "" {
			s.handleError(w, r, http.StatusBadRequest, errNonEmptyBodyRequired)
			return
		}

		finalEmail := contextUser.Email
		if userMeta.Email != "" {
			finalEmail = userMeta.Email
		}
		finalPassword := contextUser.Password
		if userMeta.Password != "" {
			finalPassword = &model.Password{
				Original: userMeta.Password,
			}
		}
		user := &model.User{
			Id:       contextUser.Id,
			Email:    finalEmail,
			Password: finalPassword,
		}

		err := (*s.store).UserRepository().Update(user)
		if err != nil {
			s.handleError(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

// @Summary DeleteUser
// @Tags common
// @Description Delete yourself after authorization
// @ID users-delete
// @Accept json
// @Produce json
// @Success 200
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /authorized/delete [delete]
func (s *Server) handleUserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		maybeContextUser := r.Context().Value(userContextKey)
		if maybeContextUser == nil {
			s.handleError(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		contextUser := maybeContextUser.(*model.User)
		err := (*s.store).UserRepository().Delete(contextUser)
		if err != nil {
			s.handleError(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *Server) handleError(w http.ResponseWriter, r *http.Request, status int, err error) {
	s.respond(w, r, status, map[string]string{"error": err.Error()})
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return
		}
	}
}
