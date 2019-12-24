package mysql

import (
  "database/sql"
  "errors"
  "golangcity.com/golangcity/pkg/models"
)

type PostModel struct {
  DB *sql.DB
}

func (p *PostModel) Insert(title, content, code string) (int, error) {
  stmt := `INSERT INTO posts (title, content, code, created)
  VALUES(?, ?, ?, UTC_TIMESTAMP())`

  result, err := p.DB.Exec(stmt, title, content, code)
  if err != nil {
    return 0, err
  }
  id, err := result.LastInsertId()
  if err != nil {
    return 0, err
  }
  return int(id), nil
}

func (p *PostModel) Get(id int) (*models.Post, error) {
  stmt := `SELECT id, title, content, code, created FROM posts WHERE id = ?`

  row := p.DB.QueryRow(stmt, id)
  px := &models.Post{}

  err := row.Scan(&px.ID, &px.Title, &px.Content, &px.Code, &px.Created)
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      return nil, models.ErrNoRecord
    } else {
      return nil, err
    }
  }
  return px, nil
}

func (p *PostModel) Latest() ([]*models.Post, error) {
  stmt := `SELECT * FROM posts` //no limit?!
  rows, err := p.DB.Query(stmt)
  if err != nil {
    return nil, err
  }
  defer rows.Close() // important!!
  posts := []*models.Post{}
  for rows.Next() {
    pz := &models.Post{}
    err = rows.Scan(&pz.ID, &pz.Title, &pz.Content, &pz.Code, &pz.Created)
    if err != nil {
      return nil, err
    }
    posts = append(posts, pz)
  }
  if err = rows.Err(); err != nil {
    return nil, err
  }
  return posts, nil
}

func (p *PostModel) Delete(id int) {
  return
}
