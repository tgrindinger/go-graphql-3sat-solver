package factories

import (
	u "github.com/google/uuid"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

type JobFactory struct {
}

func (f *JobFactory) CreateJob(newJob *model.NewJob) *model.Job {
	job := &model.Job{
		Name:    newJob.Name,
		Clauses: []*model.Clause{},
		Done:    false,
		Uuid:    u.New(),
	}
	for _, clause := range newJob.Clauses {
		job.Clauses = append(job.Clauses, createClause(clause))
	}
	return job
}

func createClause(clause *model.NewClause) *model.Clause {
	return &model.Clause{
		Var1: createVariable(clause.Var1),
		Var2: createVariable(clause.Var2),
		Var3: createVariable(clause.Var3),
	}
}

func createVariable(variable *model.NewVariable) *model.Variable {
	return &model.Variable{
		Negated: variable.Negated,
		Name:    variable.Name,
	}
}
