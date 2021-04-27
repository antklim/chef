# ramen

Ramen is a layout and components generator for Go projects.

Tastes:
- pkg - package
- app - application
  - http
  - grpc

structure:
/cmd
/cmd/main.go

/internal
/internal/app - application source code

/internal/adapter - adapters from/to app structures (for example from http-request to app struct)
/internal/provider - external services providers/clients
/internal/server - an application server
/internal/server/http - http routes and handlers definitions 
/internal/server/grpc - grpc message handlers
/test - contains testing tools

By default it creates a template for application with http server
Commands:
init - inits a new project
add <component> - adds a component

Options:
--name, -n - project name
--root, -r - project root directory
--taste, -t - pkg, app
--server, -s - http

TODO: Add CI/CD pipeline
  - test and coverage
  - linting with golangci