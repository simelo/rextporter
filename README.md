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

Trough program parameter you can refer to main config path only, if you wan to use a custom path for services config file, metrics for services or metrics you should edit the main config it self. If you point to a not existent main config file it will use a default path and crete it in such default path. If you point to a non existent (or not specify one) services, metrics for services, or metrics path a default one will be create in such path.

You have at least 4 config files:

- Main config (general definitions, like for example, load the service config from file and use "this" path)

- Services config (services definitions), the path should be maped from the main config.

- Metrics for service (metrics definitions) config, the path should be maped from the main config.

- `ServiceName+Metrics.toml` define the metrics of a giving service (`ServiceName`), the path should be mapped from metrics for service. You can have multiple ServiceNameMetric(`skycoinMetric.toml` for instance), depending on the number of service and hoh they are mapped with metrics.

Example main configuration file:
```toml
serviceConfigTransport = "file"
servicesConfigPath = "/home/adacosta/.config/simelo/rextporter/service.toml"
metricsForServicesPath = "/home/adacosta/.config/simelo/rextporter/metricsForServices.toml"
```

Example services configuration file:
```toml
# Services configuration.
[[services]]
  name = "skycoin"
  mode = "rest_api"
  scheme = "http"
  port = 8000
  basePath = ""
  authType = "CSRF"
  tokenHeaderKey = "X-CSRF-Token"
  genTokenEndpoint = "/api/v1/csrf"
  tokenKeyFromEndpoint = "csrf_token"

  [services.location]
    location = "localhost"
```

Example metrics for service configuration file:
```toml
serviceNameToMetricsConfPath = [
	{ skycoin = "/home/adacosta/.config/simelo/rextporter/skycoinMetrics.toml" },
	{ wallet = "/home/adacosta/.config/simelo/rextporter/walletMetrics.toml" },
]
```

Example metrics(for skycoin in this case) file configuration file.
```toml
# All metrics to be measured.
[[metrics]]
  name = "seq"
  url = "/api/v1/health"
  httpMethod = "GET"
  path = "/blockchain/head/seq"

  [metrics.options]
    type = "Counter"
    description = "Put a description for this metrics"
```

Example gauge vector metric configuration.
```toml
[[metrics]]
  name = "burn_factor_by_service"
  url = "/api/v1/network/connections"
  httpMethod = "GET"
  path = "/connections"

  [metrics.options]
    type = "Gauge"
    itemPath = "/unconfirmed_verify_transaction/burn_factor"
    description = "I am running since"

  [[metrics.options.labels]]
    name = "ip_port"
    path = "/address"
```