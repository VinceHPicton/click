package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handleInsertTest(t *testing.T) {
	srv := Server{
		// db: mockDatabase,
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
