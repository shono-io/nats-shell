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

    cmd, err := parseInfo(nc, i)
    if err != nil {
      return nil, fmt.Errorf("failed to parse API info: %w", err)
    }

    rootCmd.AddCommand(cmd)
  }

  return rootCmd, nil
}

func parseInfo(nc *nats.Conn, info micro.Info) (*cobra.Command, error) {
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

  return asCommand(nc, &root)
}

func asCommand(nc *nats.Conn, e *treeElement) (*cobra.Command, error) {
  cmd := &cobra.Command{
    Use:   e.Name,
    Short: e.Summary,
  }

  if e.Endpoint != nil {
    var params []Parameter
    if e.Endpoint.Metadata != nil {
      m := Metadata(e.Endpoint.Metadata)
      cmd.Short = m.Summary()
      cmd.Long = m.Description()

      params = m.Parameters()
      for _, p := range params {
        switch p.Kind {
        case ParamKindString:
          dv := ""
          if p.Default != nil {
            dv = p.Default.(string)
          }
          cmd.Flags().String(p.Name, dv, p.Summary)
        case ParamKindInt:
          dv := 0
          if p.Default != nil {
            dv = p.Default.(int)
          }
          cmd.Flags().Int(p.Name, dv, p.Summary)
        case ParamKindBool:
          dv := false
          if p.Default != nil {
            dv = p.Default.(bool)
          }
          cmd.Flags().Bool(p.Name, dv, p.Summary)
        case ParamKindFloat:
          dv := 0.0
          if p.Default != nil {
            dv = p.Default.(float64)
          }
          cmd.Flags().Float64(p.Name, dv, p.Summary)
        case ParamKindDuration:
          dv := 0 * time.Second
          if p.Default != nil {
            pd, err := time.ParseDuration(p.Default.(string))
            if err == nil {
              dv = pd
            }
          }
          cmd.Flags().Duration(p.Name, dv, p.Summary)
        default:
          return nil, fmt.Errorf("unknown parameter kind: %v", p.Kind)
        }
      }
    }

    cmd.Run = func(cmd *cobra.Command, args []string) {
      // -- retrieve the parameters
      ps := map[string]any{}
      for _, p := range params {
        var err error
        switch p.Kind {
        case ParamKindString:
          ps[p.Name], err = cmd.Flags().GetString(p.Name)
          if p.Required && ps[p.Name].(string) == "" {
            fmt.Printf("parameter %s is required\n", p.Name)
            return
          }
        case ParamKindInt:
          ps[p.Name], err = cmd.Flags().GetInt(p.Name)
          if p.Required && ps[p.Name].(int) == 0 {
            fmt.Printf("parameter %s is required\n", p.Name)
            return
          }
        case ParamKindBool:
          ps[p.Name], err = cmd.Flags().GetBool(p.Name)
        case ParamKindFloat:
          ps[p.Name], err = cmd.Flags().GetFloat64(p.Name)
          if p.Required && ps[p.Name].(float64) == 0 {
            fmt.Printf("parameter %s is required\n", p.Name)
            return
          }
        case ParamKindDuration:
          ps[p.Name], err = cmd.Flags().GetDuration(p.Name)
          if p.Required && ps[p.Name].(time.Duration) == 0 {
            fmt.Printf("parameter %s is required\n", p.Name)
            return
          }
        default:
          err = fmt.Errorf("unknown parameter kind: %v", p.Kind)
        }

        if err != nil {
          fmt.Printf("failed to get parameter %s: %v\n", p.Name, err)
          return
        }

        resp, err := doReq(ps, e.Endpoint.Subject, 0, 10*time.Second, nc)
        if err != nil {
          fmt.Printf("failed to send request: %v\n", err)
          return
        }

        for _, r := range resp {
          fmt.Println(string(r))
        }
      }
    }
  }

  for _, c := range e.Children {
    sub, err := asCommand(nc, c)
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
