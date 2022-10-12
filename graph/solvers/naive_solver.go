package solvers

import (
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/factories"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

type naiveSolver struct {
	solutionFactory *factories.SolutionFactory
}

func NewNaiveSolver(
	solutionFactory *factories.SolutionFactory,
) *naiveSolver {
	return &naiveSolver{
		solutionFactory: solutionFactory,
	}
}

func (s *naiveSolver) Solve(job *model.Job) *model.Solution {
	variables := map[string]bool{}
	for _, c := range job.Clauses {
		variables[c.Var1.Name] = true
		variables[c.Var2.Name] = true
		variables[c.Var3.Name] = true
	}
	return s.solutionFactory.ConstructSolution(variables, job)
}
