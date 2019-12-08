package main

import (
  "net/http"
  "github.com/bmizerany/pat"
  "github.com/justinas/alice"
)
// page 179: pat's pattern matching order
func (app *application) routes() http.Handler {
  standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

  mux := pat.New()
  mux.Get("/", http.HandlerFunc(app.home))
  mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
  mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
  mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))
  mux.Post("/snippet/:id/delete", http.HandlerFunc(app.deleteSnippet))

  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Get("/static/", http.StripPrefix("/static", fileServer))

  return standardMiddleware.Then(mux)
}
