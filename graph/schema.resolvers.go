package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	u "github.com/google/uuid"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/generated"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

// Clauses is the resolver for the clauses field.
func (r *jobResolver) Clauses(ctx context.Context, obj *model.Job) ([]*model.Clause, error) {
	return obj.Clauses, nil
}

// UUID is the resolver for the uuid field.
func (r *jobResolver) UUID(ctx context.Context, obj *model.Job) (string, error) {
	return obj.Uuid.String(), nil
}

// CreateJob is the resolver for the createJob field.
func (r *mutationResolver) CreateJob(ctx context.Context, input model.NewJob) (*model.Job, error) {
	return r.JobDispatcher.DispatchJob(&input), nil
}

// Job is the resolver for the job field.
func (r *queryResolver) Job(ctx context.Context, uuid string) (*model.Job, error) {
	actualUuid, err := u.Parse(uuid)
	if err != nil {
		return nil, err
	}
	return r.JobDispatcher.FindJob(actualUuid)
}

// Solution is the resolver for the solution field.
func (r *queryResolver) Solution(ctx context.Context, uuid string) (*model.Solution, error) {
	actualUuid, err := u.Parse(uuid)
	if err != nil {
		return nil, err
	}
	return r.JobDispatcher.FindSolution(actualUuid)
}

// UUID is the resolver for the uuid field.
func (r *solutionResolver) UUID(ctx context.Context, obj *model.Solution) (string, error) {
	return obj.Uuid.String(), nil
}

// Variables is the resolver for the variables field.
func (r *solutionResolver) Variables(ctx context.Context, obj *model.Solution) ([]*model.SolvedVariable, error) {
	return obj.Variables, nil
}

// Elapsed is the resolver for the elapsed field.
func (r *solutionResolver) Elapsed(ctx context.Context, obj *model.Solution) (int, error) {
	return int(obj.Elapsed.Milliseconds()), nil
}

// Job returns generated.JobResolver implementation.
func (r *Resolver) Job() generated.JobResolver { return &jobResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Solution returns generated.SolutionResolver implementation.
func (r *Resolver) Solution() generated.SolutionResolver { return &solutionResolver{r} }

type jobResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type solutionResolver struct{ *Resolver }
