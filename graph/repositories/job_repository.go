package repositories

import (
	u "github.com/google/uuid"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

type JobRepository interface {
	FindJob(uuid u.UUID) (*model.Job, error)
	InsertJob(job *model.Job)
	MarkDone(job *model.Job)
}
