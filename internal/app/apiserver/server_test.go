package apiserver_test

import (
	"awesomeProject/internal/app/apiserver"
	"awesomeProject/internal/app/store"
	"awesomeProject/internal/app/store/teststore"
	"bytes"
	"encoding/json"
	sessions2 "github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_handleUsersCreate(t *testing.T) {
	testCases := []struct {
		key              string
		payload          interface{}
		expectedHttpCode int
	}{
		{
			key: "valid",
			payload: map[string]string{
				"email":    "abc@mail.com",
				"password": "1234567890",
			},
			expectedHttpCode: http.StatusCreated,
		},
		{
			key:              "invalid payload",
			payload:          "invalid payload format",
			expectedHttpCode: http.StatusBadRequest,
		},
		{
			key: "invalid data in correct payload",
			payload: map[string]string{
				"email":    "abc@mail.com",
				"password": "",
			},
			expectedHttpCode: http.StatusUnprocessableEntity,
		},
	}
	s := teststore.NewStore()
	server := apiserver.NewServer(s, sessions2.NewCookieStore([]byte("xxx")))

	for _, testCase := range testCases {
		t.Run(testCase.key, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			buf := &bytes.Buffer{}
			err := json.NewEncoder(buf).Encode(testCase.payload)
			if err != nil {
				t.Fatal(err)
			}
			request, _ := http.NewRequest(http.MethodPost, "/users", buf)
			server.ServeHTTP(recorder, request)
			assert.Equal(t, testCase.expectedHttpCode, recorder.Code)
		})
	}
}

func TestServer_handleSessionsCreate(t *testing.T) {
	userGen := store.TestUserHelper(t)
	s := teststore.NewStore()
	user := userGen()
	err := s.UserRepository().Create(user)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		key              string
		payload          interface{}
		expectedHttpCode int
	}{
		{
			key: "valid",
			payload: map[string]string{
				"email":    user.Email,
				"password": user.Password.Original,
			},
			expectedHttpCode: http.StatusOK,
		},
		{
			key:              "invalid payload",
			payload:          "invalid payload format",
			expectedHttpCode: http.StatusBadRequest,
		},
		{
			key: "invalid data in correct payload",
			payload: map[string]string{
				"email":    "abc@mail.com",
				"password": user.Password.Original + "xxx",
			},
			expectedHttpCode: http.StatusUnauthorized,
		},
	}
	server := apiserver.NewServer(s, sessions2.NewCookieStore([]byte("xxx")))

	for _, testCase := range testCases {
		t.Run(testCase.key, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			buf := &bytes.Buffer{}
			err := json.NewEncoder(buf).Encode(testCase.payload)
			if err != nil {
				t.Fatal(err)
			}
			request, _ := http.NewRequest(http.MethodPost, "/sessions", buf)
			server.ServeHTTP(recorder, request)
			assert.Equal(t, testCase.expectedHttpCode, recorder.Code)
		})
	}
}
