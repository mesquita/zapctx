# zapctx

Opinionated logger setup for [zap](https://github.com/uber-go/zap) that should make it easier for your project.
This project only gives you a configured logger, there is no wrapper around zap, so all zap features are available, e.g.
fields. Full zap documentation available [here](https://pkg.go.dev/go.uber.org/zap).

### Usage

Typical usage will be done by calling the log function from the context:

```go
ctx := context.Background() // ideally you will have a proper context available, this is just an example.
zapctx.Error(ctx, "error message")
zapctx.Info(ctx, "info message")
zapctx.Debug(ctx, "debug message")
```

Another provided shorthand function is available to add multiple fields to the logger in the context:

```go
ctx := context.Background() // ideally you will have a proper context available, this is just an example.
ctx = WithFields(ctx, zap.String("key", "value"), zap.String("other-key", "other-value"))
```

#### Customizations

##### Log level
If you want to change the log level, default is `Info`, there is a function available to change the log level in runtime:

```go
zapctx.SetLogLevel(zapcore.DebugLevel)
```

This might be useful if you want to use a different log level for development. It's up to each project to define it, the 
default log level is `info`

##### Hooks

If you want to run custom code before each log entry you can provide Hook functions like this:

```go
zapctx.AddHook(func(entry zapcore.Entry) error {
	// your handling goes here
	return nil
})
```

##### Logger override

If you want to completely override the logger in context use the `With` function, bear in mind that this will only be valid
for the ctx where you do that, any other ctx will not use that custom logger.

### NFRs

- Logger is configured to output json.

- Provides a function to set a TraceID field in the context so that your TraceID is propagated to all logs. Your project
  code typically shouldn't need this function, it's up to drivers and middlewares to set this without intervention.

```go
ctx := context.Background() // ideally you will have a proper context available, this is just an example.
ctx = zapctx.WithTraceID(ctx, "trace-id")
```

### Tests

```go
go test -v ./...
```

### Linting

This projects uses [golangci-lint](https://golangci-lint.run).

```
$ golangci-lint run
```

If you want to have the linter automatically fix the issues for you (it won't fix all of them) you can run:

```
$ golangci-lint run --fix 
```
