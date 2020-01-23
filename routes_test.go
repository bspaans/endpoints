package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type suite struct{}

var _ = Suite(&suite{})

func (s *suite) SetUpTest(c *C) {
}

func getRouter_Add_Endpoints(checkSucceeds bool) *mux.Router {
	unit := NewRoutes()
	handlerFunc := func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
	}
	handlerFunc2 := func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
	}
	handlerFunc3 := func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
	}
	handlerFunc4 := func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
	}
	unit.AddEndpoints([]*Endpoint{
		NewEndpoint("/test", GET(handlerFunc), POST(handlerFunc2)),
		NewEndpoint("/test2", DELETE(handlerFunc3), PUT(handlerFunc4)),
	})
	return unit.GetMux()
}

func getRouter_Add_One_Route(checkSucceeds bool) *mux.Router {
	unit := NewRoutes()
	handlerFunc := func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
	}
	handlerFunc2 := func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
	}
	unit.AddEndpoint(NewEndpoint("/test",
		GET(handlerFunc),
		POST(handlerFunc2),
	))
	return unit.GetMux()
}

func (s *suite) Test_Routes_happy_path(c *C) {
	router := getRouter_Add_One_Route(true)
	test := NewRoutesTest(200)
	test.RunRouter(c, router)
}

func (s *suite) Test_Routes_AddRoutes(c *C) {
	router := getRouter_Add_Endpoints(true)
	test := NewRoutesTest(200)
	test.ExpectTest2ToSucceed = true
	test.RunRouter(c, router)
}

type RoutesTest struct {
	ExpectedCode         int
	ExpectTest2ToSucceed bool
}

func NewRoutesTest(statusCode int) *RoutesTest {
	return &RoutesTest{
		ExpectedCode: statusCode,
	}
}

func (r *RoutesTest) DoRequest(router *mux.Router, method, url string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, nil)
	router.ServeHTTP(rr, req)
	return rr
}

func (r *RoutesTest) RunRouter(c *C, router *mux.Router) {
	rr := r.DoRequest(router, "GET", "/test")
	c.Assert(rr.Code, Equals, r.ExpectedCode)

	rr = r.DoRequest(router, "POST", "/test")
	c.Assert(rr.Code, Equals, r.ExpectedCode)

	rr = r.DoRequest(router, "PUT", "/test")
	c.Assert(rr.Code, Equals, 405)
	c.Assert(rr.Header().Get("Allow"), Equals, "GET,POST")

	rr = r.DoRequest(router, "DELETE", "/test")
	c.Assert(rr.Code, Equals, 405)
	c.Assert(rr.Header().Get("Allow"), Equals, "GET,POST")

	rr = r.DoRequest(router, "DELETE", "/test2")
	expected := 404
	if r.ExpectTest2ToSucceed {
		expected = 200
	}
	c.Assert(rr.Code, Equals, expected)
}
