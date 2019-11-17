HANDLER ==> func (w http.ResponseWriter, r *http.Request)
The request handler alone can not accept any HTTP connections from the outside.
An HTTP server has to listen on a port to pass connections on to the
request handler.
handler is an object which satisfies the http.Handler interface
[HTTP server]
==> process dynamic requests - http.HandleFunc
==> serve static assets - http.FileServer
==> accept connections - listen on a port
http://<servername>/<handlername>?<parameters>
multiplexer ==> the piece of code that redirects a request to a handler.


```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```


<hr><hr>

**Resources**

https://cryptic.io/go-http/

https://www.alexedwards.net/blog/a-recap-of-request-handling

https://www.alexedwards.net/blog/interfaces-explained
