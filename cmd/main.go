package main

import (
  "errors"
  "flag"
  "os"

  "elf_cmp/cmd/internal/compare"
)

type Config struct {
  Fname1 string
  Fname2 string
}

var config Config

func parseArgs() {
  flag.Parse()
  args := flag.Args()

  if len(args) != 2 {
    println("Need two arguments")
    return
  }

  config.Fname1 = flag.Arg(0)
  config.Fname2 = flag.Arg(1)
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
  parseArgs()
  checkConfig()
  err := compare.Compare(config.Fname1, config.Fname2)
  if err != nil {
    println("Error")
  }
}
