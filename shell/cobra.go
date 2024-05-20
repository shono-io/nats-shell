package shell

import (
  "encoding/json"
  "fmt"
  "github.com/nats-io/nats.go"
  "github.com/nats-io/nats.go/micro"
  "github.com/spf13/cobra"
  "strings"
  "time"
)

func Command(nc *nats.Conn) (*cobra.Command, error) {
  subject := fmt.Sprintf("%s.INFO", micro.APIPrefix)

  resp, err := doReq(nil, subject, 0, 10*time.Second, nc)
  if err != nil {
    return nil, fmt.Errorf("failed to get API info: %w", err)
  }

  rootCmd := &cobra.Command{
    Use:   "nats-shell",
    Short: "nats-shell is a CLI for interacting with the nats micro services",
  }

  for _, m := range resp {
    var i micro.Info
    if err := json.Unmarshal(m, &i); err != nil {
      return nil, fmt.Errorf("failed to unmarshal API info: %w", err)
    }

    cmd, err := parseInfo(i)
    if err != nil {
      return nil, fmt.Errorf("failed to parse API info: %w", err)
    }

    rootCmd.AddCommand(cmd)
  }

  return rootCmd, nil
}

func parseInfo(info micro.Info) (*cobra.Command, error) {
  root := treeElement{
    Name:    info.Name,
    Summary: info.Description,
  }

  for _, e := range info.Endpoints {
    cur := &root
    parts := strings.Split(e.Subject, ".")[1:]
    for _, p := range parts {
      if cur.Children == nil {
        cur.Children = make(map[string]*treeElement)
      }

      if _, ok := cur.Children[p]; !ok {
        cur.Children[p] = &treeElement{
          Name: p,
        }
      }

      cur = cur.Children[p]
    }

    cur.Endpoint = &e
  }

  return asCommand(&root)
}

func asCommand(e *treeElement) (*cobra.Command, error) {
  cmd := &cobra.Command{
    Use:   e.Name,
    Short: e.Summary,
  }

  if e.Endpoint != nil {
    if e.Endpoint.Metadata != nil {
      m := Metadata(e.Endpoint.Metadata)
      cmd.Short = m.Summary()
      cmd.Long = m.Description()
    }

    cmd.Run = func(cmd *cobra.Command, args []string) {

    }
  }

  for _, c := range e.Children {
    sub, err := asCommand(c)
    if err != nil {
      return nil, err
    }

    cmd.AddCommand(sub)
  }

  return cmd, nil
}

type treeElement struct {
  Name     string
  Summary  string
  Endpoint *micro.EndpointInfo
  Children map[string]*treeElement
}
