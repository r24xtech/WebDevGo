package models

import (
  "errors"
  "time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type List struct { //need uppercase letter here to import
  ID int
  Item string
  Created time.Time
}
