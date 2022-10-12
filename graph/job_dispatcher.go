package graph

import (
	"github.com/google/uuid"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/factories"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/repositories"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/solvers"
)

type JobDispatcher struct {
	solver solvers.Solver
	jobRepository repositories.JobRepository
	solutionRepository repositories.SolutionRepository
	jobFactory *factories.JobFactory
}

func NewJobDispatcher(
	solver solvers.Solver,
	jobRepository repositories.JobRepository,
	solutionRepository repositories.SolutionRepository,
	jobFactory *factories.JobFactory,
) *JobDispatcher {
	return &JobDispatcher{
		solver: solver,
		jobRepository: jobRepository,
		solutionRepository: solutionRepository,
		jobFactory: jobFactory,
	}
}

func (d *JobDispatcher) DispatchJob(newJob *model.NewJob) *model.Job {
	job := d.jobFactory.CreateJob(newJob)
	d.jobRepository.InsertJob(job)
	go d.dispatchJobAsync(job)
	return job
}

func (d *JobDispatcher) dispatchJobAsync(job *model.Job) {
	solution := d.solver.Solve(job)
	d.solutionRepository.InsertSolution(solution)
	d.jobRepository.MarkDone(job)
}

func (d *JobDispatcher) FindJob(uuid uuid.UUID) (*model.Job, error) {
	return d.jobRepository.FindJob(uuid)
}

func (d *JobDispatcher) FindSolution(uuid uuid.UUID) (*model.Solution, error) {
	return d.solutionRepository.FindSolution(uuid)
}
