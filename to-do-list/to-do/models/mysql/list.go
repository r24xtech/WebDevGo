package mysql

import (
  "database/sql"
  "to-do-list/to-do/models"
)

type ListModel struct {
  DB *sql.DB
}

func (l *ListModel) Insert(item string) error {
  stmt := `INSERT INTO list (item, created) VALUES(?, UTC_TIMESTAMP())`

  _, err := l.DB.Exec(stmt, item)
  if err != nil {
    return err
  }
  return nil
}

func (l *ListModel) GetAll() ([]*models.List, error) {
  stmt := `SELECT * FROM list`
  rows, err := l.DB.Query(stmt)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  listItems := []*models.List{}
  for rows.Next() {
    item := &models.List{}
    err = rows.Scan(&item.ID, &item.Item, &item.Created)
    if err != nil {
      return nil, err
    }
    listItems = append(listItems, item)
  }
  if err = rows.Err(); err != nil {
    return nil, err
  }
  return listItems, nil
}

func (l *ListModel) Delete(id int) error {
  stmt := `DELETE FROM list WHERE id = ?`

  _, err := l.DB.Exec(stmt, id)
  if err != nil {
    return err
  }
  return nil
}
