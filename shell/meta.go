package shell

type Option func(map[string]string)

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
