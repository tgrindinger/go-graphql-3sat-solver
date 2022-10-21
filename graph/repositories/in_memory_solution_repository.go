package repositories

import (
	"fmt"
	"sync"

	u "github.com/google/uuid"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

type InMemorySolutionRepository struct {
	solutions []*model.Solution
	m sync.RWMutex
}

func (r* InMemorySolutionRepository) FindSolution(uuid u.UUID) (*model.Solution, error) {
	r.m.RLock()
	for _, j := range r.solutions {
		if j.Uuid == uuid {
			r.m.RUnlock()
			return j, nil
		}
	}
	r.m.RUnlock()
	return nil, fmt.Errorf("unable to find solution with uuid %s", uuid.String())
}

func (r* InMemorySolutionRepository) InsertSolution(solutions *model.Solution) {
	r.m.Lock()
	r.solutions = append(r.solutions, solutions)
	r.m.Unlock()
}
