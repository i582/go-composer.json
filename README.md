# go-composer.json

go-composer.json is a small library for manipulating the `composer.json` configuration file.

### Supported

1. Resolving of namespaces for PSR-4 autoload.
2. Working with local dependencies, resolving paths to them.
3. Custom checks for config.

#### PSR-4

To resolve the path to the namespace, use the `Psr4PathForNamespace` method.

#### Custom checks

To add a custom check, use the `AddCheck` method. 

Example:

```go
cfg.AddCheck(func(config *composer.Config) *composer.ConfigError {
    if !strings.HasPrefix(config.Name, "my/") {
        return &composer.ConfigError{
            Msg:      "name must starts with prefix my/",
            Critical: true,
        }
    }
    return nil
})
```

### License

MIT

