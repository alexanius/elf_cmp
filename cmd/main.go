package main

import (
  "errors"
  "flag"
  "fmt"
  "os"

  "elf_cmp/cmd/internal/compare"
)

type Action int

const (
    Analyze Action = iota
    Compare
)

type Config struct {
  Action    Action
  Fname1    string
  Fname2    string
  Log       string // File with different runtime logs
  Html      bool
}

var config Config

func parseArgs() error {
  flag.BoolVar(&config.Html, "html", false, "generate html report")
  flag.StringVar(&config.Log, "log", "", "analyze runtime log file")

  flag.Parse()
  args := flag.Args()

  if len(args) == 0 {
    return errors.New("Choose action 'analyze' or 'compare'")
  }

  action := args[0]
  switch action {
  case "analyze":
    config.Action = Analyze
    if len(args) != 2 {
      if config.Log == "" {
        return errors.New("For analysis you need one file as an argument")
      } else {
        if _, err := os.Stat(config.Log); errors.Is(err, os.ErrNotExist) {
          return errors.New("Error: file '" + config.Log + "' does not exist")
        }
      }
    }
  case "compare":
    config.Action = Compare
    if len(args) != 3 {
      return errors.New("For compare you need two files as an argument")
    }
  default:
    return errors.New("Unknown action: " + action)
  }

  if len(args) == 1 {
    return nil
  }
  config.Fname1 = args[1]
  if _, err := os.Stat(config.Fname1); errors.Is(err, os.ErrNotExist) {
    return errors.New("Error: file '" + config.Fname1 + "' does not exist")
  }

  if config.Action == Compare {
    config.Fname2 = args[2]
    if _, err := os.Stat(config.Fname2); errors.Is(err, os.ErrNotExist) {
      return errors.New("Error: file '" + config.Fname2 + "' does not exist")
    }
  }
  return nil
}

func main() {
  err := parseArgs()
  if err != nil {
    fmt.Println(err)
    return
  }

  switch config.Action {
  case Analyze:
    if config.Log != "" {
      err = compare.AnalyzeLog(config.Log)
    }
  case Compare:
    err = compare.Compare(config.Fname1, config.Fname2, config.Html)
  }
  if err != nil {
    fmt.Println(err)
  }
}

