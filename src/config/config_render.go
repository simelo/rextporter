package config

const mainConfigFileContentTemplate = `
serviceConfigTransport = "file" # "file" | "consulCatalog"
# render a template with a portable path
serviceConfigPath = "{{.ServiceConfigPath}}"
metricsConfigPath = "{{.MetricsConfigPath}}"
`

const serviceConfigFileContentTemplate = `
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
`
const metricsConfigFileContentTemplate = `
# All metrics to be measured.
[[metrics]]
  name = "seq"
  url = "/api/v1/health.json"
  httpMethod = "GET"
  path = "/blockchain/head/seq"

  [metrics.options]
    type = "Counter"
    description = "I am running since"

# [[metrics]]
#   name = "openConnections"
#   url = "/api/v1/network/connections"
#   httpMethod = "GET"
#   path = "/"

#   [metrics.options]
#     type = "Histogram"
#     description = "Connections ammount"

#   [metrics.histogramOptions]
#     buckets = [
#       1,
#       2, 
#       3
#     ]




# TODO(denisacostaq@gmail.com):
# if you refer(under "metrics_for_host") to a not previously defined host or metric it will be raise an error and the process will not start
# if in all your definition you not use some host or metric the process will raise a warning and the process will start normally.
`
