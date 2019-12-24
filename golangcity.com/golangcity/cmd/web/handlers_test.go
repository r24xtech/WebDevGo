package main

import (
  "bytes"
  "net/http"
  "net/url"
  "testing"
)

func TestPing(t *testing.T) {
  app := newTestApplication(t)
  ts := newTestServer(t, app.routes())
  defer ts.Close()

  code, _, body := ts.get(t, "/ping")

  if code != http.StatusOK {
    t.Errorf("want %d; got %d", http.StatusOK, code)
  }

  if string(body) != "OK" {
    t.Errorf("want body equal %q", "OK")
  }
}

func TestShowPost(t *testing.T) {
  app := newTestApplication(t)
  ts := newTestServer(t, app.routes())
  defer ts.Close()

  tests := []struct {
    name string
    urlPath string
    wantCode int
    wantBody []byte
  }{
    {"Valid ID", "/post/1", http.StatusOK, []byte("Print Function")},
    {"Non-existent ID", "/post/2", http.StatusNotFound, nil},
    {"Negative ID", "/post/-1", http.StatusNotFound, nil},
    {"Decimal ID", "/post/1.23", http.StatusNotFound, nil},
    {"String ID", "/post/foo", http.StatusNotFound, nil},
    {"Empty ID", "/post/", http.StatusNotFound, nil},
    {"Trailing slash", "/post/1/", http.StatusNotFound, nil},
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T){
      code, _, body := ts.get(t, tt.urlPath)

      if code != tt.wantCode {
        t.Errorf("want %d; got %d", tt.wantCode, code)
      }

      if !bytes.Contains(body, tt.wantBody) {
        t.Errorf("want body to contain %q", tt.wantBody)
      }
    })
  }
}

func TestSignupUser(t *testing.T) {
  app := newTestApplication(t)
  ts := newTestServer(t, app.routes())
  defer ts.Close()

  _, _, body := ts.get(t, "/user/signup")
  csrfToken := extractCSRFToken(t, body)

  tests := []struct {
    name string
    userName string
    userEmail string
    userPassword string
    csrfToken string
    wantCode int
    wantBody []byte
  }{
    {"Valid submission", "Bob", "bob@example.com", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
    {"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
    {"Empty email", "Bob", "", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
    {"Empty password", "Bob", "bob@example.com", "", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
    {"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
    {"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
    {"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
    {"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, http.StatusOK, []byte("This field is too short (minimum is 10 characters)")},
    {"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Address is already in use")},
    {"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      form := url.Values{}
      form.Add("name", tt.userName)
      form.Add("email", tt.userEmail)
      form.Add("password", tt.userPassword)
      form.Add("csrf_token", tt.csrfToken)

      code, _, body := ts.postForm(t, "/user/signup", form)

      if code != tt.wantCode {
        t.Errorf("want %d; got %d", tt.wantCode, code)
      }

      if !bytes.Contains(body, tt.wantBody) {
        t.Errorf("want body %s to contain %q", body, tt.wantBody)
      }
    })
  }
}

func TestCreatePostForm(t *testing.T) {
  app := newTestApplication(t)
  ts := newTestServer(t, app.routes())
  defer ts.Close()

  t.Run("Unauthenticated", func(t *testing.T) {
    code, headers, _ := ts.get(t, "/post/create")
    if code != http.StatusSeeOther {
      t.Errorf("want %d; got %d", http.StatusSeeOther, code)
    }
    if headers.Get("Location") != "/user/login" {
      t.Errorf("want %s; got %s", "/user/login", headers.Get("Location"))
    }
  })

  t.Run("Authenticated", func(t *testing.T) {
    _, _, body := ts.get(t, "/user/login")
    csrfToken := extractCSRFToken(t, body)

    form := url.Values{}
    form.Add("email", "alice@example.com")
    form.Add("password", "")
    form.Add("csrf_token", csrfToken)
    ts.postForm(t, "/user/login", form)

    code, _, body := ts.get(t, "/post/create")
    if code != 200 {
      t.Errorf("want %d; got %d", 200, code)
    }
    formTag := "<form action=\"/post/create\" method=\"POST\">"
    if !bytes.Contains(body, []byte(formTag)) {
      t.Errorf("want body %s to contain %s", body, formTag)
    }
  })
}
