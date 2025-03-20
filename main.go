package main

import (
	"fmt"

	"github.com/ben-lehman/gorss/internal/config"
)
func main() {
  conf, err := config.Read()
  if err != nil {
    fmt.Printf("Error: %v\n", err)
  }
  conf.SetUser("ben")
  conf, err = config.Read()
  if err != nil {
    fmt.Printf("Error: %v\n", err)
  }
  fmt.Printf("db url: %v\ncurrent user: %v\n", conf.DbURL, conf.CurrentUsername)
  return
}
