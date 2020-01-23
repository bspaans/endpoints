# Endpoints

This is a go library to programmatically handle REST endpoints, which can be
handy when you want to generate code from an OpenAPI specification (see also:
[sheepcounter](https://github.com/bspaans/sheepcounter)). It uses `gorilla/mux`
for the actual routing.


## Usage


```
import (
	"github.com/bspaans/endpoints"
)

func main() {
    routes := endpoints.NewRoutes()
    method1 := endpoints.NewMethod("GET", myHttpHandlerFunc)
    method2 := endpoints.NewMethod("POST", myOtherHttpHandlerFunc)
    endpoint := endpoints.NewEndpoint("/test", method1, method2)
    routes.AddEndpoint(endpoints)
    handler := routes.GetMux()
    http.Handle("/", handler)
}
```
