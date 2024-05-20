package shell

// Much of the code in this file is copied from the natscli (https://github.com/nats-io/natscli/) project since I
// didn't want to reinvent the excellent stuff they did on this already. I've made some modifications to the code
// to fit my needs.

import (
  "bytes"
  "context"
  "encoding/json"
  "fmt"
  "github.com/klauspost/compress/s2"
  "github.com/nats-io/nats.go"
  "io"
  "strings"
  "sync"
  "time"
)

func doReqAsync(req any, subj string, waitFor int, timeout time.Duration, nc *nats.Conn, cb func([]byte)) error {
  jreq := []byte("{}")
  var err error

  if req != nil {
    switch val := req.(type) {
    case string:
      jreq = []byte(val)
    default:
      jreq, err = json.Marshal(req)
      if err != nil {
        return err
      }
    }
  }

  var (
    mu       sync.Mutex
    ctr      = 0
    finisher *time.Timer
  )

  ctx, cancel := context.WithTimeout(context.Background(), timeout)
  defer cancel()

  if waitFor == 0 {
    finisher = time.NewTimer(timeout)
    go func() {
      select {
      case <-finisher.C:
        cancel()
      case <-ctx.Done():
        return
      }
    }()
  }

  errs := make(chan error)
  sub, err := nc.Subscribe(nc.NewRespInbox(), func(m *nats.Msg) {
    mu.Lock()
    defer mu.Unlock()

    data := m.Data
    if m.Header.Get("Content-Encoding") == "snappy" {
      ud, err := io.ReadAll(s2.NewReader(bytes.NewBuffer(data)))
      if err != nil {
        errs <- err
        return
      }
      data = ud
    }

    if finisher != nil {
      finisher.Reset(300 * time.Millisecond)
    }

    if m.Header.Get("Status") == "503" {
      errs <- nats.ErrNoResponders
      return
    }

    cb(data)
    ctr++

    if waitFor > 0 && ctr == waitFor {
      cancel()
    }
  })
  if err != nil {
    return err
  }
  defer sub.Unsubscribe()

  if waitFor > 0 {
    sub.AutoUnsubscribe(waitFor)
  }

  msg := nats.NewMsg(subj)
  msg.Data = jreq
  if subj != "$SYS.REQ.SERVER.PING" && !strings.HasPrefix(subj, "$SYS.REQ.ACCOUNT") {
    msg.Header.Set("Accept-Encoding", "snappy")
  }
  msg.Reply = sub.Subject

  err = nc.PublishMsg(msg)
  if err != nil {
    return err
  }

  select {
  case err = <-errs:
    if err == nats.ErrNoResponders && strings.HasPrefix(subj, "$SYS") {
      return fmt.Errorf("server request failed, ensure the account used has system privileges and appropriate permissions")
    }

    return err
  case <-ctx.Done():
  }

  return nil
}

func doReq(req any, subj string, waitFor int, timeout time.Duration, nc *nats.Conn) ([][]byte, error) {
  res := [][]byte{}
  mu := sync.Mutex{}

  err := doReqAsync(req, subj, waitFor, timeout, nc, func(r []byte) {
    mu.Lock()
    res = append(res, r)
    mu.Unlock()
  })

  return res, err
}
