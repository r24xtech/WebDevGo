package main

import (
  "database/sql"
  "log"
  "net/http"
  "to-do-list/to-do/models/mysql"
  "to-do-list/to-do/models"
  _ "github.com/go-sql-driver/mysql"
)

type application struct {
  list *mysql.ListModel
}

type htmlData struct {
  ListData []*models.List // must be uppercase!!
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

func main() {
  dsn := "test:test@/to_do_list?parseTime=true"
  db, err := openDB(dsn)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  app := &application{
    list: &mysql.ListModel{DB: db},
  }
  srv := &http.Server{
    Addr: ":8080",
    Handler: app.routes(),
  }
  log.Print("Starting server on :8080")
  log.Fatal(srv.ListenAndServe())
}


/*
go mod init to-do-list/to-do
go get github.com/go-sql-driver/mysql@v1
*/

/*
[MySQL]
mysql -u root -p
show databases;
use to_do_list;
show tables;
select * from list;
create user 'test'@'localhost';
grant select, insert, update, delete on list.* to 'test'@'localhost';
grant select, insert, update, delete on to_do_list.* to 'test'@'localhost';
alter user 'test'@'localhost' identified by 'test';

db: to_do_list
table: list
(id)
(item)
(created)
*/
