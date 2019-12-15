package main

import (
  "net/http"
  "github.com/bmizerany/pat"
)

func (app *application) routes() *pat.PatternServeMux{
  mux := pat.New()
  mux.Get("/", http.HandlerFunc(app.home))
  mux.Post("/create", http.HandlerFunc(app.create))
  mux.Post("/delete/:id", http.HandlerFunc(app.delete))
  return mux
}
