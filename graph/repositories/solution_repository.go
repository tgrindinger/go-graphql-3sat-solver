package repositories

import (
	u "github.com/google/uuid"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

type SolutionRepository interface {
	FindSolution(uuid u.UUID) (*model.Solution, error)
	InsertSolution(solution *model.Solution)
}
