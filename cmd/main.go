package main

import (
  "errors"
  "flag"
  "fmt"
  "os"

  "elf_cmp/cmd/internal/compare"
)

type Config struct {
  Fname1    string
  Fname2    string
  Html      bool
}

var config Config

func parseArgs() error {
  flag.BoolVar(&config.Html, "html", false, "generate html report")

  flag.Parse()
  args := flag.Args()

  if len(args) != 2 {
    return errors.New("Need two arguments")
  }

  config.Fname1 = flag.Arg(0)
  config.Fname2 = flag.Arg(1)
  return nil
}

func checkConfig() {
  if _, err := os.Stat(config.Fname1); errors.Is(err, os.ErrNotExist) {
    println("Error: ", config.Fname1, "does not exist")
  }

  if _, err := os.Stat(config.Fname2); errors.Is(err, os.ErrNotExist) {
    println("Error: ", config.Fname2, "does not exist")
  }
}

func main() {
  err := parseArgs()
  if err != nil {
    fmt.Println(err)
    return
  }
  checkConfig()
  err = compare.Compare(config.Fname1, config.Fname2, config.Html)
  if err != nil {
    fmt.Println(err)
  }
}
