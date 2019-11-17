# WebDevGo
Web development in Go


HANDLER ==> `func (w http.ResponseWriter, r *http.Request)`
The request handler alone can not accept any HTTP connections from the outside. An HTTP server has to listen on a port to pass connections on to the request handler. Handler is an object which satisfies the `http.Handler` interface.<br>

**[HTTP server]**
* process dynamic requests:: http.HandleFunc
* serve static assets:: http.FileServer
* accept connections:: listen on a port
`http://<servername>/<handlername>?<parameters>`
* multiplexer:: the piece of code that redirects a request to a handler.


```go
// The building block of the entire net/http package is the http.Handler interface
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

Once implemented the `http.Handler` can be passed to `http.ListenAndServe`, which will call the `ServeHTTP` method on every incoming request.

Often defining a full type to implement the `http.Handler` interface is a bit overkill, especially for extremely simple `ServeHTTP` functions. The net/http package provides a helper function, `http.HandlerFunc`, which wraps a function which has the signature `func(w http.ResponseWriter, r *http.Request)`, returning an `http.Handler` which will call it in all cases.

The `http.ServeMux` is itself an `http.Handler`, so it can be passed into `http.ListenAndServe`. When it receives a request it will check if the request’s path is prefixed by any of its known paths, choosing the longest prefix match it can find. We use the `/` endpoint as a catch-all to catch any requests to unknown endpoints.

`http.ServeMux` has both `Handle` and `HandleFunc` methods. These do the same thing, except that `Handle` takes in an `http.Handler` while `HandleFunc` merely takes in a function, implicitly wrapping it just as `http.HandlerFunc` does.
<br><br>

**Processing HTTP requests with Go is primarily about two things: ServeMuxes and Handlers.**<br>
A ServeMux is essentially a HTTP request router (or multiplexor). It compares incoming requests against a list of predefined URL paths, and calls the associated handler for the path whenever a match is found.

Handlers are responsible for writing response headers and bodies. Almost any object can be a handler, so long as it satisfies the http.Handler interface. 





![GitHub Logo](/listenandserve.png)
![GitHub Logo](/handlerfunc.png)

<br>
<hr><hr>
<br>

**Resources**

https://cryptic.io/go-http/

https://www.alexedwards.net/blog/a-recap-of-request-handling

https://www.alexedwards.net/blog/interfaces-explained

https://golang.org/pkg/net/http/#HandlerFunc


