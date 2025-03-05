package main

import (
  "fmt"
  "os"

  "github.com/dgutierrez1287/vault-util/cmd"
)

func main() {
  if err := cmd.Execute(); err != nil {
    fmt.Println("Error executing cmd")
    os.Exit(1)
  }
}
