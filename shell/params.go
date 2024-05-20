package shell

const (
  ParamKindString ParamKind = iota
  ParamKindInt
  ParamKindFloat
  ParamKindBool
  ParamKindDuration
)

type ParamKind int64

func (p ParamKind) String() string {
  switch p {
  case ParamKindString:
    return "string"
  case ParamKindInt:
    return "int"
  case ParamKindFloat:
    return "float"
  case ParamKindBool:
    return "bool"
  case ParamKindDuration:
    return "duration"
  default:
    return "unknown"
  }
}

type ParamOption func(*Parameter)

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
