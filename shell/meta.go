package shell

import "encoding/json"

type Option func(map[string]string)
type ParamOption func(*Parameter)

func NewMetadata(opts ...Option) Metadata {
  m := map[string]string{}

  for _, opt := range opts {
    opt(m)
  }

  return m
}

func WithSummary(summary string) Option {
  return func(m map[string]string) {
    m["_nats.shell.summary"] = summary
  }
}

func WithDescription(desc string) Option {
  return func(m map[string]string) {
    m["_nats.shell.description"] = desc
  }
}

func WithParameters(param ...Parameter) Option {
  return func(m map[string]string) {
    ps, fnd := m["_nats.shell.parameters"]
    if !fnd {
      ps = "[]"
    }

    var params []Parameter
    _ = json.Unmarshal([]byte(ps), &params)

    params = append(params, param...)

    b, _ := json.Marshal(params)
    m["_nats.shell.parameters"] = string(b)
  }
}

func WithEntries(m map[string]string) Option {
  return func(m map[string]string) {
    for k, v := range m {
      m[k] = v
    }
  }
}

type Metadata map[string]string

func (m Metadata) Summary() string {
  return m["_nats.shell.summary"]
}

func (m Metadata) Description() string {
  return m["_nats.shell.description"]
}

func (m Metadata) Parameters() []Parameter {
  ps, fnd := m["_nats.shell.parameters"]
  if !fnd {
    return nil
  }

  var params []Parameter
  _ = json.Unmarshal([]byte(ps), &params)

  return params
}

func WithParameterSummary(summary string) ParamOption {
  return func(p *Parameter) {
    p.Summary = summary
  }
}

func Required() ParamOption {
  return func(p *Parameter) {
    p.Required = true
  }
}

func WithDefaultValue(v any) ParamOption {
  return func(p *Parameter) {
    p.Default = v
  }
}

func NewParameter(name string, kind ParamKind, opts ...ParamOption) Parameter {
  p := Parameter{Name: name, Kind: kind}
  for _, opt := range opts {
    opt(&p)
  }
  return p
}

type Parameter struct {
  Name     string    `json:"name"`
  Kind     ParamKind `json:"kind"`
  Summary  string    `json:"summary,omitempty"`
  Required bool      `json:"required,omitempty"`
  Default  any       `json:"default"`
}
