// go mod init golangcity.com/golangcity
package main

import (
  "crypto/tls"
  "database/sql"
  "flag"
  "html/template"
  "log"
  "net/http"
  "os"
  "time"
  "golangcity.com/golangcity/pkg/models"
  "golangcity.com/golangcity/pkg/models/mysql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/golangcollege/sessions"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
  debug bool
  errorLog *log.Logger
  infoLog *log.Logger
  session *sessions.Session
  posts interface {
    Insert(string, string, string) (int, error)
    Get(int) (*models.Post, error)
    Latest() ([]*models.Post, error)
    Delete(int)
  }
  templateCache map[string]*template.Template
  users interface {
    Insert(string, string, string) error
    Authenticate(string, string) (int, error)
    Get(int) (*models.User, error)
    ChangePassword(int, string, string) error
  }
}

func main() {
  dsn := flag.String("dsn", "webcity:9jaZxapeoy4n1qk2w@/golangcity?parseTime=true", "MySQL data source name")
  debug := flag.Bool("debug", false, "Enable debug mode")
  addr := flag.String("addr", ":4000", "HTTP network address")
  secret := flag.String("secret", "s1Ndh+pPbnxHbS*+9Ak8qGWhTQbpa@ge", "Secret key")
  flag.Parse()

  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  db, err := openDB(*dsn)
  if err != nil {
    errorLog.Fatal(err)
  }
  defer db.Close()

  templateCache, err := newTemplateCache("./ui/html/")
  if err != nil {
    errorLog.Fatal(err)
  }

  session := sessions.New([]byte(*secret))
  session.Lifetime = 12 *time.Hour
  session.Secure = true

  app := &application{
    debug: *debug,
    errorLog: errorLog,
    infoLog: infoLog,
    session: session,
    posts: &mysql.PostModel{DB: db},
    templateCache: templateCache,
    users: &mysql.UserModel{DB: db},
  }

  tlsConfig := &tls.Config{
    PreferServerCipherSuites: true,
    CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
  }

  srv := &http.Server{
    Addr: *addr,
    ErrorLog: errorLog,
    Handler: app.routes(),
    TLSConfig: tlsConfig,
    IdleTimeout: time.Minute,
    ReadTimeout: 5 *time.Second,
    WriteTimeout: 10 *time.Second,
  }

  infoLog.Printf("Starting server on %s", *addr)
  err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
  errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
  db, err := sql.Open("mysql", dsn)
  if err != nil {
    return nil, err
  }
  if err = db.Ping(); err != nil {
    return nil, err
  }
  return db, nil
}

/*
CREATE DATABASE golangcity CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE posts (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  code TEXT NOT NULL,
  created DATETIME NOT NULL
);
flag ==> question, feedback, interesting,...
upvote == find similar

CREATE INDEX idx_posts_created ON posts(created);

CREATE USER 'webcity'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON golangcity.* TO 'webcity'@'localhost';
ALTER USER 'webcity'@'localhost' IDENTIFIED BY '9jaZxapeoy4n1qk2w';

go get github.com/go-sql-driver/mysql@v1
*/

/*
USE golangcity;

CREATE TABLE users (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  hashed_password CHAR(60) NOT NULL,
  created DATETIME NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE(email);
*/

/*
CREATE DATABASE test_golangcity CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER 'test_web'@'localhost';
GRANT CREATE, DROP, ALTER, INDEX, SELECT, INSERT, UPDATE, DELETE ON test_golangcity.* TO 'test_web'@'localhost';
ALTER USER 'test_web'@'localhost' IDENTIFIED BY 'pass';
*/
