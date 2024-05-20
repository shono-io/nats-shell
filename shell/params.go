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
