package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func Test_handleInsertTest(t *testing.T) {
	srv := Server{
		// DB: *MockDB,
	}
	srv.Routes()

	dbtestURL, err := srv.Router.Get(dbTestRouteName).URL()
	if err != nil {
		t.Errorf(err.Error())
	}

	req := httptest.NewRequest(http.MethodGet, dbtestURL.Path, nil)
	w := httptest.NewRecorder()
	srv.Router.ServeHTTP(w, req)
}
