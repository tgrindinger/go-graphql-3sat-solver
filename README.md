# go-graphql-3sat-solver
A GraphQL API for attempting to generate solutions to the NP-Complete 3SAT problem written in Go

`go run server.go` will start the server.

`go generate ./...` will rebuild generated files if you want to make changes to the graphql schema.

`go test ./graph/...` will run the unit tests. The Go extension for VSCode has nice integration for running individual unit tests.

`go test ./end_to_end_tests/...` will run the end-to-end tests. Make sure the server is already running.

Roadmap:
Integrate a persistent database for storing jobs and solutions.
Replace the naive solver with one based on a genetic algorithm. Perhaps more algorithms to come!
Improve end-to-end tests.
Add user management.
Improve the job management to limit duration time and number of concurrent goroutines.
