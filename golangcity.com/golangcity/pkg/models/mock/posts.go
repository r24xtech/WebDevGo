package mock

import (
  "time"
  "golangcity.com/golangcity/pkg/models"
)

var mockPost = &models.Post{
  ID: 1,
  Title: "Print Function",
  Content: "This is how you print text",
  Code: "fmt.Println(\"Printing to the console!\")",
  Created: time.Now(),
}

type PostModel struct {}

func (m *PostModel) Insert(title, content, code string) (int, error) {
  return 2, nil
}

func (m *PostModel) Get(id int) (*models.Post, error) {
  switch id {
  case 1:
    return mockPost, nil
  default:
    return nil, models.ErrNoRecord
  }
}

func (m *PostModel) Latest() ([]*models.Post, error) {
  return []*models.Post{mockPost}, nil
}

func (m *PostModel) Delete(id int) {
  return
}
