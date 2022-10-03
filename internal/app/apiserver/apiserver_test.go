package apiserver

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIServer_HandleHello(t *testing.T) {
	server := NewServer(NewDefaultConfig())
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/hello", nil)
	server.HandleHello().ServeHTTP(recorder, request)
	assert.Equal(t, "Hello!", recorder.Body.String())
}
