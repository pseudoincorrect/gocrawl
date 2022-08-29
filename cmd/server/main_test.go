package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("TestMain")
	app = App{}
	app.initializeRouter()
	fmt.Println(app)
	code := m.Run()
	os.Exit(code)
}

func TestCrawlAnURL(t *testing.T) {
	data := []struct {
		body string
		code int
	}{
		{body: `{"url":"http://www.data.com"}`, code: 200},
		{body: `{"url":"www.data.com"}`, code: 400},
		{body: `{"ur":"www.data.com"}`, code: 400},
		{body: "", code: 400},
	}
	for _, d := range data {
		req, _ := http.NewRequest("POST", "/crawl", strings.NewReader(d.body))
		res := executeRequest(req)
		if res.Code != d.code {
			t.Logf("Body: %s", d.body)
			t.Errorf("Expected error code to be %d, got %d", d.code, res.Code)
		}
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}
