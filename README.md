# Rextporter
This is the executable entry point for the server.
- [Install](#install)
- [Run](#run)
- [Config file](#config-file)

## Install

```bash
$ cd $GOPATH/src/github.com/simelo/rextporter/cmd/rextporter
$ go install ./...
```

## Run

You can run the program (`rextporter`, make sure you have it accessible trough your `PATH` env variable) by calling it in the console and you have the following parameters options.

 - `-config` Metrics main config file path. (default to your home config folder + simelo -> rextporter -> main.toml).
 - `-handler` Handler to expose metric. (default "/metrics").
 - `-port` Listen port. (default 8080)

### Config file

You have 3 config files, main config(general definitions, like for example, load the service config from file and use "this" path), service config(service definitions) and metrics(metrics definitions) config.
Trough program parameter you can refer to main config path only, if you wan to use a custom path for service config file and/or metrics config file you should edit the main config it self. If you point to a not existent main config file it will use a default path and crete it in such default path. If you point to a non existent metrics or service path a default one will be create in such path.

Example main configuration file:
```toml
serviceConfigTransport = "file"
serviceConfigPath = "/home/adacosta/.config/simelo/rextporter/service.toml"
metricsConfigPath = "/home/adacosta/.config/simelo/rextporter/metrics.toml"
```

Example service configuration file:
```toml
# Service configuration.
name = "wallet"
scheme = "http"
port = 8000
basePath = ""
authType = "CSRF"
tokenHeaderKey = "X-CSRF-Token"
genTokenEndpoint = "/api/v1/csrf.json"
tokenKeyFromEndpoint = "csrf_token"

[location]
  location = "localhost"
```

Example metrics file configuration file.
```toml
# All metrics to be measured.
[[metrics]]
  name = "seq"
  url = "/api/v1/health.json"
  httpMethod = "GET"
  path = "/blockchain/head/seq"

  [metrics.options]
    type = "Counter"
    description = "I am running since"
```