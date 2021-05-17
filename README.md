# chef

[![codecov](https://codecov.io/gh/antklim/chef/branch/master/graph/badge.svg?token=EMWCS55TZR)](https://codecov.io/gh/antklim/chef)

Chef is a layout and components generator for Go projects.

Categories:
- pkg - package
- app - application
  - http
  - grpc

App/Http structure:
- /app - application source code
- /adapter - adapters from/to app structures (for example from http-request to app struct)
- /handler - http routes and handlers
- /provider - external services providers/clients
- /test - contains testing tools
- main.go

Others:
- /cmd
- /cmd/main.go

- /internal
- /internal/app - application source code

- /internal/adapter - adapters from/to app structures (for example from http-request to app struct)
- /internal/provider - external services providers/clients
- /internal/server - an application server
- /internal/server/http - http routes and handlers definitions
- /internal/server/grpc - grpc message handlers
- /test - contains testing tools

By default it creates a template for application with http server
Commands:
init - inits a new project
add <component> - adds a component

Options:
--name, -n - project name
--root, -r - project root directory
--category, -c - pkg, app
--server, -s - http, grpc
