package shell

type Option func(map[string]string)

func NewMetadata(m map[string]string, opts ...Option) map[string]string {
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

type Metadata map[string]string

func (m Metadata) Summary() string {
  return m["_nats.shell.summary"]
}

func (m Metadata) Description() string {
  return m["_nats.shell.description"]
}
