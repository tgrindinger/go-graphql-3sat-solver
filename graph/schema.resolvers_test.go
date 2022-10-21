package graph

import (
	"context"
	"testing"

	u "github.com/google/uuid"

	"github.com/tgrindinger/go-graphql-3sat-solver/graph/factories"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/repositories"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/solvers"
)

type mutationResolverContext struct {
	solver solvers.Solver
	jobRepository *repositories.InMemoryJobRepository
	solutionRepository *repositories.InMemorySolutionRepository
	jobFactory *factories.JobFactory
	solutionFactory *factories.SolutionFactory
	jobDispatcher *JobDispatcher
	resolver *Resolver
	mutationResolver *mutationResolver
	queryResolver *queryResolver
}

func newMutationResolverContext() *mutationResolverContext {
	jobRepository := &repositories.InMemoryJobRepository{}
	solutionRepository := &repositories.InMemorySolutionRepository{}
	jobFactory := &factories.JobFactory{}
	solutionFactory := &factories.SolutionFactory{}
	solver := solvers.NewNaiveSolver(
		solutionFactory,
	)
	jobDispatcher := NewJobDispatcher(
		solver,
		jobRepository,
		solutionRepository,
		jobFactory,
	)
	resolver := &Resolver{
		JobDispatcher: jobDispatcher,
	}
	mutationResolver := &mutationResolver{
		Resolver: resolver,
	}
	queryResolver := &queryResolver{
		Resolver: resolver,
	}
	return &mutationResolverContext{
		solver: solver,
		jobRepository: jobRepository,
		solutionRepository: solutionRepository,
		jobFactory: jobFactory,
		solutionFactory: solutionFactory,
		jobDispatcher: jobDispatcher,
		resolver: resolver,
		mutationResolver: mutationResolver,
		queryResolver: queryResolver,
	}
}

func TestCreateJob(t *testing.T) {
	cases := []struct {
		desc string
		input model.NewJob
		want *model.Job
		err error
	}{
		{ "empty newjob returns empty done job", model.NewJob{}, &model.Job{}, nil },
		{ "newjob with one clause returns done job with one clause", newJobWithOneClause(), jobWithOneClause(), nil },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			mutationResolverContext := newMutationResolverContext()
			job, err := mutationResolverContext.mutationResolver.CreateJob(context.TODO(), tc.input)
			assertJobsAreEqual(t, job, tc.want)
			if err != tc.err {
				t.Errorf("got '%v' want '%v'", err, tc.err)
			}
		})
	}
}

func TestJobWhenGivenInvalidInput(t *testing.T) {
	cases := []struct {
		desc string
		uuid string
		job *model.Job
		err string
	}{
		{ "error on job not found", "b2312d3c-b09d-4d35-9528-a70104c70738", nil, "unable to find job with uuid b2312d3c-b09d-4d35-9528-a70104c70738" },
		{ "error on invalid uuid", "stuff", nil, "invalid UUID length: 5" },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			mutationResolverContext := newMutationResolverContext()
			job, err := mutationResolverContext.queryResolver.Job(context.TODO(), tc.uuid)
			if job != nil {
				t.Fatalf("got a job when should be error")
			}
			if err.Error() != tc.err {
				t.Errorf("got '%v' want '%v'", err, tc.err)
			}
		})
	}
}

func TestJobWhenGivenValidInput(t *testing.T) {
	cases := []struct {
		desc string
		uuid string
		job *model.Job
	}{
		{ "returns found job", uuidOfJobWithKnownUuid(), jobWithKnownUuid() },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			mutationResolverContext := newMutationResolverContext()
			mutationResolverContext.jobRepository.InsertJob(tc.job)
			job, err := mutationResolverContext.queryResolver.Job(context.TODO(), tc.uuid)
			assertJobsAreEqual(t, job, tc.job)
			if err != nil {
				t.Errorf("returned an error")
			}
		})
	}
}

func TestSolutionWhenGivenInvalidInput(t *testing.T) {
	cases := []struct {
		desc string
		uuid string
		solution *model.Solution
		err string
	}{
		{ "error on solution not found", "b2312d3c-b09d-4d35-9528-a70104c70738", nil, "unable to find solution with uuid b2312d3c-b09d-4d35-9528-a70104c70738" },
		{ "error on invalid uuid", "stuff", nil, "invalid UUID length: 5" },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			mutationResolverContext := newMutationResolverContext()
			solution, err := mutationResolverContext.queryResolver.Solution(context.TODO(), tc.uuid)
			if solution != nil {
				t.Fatalf("got a solution when should be error")
			}
			if err.Error() != tc.err {
				t.Errorf("got '%v' want '%v'", err, tc.err)
			}
		})
	}
}

func TestSolutionWhenGivenValidInput(t *testing.T) {
	cases := []struct {
		desc string
		uuid string
		solution *model.Solution
	}{
		{ "returns found solution", uuidOfSolutionWithKnownUuid(), solutionWithKnownUuid() },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			mutationResolverContext := newMutationResolverContext()
			mutationResolverContext.solutionRepository.InsertSolution(tc.solution)
			solution, err := mutationResolverContext.queryResolver.Solution(context.TODO(), tc.uuid)
			assertSolutionsAreEqual(t, solution, tc.solution)
			if err != nil {
				t.Errorf("returned an error")
			}
		})
	}
}

