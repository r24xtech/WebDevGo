package main

import (
  "net/http"
  "github.com/bmizerany/pat"
  "github.com/justinas/alice"
)

func (app *application) routes() http.Handler{
  standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
  dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

  mux := pat.New()
  mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
  mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))
  mux.Get("/post/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createPostForm))
  mux.Post("/post/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createPost))
  mux.Get("/post/:id", dynamicMiddleware.ThenFunc(app.showPost))

  mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
  mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
  mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
  mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
  mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))
  mux.Get("/user/profile", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.userProfile))
  mux.Get("/user/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.changePasswordForm))
  mux.Post("/user/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.changePassword))

  mux.Get("/ping", http.HandlerFunc(ping))

  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Get("/static/", http.StripPrefix("/static", fileServer))

  return standardMiddleware.Then(mux)
}
