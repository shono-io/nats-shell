# nats-shell
Nats micro CLI based on micro metadata.

## Why?
Building a CLI for your nats micro services is a bit of a pain. This project provides an api by introspecting the
micro services metadata and provides a CLI to interact with the services.

## How?
When nats-shell is started, it will connect your currently active context. It will introspect the services available
on the nats node and parse additional information present in the metadata of those services. It will then move on to
create the commands for the endpoints it discovered.

## Installation
```bash
go get github.com/shonoio/nats-shell
```

## Getting Started
You will need to have an active nats context in order to work with nats-shell. Therefor, you might want to rely on the
nats-cli to create a context for you:
```shell
nats context create my-context --server nats://localhost:4222 --select
```

Once you have an active context, you can start nats-shell:
```shell
nats-shell
```

The default will show all service which have been detected

### Augmenting your services
nats-shell will work with any services which have been declared through micro, but there is some additional metadata 
that can be provided to steer the way nats-shell will interpret your services. This metadata can be provided when your
endpoint is being registered:
```go
_ = svc.AddEndpoint(name,
    handler,
    micro.WithEndpointMetadata(shell.NewMetadata(
      shell.WithSummary("This is the short description of the endpoint")
      shell.WithDescription("This is a the long description of the endpoint"),
      shell.WithParameters(
        shell.NewParameter("param1", shell.ParamKindString,
          shell.WithParameterSummary("the first argument"),
          shell.Required()),
        shell.NewParameter("param2", shell.ParamKindString,
          shell.WithParameterSummary("the second argument"),
          shell.WithDefault("default value")),
      ),
    )),
  )
```