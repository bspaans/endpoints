package routes

import (
	"net/http"
)

type Endpoint struct {
	URL     string
	Methods []*Method
}

func NewEndpoint(url string, methods ...*Method) *Endpoint {
	return &Endpoint{
		URL:     url,
		Methods: methods,
	}
}

type Method struct {
	Method      string
	HandlerFunc http.HandlerFunc
}

func NewMethod(method string, handler http.HandlerFunc) *Method {
	return &Method{
		Method:      method,
		HandlerFunc: handler,
	}
}

func GET(handler http.HandlerFunc) *Method {
	return NewMethod("GET", handler)
}
func POST(handler http.HandlerFunc) *Method {
	return NewMethod("POST", handler)
}
func PUT(handler http.HandlerFunc) *Method {
	return NewMethod("PUT", handler)
}
func DELETE(handler http.HandlerFunc) *Method {
	return NewMethod("DELETE", handler)
}
func OPTIONS(handler http.HandlerFunc) *Method {
	return NewMethod("OPTIONS", handler)
}
