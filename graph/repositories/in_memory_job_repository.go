package repositories

import (
	"fmt"
	"sync"

	u "github.com/google/uuid"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

type InMemoryJobRepository struct {
	jobs []*model.Job
	m sync.Mutex
}

func (r* InMemoryJobRepository) FindJob(uuid u.UUID) (*model.Job, error) {
	r.m.Lock()
	for _, j := range r.jobs {
		if j.Uuid == uuid {
			r.m.Unlock()
			return j, nil
		}
	}
	r.m.Unlock()
	return nil, fmt.Errorf("unable to find job with uuid %s", uuid.String())
}

func (r* InMemoryJobRepository) InsertJob(job *model.Job) {
	r.m.Lock()
	r.jobs = append(r.jobs, job)
	r.m.Unlock()
}

func (r* InMemoryJobRepository) MarkDone(job *model.Job) {
	r.m.Lock()
	job.Done = true
	r.m.Unlock()
}
