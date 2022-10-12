package solvers

import "github.com/tgrindinger/go-graphql-3sat-solver/graph/model"

type Solver interface {
	Solve(job *model.Job) *model.Solution
}
