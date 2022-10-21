package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/factories"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/generated"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/repositories"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/solvers"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolver := buildResolver()
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func buildResolver() *graph.Resolver {
	jobRepository := &repositories.InMemoryJobRepository{}
	solutionRepository := &repositories.InMemorySolutionRepository{}
	jobFactory := &factories.JobFactory{}
	solutionFactory := &factories.SolutionFactory{}
	// solver := solvers.NewNaiveSolver(solutionFactory)
	duration, _ := time.ParseDuration("10s")
	randomFactory := &factories.TimeRandomFactory{}
	populationGenerator := solvers.NewPopulationGenerator(randomFactory)
	solver := solvers.NewGeneticSolver(10, duration, solutionFactory, populationGenerator, randomFactory)
	return &graph.Resolver{
		JobDispatcher: graph.NewJobDispatcher(
			solver,
			jobRepository,
			solutionRepository,
			jobFactory,
		),
	}
}
