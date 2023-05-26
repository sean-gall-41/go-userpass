package main

import(
  "os"
  "log"
  "github.com/sean-gall-41/go-userpass/internal"
)

func main() {
  if err := internal.StartMySQL(); err != nil {
    log.Fatal(err)
    os.Exit(-1)
  }
  if err := startServer(); err != nil {
    log.Fatal(err)
    os.Exit(-1)
  }
}

