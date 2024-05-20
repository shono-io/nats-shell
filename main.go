/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
  "github.com/nats-io/jsm.go/natscontext"
  "github.com/shono-io/nats-shell/shell"
  "log"
  "os"
)

func main() {
  nc, err := natscontext.Connect(natscontext.SelectedContext())
  if err != nil {
    log.Panicln(err)
  }

  rootCmd, err := shell.Command(nc)
  if err != nil {
    log.Panicln(err)
  }

  if err := rootCmd.Execute(); err != nil {
    os.Exit(1)
  }
}