func newJobWithOneClause() model.NewJob {
	return model.NewJob{
		Clauses: []*model.NewClause{
			{
				Var1: &model.NewVariable{ Name: "v1", Negated: true },
				Var2: &model.NewVariable{ Name: "v2", Negated: false },
				Var3: &model.NewVariable{ Name: "v3", Negated: true },
			},
		},
	}
}

func jobWithOneClause() *model.Job {
	return &model.Job{
		Clauses: []*model.Clause{
			{
				Var1: &model.Variable{ Name: "v1", Negated: true },
				Var2: &model.Variable{ Name: "v2", Negated: false },
				Var3: &model.Variable{ Name: "v3", Negated: true },
			},
		},
	}
}

func uuidOfJobWithKnownUuid() string {
	return "e2be8104-4770-44fe-ad38-7f85088700f7"
}

func jobWithKnownUuid() *model.Job {
	uuid, _ := u.Parse(uuidOfJobWithKnownUuid())
	return &model.Job{
		Uuid: uuid,
		Clauses: []*model.Clause{
			{
				Var1: &model.Variable{ Name: "v1", Negated: true },
				Var2: &model.Variable{ Name: "v2", Negated: false },
				Var3: &model.Variable{ Name: "v3", Negated: true },
			},
		},
	}
}

func uuidOfSolutionWithKnownUuid() string {
	return "e2be8104-4770-44fe-ad38-7f85088700f7"
}

func solutionWithKnownUuid() *model.Solution {
	uuid, _ := u.Parse(uuidOfJobWithKnownUuid())
	return &model.Solution{
		Uuid: uuid,
		Variables: []*model.SolvedVariable{
			{ Name: "v1", Value: false },
			{ Name: "v2", Value: false },
			{ Name: "v3", Value: false },
		},
	}
}

func assertJobsAreEqual(t testing.TB, got *model.Job, want *model.Job) {
	if (got == nil) != (want == nil) {
		t.Fatalf("nil expectations violated: got '%t' want '%t'", got == nil, want == nil)
	}
	if want == nil {
		return
	}
	if got.Name != want.Name {
		t.Errorf("wrong Name value: got '%s' want '%s'", got.Name, want.Name)
	}
	if got.Done != want.Done {
		t.Errorf("wrong Done value: got '%t' want '%t'", got.Done, want.Done)
	}
	if len(got.Clauses) != len(want.Clauses) {
		t.Fatalf("wrong number of clauses: got '%d' want '%d'", len(got.Clauses), len(want.Clauses))
	}
	for i := range got.Clauses {
		assertVariablesAreEqual(t, i, 1, got.Clauses[i].Var1, want.Clauses[i].Var1)
		assertVariablesAreEqual(t, i, 2, got.Clauses[i].Var2, want.Clauses[i].Var2)
		assertVariablesAreEqual(t, i, 3, got.Clauses[i].Var3, want.Clauses[i].Var3)
	}
}

func assertSolutionsAreEqual(t testing.TB, got *model.Solution, want *model.Solution) {
	if (got == nil) != (want == nil) {
		t.Fatalf("nil expectations violated: got '%t' want '%t'", got == nil, want == nil)
	}
	if want == nil {
		return
	}
	if got.Score != want.Score {
		t.Errorf("wrong Score value: got %f want %f", got.Score, want.Score)
	}
	if got.Uuid != want.Uuid {
		t.Errorf("wrong Uuid value: got %s want %s", got.Uuid.String(), want.Uuid.String())
	}
	if len(got.Variables) != len(want.Variables) {
		t.Errorf("wrong number of variables: got %d want %d", len(got.Variables), len(want.Variables))
	}
	for i := range got.Variables {
		assertSolvedVariablesAreEqual(t, i, got.Variables[i], want.Variables[i])
	}
}

func assertVariablesAreEqual(t testing.TB, clauseIndex int, varIndex int, got *model.Variable, want *model.Variable) {
	if got.Name != want.Name {
		t.Errorf("wrong name for variable %d in clause %d: got '%s' want '%s'", clauseIndex, varIndex, got.Name, want.Name)
	}
	if got.Negated != want.Negated {
		t.Errorf("wrong negated for variable %d in clause %d: got '%t' want '%t'", clauseIndex, varIndex, got.Negated, want.Negated)
	}
}

func assertSolvedVariablesAreEqual(t testing.TB, index int, got *model.SolvedVariable, want *model.SolvedVariable) {
	if got.Name != want.Name {
		t.Errorf("wrong Name for variable %d: got '%s' want '%s'", index, got.Name, want.Name)
	}
	if got.Value != want.Value {
		t.Errorf("wrong Value for variable %d: got '%t' want '%t'", index, got.Value, want.Value)
	}
}
