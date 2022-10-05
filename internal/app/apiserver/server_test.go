package apiserver_test

import (
	"awesomeProject/internal/app/apiserver"
	"awesomeProject/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleUsersCreate(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/users", nil)
	store := teststore.NewStore()
	server := apiserver.NewServer(store)
	server.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)

}
