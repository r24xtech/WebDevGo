package main

import (
  "fmt"
  "strconv"
  "html/template"
  "net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.Error(w, "not found", http.StatusNotFound)
    return
  }
  allItems, err := app.list.GetAll()
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }
  data := &htmlData{ListData: allItems}
  temp, err := template.ParseFiles("./html/index.html")
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }
  err = temp.Execute(w, data)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
  }
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    w.Header().Set("Allow", http.MethodPost)
    http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
    return
  }
  if err := r.ParseForm(); err != nil {
    fmt.Fprintf(w, "ParseForm() err: %v", err)
    return
  }
  item := r.FormValue("item")
  err := app.list.Insert(item)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    w.Header().Set("Allow", http.MethodPost)
    http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
    return
  }
  id, err := strconv.Atoi(r.URL.Query().Get(":id"))
  if err != nil || id < 1 {
    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
    return
  }
  err = app.list.Delete(id)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/", http.StatusSeeOther)
}
