# Changelog 
All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](http://www.keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](http://www.semver.org/spec/v2.0.0.html).

[Unreleased](https://github.com/skycoin/skycoin/compare/master...develop)
- Exporting configured metric under the '/metrics' endpoint.


## [0.0.2](https://github.com/simelo/rexporter/releases...) 2019-01-25

### Added
 * [\#19](https://github.com/simelo/rextporter/issues/19)
   - Multiple services definition to be monitored.

 * [\#18](https://github.com/simelo/rextporter/issues/18)
   - Define a service type, `apiRest` (get values trough API and make this the metrics as defined in metric conf), or `proxy` (forward some exposed metrics with the original metric name changed with the service name as prefix).

 * [\#14](https://github.com/simelo/rextporter/issues/14)
   - Use a default configuration file path and initialize it if not exist.
   
 * [\#13](https://github.com/simelo/rextporter/issues/13)
   - Be able to read service configuration from file.
