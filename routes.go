package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Routes struct {
	OutputRoutes        bool
	endpoints           []*Endpoint
	registeredEndpoints map[string]bool
	result              *mux.Router
	methodRouters       map[string]*mux.Router
}

func NewRoutes() *Routes {
	return &Routes{
		endpoints:           []*Endpoint{},
		registeredEndpoints: map[string]bool{},
		result:              mux.NewRouter(),
		methodRouters:       map[string]*mux.Router{},
	}
}

func (e *Routes) AddEndpoint(endpoint *Endpoint) {
	if _, exists := e.registeredEndpoints[endpoint.URL]; exists {
		panic("Endpoint for " + endpoint.URL + " has already been registered")
	}
	e.registeredEndpoints[endpoint.URL] = true
	e.endpoints = append(e.endpoints, endpoint)
}

func (r *Routes) AddEndpoints(endpoints []*Endpoint) {
	for _, endpoint := range endpoints {
		r.AddEndpoint(endpoint)
	}
}

func (r *Routes) getMethodRouter(method string) *mux.Router {
	router, found := r.methodRouters[method]
	if !found {
		router = r.result.Methods(method).Subrouter()
		r.methodRouters[method] = router
	}
	return router
}

func (r *Routes) GetMux() *mux.Router {
	for _, endpoint := range r.endpoints {
		supportedMethods := map[string]bool{
			"GET": false, "POST": false, "PUT": false, "DELETE": false, "OPTIONS": false,
		}
		for _, method := range endpoint.Methods {
			if r.OutputRoutes {
				log.Println("INFO: Adding endpoint", method.Method, endpoint.URL)
			}
			router := r.getMethodRouter(method.Method)
			supportedMethods[method.Method] = true
			router.HandleFunc(endpoint.URL, method.HandlerFunc)
		}
		r.addMissingMethods(endpoint.URL, supportedMethods)
	}
	return r.result
}

func (r *Routes) addMissingMethods(endpoint string, supportedMethods map[string]bool) {
	order := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	allowed := []string{}
	for _, method := range order {
		if supportedMethods[method] {
			allowed = append(allowed, method)
		}
	}
	for method, enabled := range supportedMethods {
		if !enabled {
			router := r.getMethodRouter(method)
			handler := r.MethodNotSupportedHandler(method, endpoint, allowed)
			router.HandleFunc(endpoint, handler)
		}
	}

}

func (r *Routes) MethodNotSupportedHandler(method, url string, allowed []string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == url {
			rw.WriteHeader(405)
			rw.Header().Add("Allow", strings.Join(allowed, ","))
		} else {
			rw.WriteHeader(404)
		}
	}
}
