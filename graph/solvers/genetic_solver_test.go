package solvers

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	u "github.com/google/uuid"

	"github.com/tgrindinger/go-graphql-3sat-solver/graph/factories"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

func TestSolve(t *testing.T) {
	cases := []struct {
		desc string
		job *model.Job
		want *model.Solution
	}{
		{ "single clause is fully solved", singleClauseJob(), singleClauseSolution() },
		{ "two clauses are fully solved", twoClauseJob(), twoClauseSolution() },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			// arrange
			tc.want.Uuid = tc.job.Uuid
			maxPopulation := 10
			maxTime, _ := time.ParseDuration("1ms")
			factory := &factories.SolutionFactory{}
			randomFactory := &factories.ZeroRandomFactory{}
			populationGenerator := &PopulationGenerator{randomFactory: randomFactory}
			sut := NewGeneticSolver(maxPopulation, maxTime, factory, populationGenerator, randomFactory)

			// act
			got := sut.Solve(tc.job)

			// assert
			assertSolutionsAreEqual(t, got, tc.want)
		})
	}
}

func TestBigSolvable(t *testing.T) {
	randomFactory := &factories.ZeroRandomFactory{}
	cases := []struct {
		desc string
		job *model.Job
		want *model.Solution
	}{
		{ "big solvable job is fully solved", bigSolvableJob(randomFactory.Build()), bigSolvableSolution() },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			// arrange
			tc.want.Uuid = tc.job.Uuid
			maxPopulation := 10
			maxTime, _ := time.ParseDuration("10s")
			factory := &factories.SolutionFactory{}
			populationGenerator := &PopulationGenerator{randomFactory: randomFactory}
			sut := NewGeneticSolver(maxPopulation, maxTime, factory, populationGenerator, randomFactory)

			// act
			got := sut.Solve(tc.job)

			// assert
			assertBigSolutionsAreEqual(t, got, tc.want)
			fmt.Printf("cycles: %d\n", got.Cycles)
		})
	}
}

func TestBigUnsolvable(t *testing.T) {
	randomFactory := &factories.ZeroRandomFactory{}
	cases := []struct {
		desc string
		job *model.Job
		want *model.Solution
	}{
		{ "big unsolvable job expires", bigUnsolvableJob(randomFactory.Build()), bigUnsolvableSolution() },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			// arrange
			tc.want.Uuid = tc.job.Uuid
			maxPopulation := 10
			maxTime, _ := time.ParseDuration("1s")
			factory := &factories.SolutionFactory{}
			populationGenerator := &PopulationGenerator{randomFactory: randomFactory}
			sut := NewGeneticSolver(maxPopulation, maxTime, factory, populationGenerator, randomFactory)

			// act
			got := sut.Solve(tc.job)

			// assert
			assertBigSolutionsAreEqual(t, got, tc.want)
		})
	}
}

func assertSolutionsAreEqual(t testing.TB, got *model.Solution, want *model.Solution) {
	if got.Uuid.String() != want.Uuid.String() {
		t.Errorf("failed to match uuid on solutions: got '%s' want '%s'", got.Uuid.String(), want.Uuid.String())
	}
	if got.Score != want.Score {
		t.Errorf("failed to match score: got %f want %f", got.Score, want.Score)
	}
	if len(got.Variables) != len(want.Variables) {
		t.Fatalf("failed to match number of variables: got %d want %d", len(got.Variables), len(want.Variables))
	}
	assertVariablesAreEqual(t, got, want)
}

func assertBigSolutionsAreEqual(t testing.TB, got *model.Solution, want *model.Solution) {
	if got.Uuid.String() != want.Uuid.String() {
		t.Errorf("failed to match uuid on solutions: got '%s' want '%s'", got.Uuid.String(), want.Uuid.String())
	}
	if (want.Score == 1.0 && got.Score < 1.0) ||
	    (want.Score < 1.0 && got.Score == 1.0) {
		t.Errorf("failed to match score: got %f want %f", got.Score, want.Score)
	}
	if len(got.Variables) != len(want.Variables) {
		t.Fatalf("failed to match number of variables: got %d want %d", len(got.Variables), len(want.Variables))
	}
}

func assertVariablesAreEqual(t testing.TB, got *model.Solution, want *model.Solution) {
	for _, wantVar := range want.Variables {
		gotVar := findMatchingVariable(t, got, wantVar)
		if gotVar.Value != wantVar.Value {
			t.Errorf("failed to match value of variable '%s': got %t want %t", gotVar.Name, gotVar.Value, wantVar.Value)
		}
	}
}

func findMatchingVariable(t testing.TB, got *model.Solution, wantVar *model.SolvedVariable) *model.SolvedVariable {
		for _, gotVar := range got.Variables {
			if wantVar.Name == string(gotVar.Name) {
				return gotVar
			}
		}
		gotNames := []string{}
		for _, gotv := range got.Variables {
			gotNames = append(gotNames, string(gotv.Name))
		}
		t.Fatalf("failed to find matching variable in solution: got '%v' want '%s'", gotNames, wantVar.Name)
		return &model.SolvedVariable{}
}

func singleClauseJob() *model.Job {
	return &model.Job{
		Name: u.NewString(),
		Clauses: []*model.Clause{
			{
				Var1: &model.Variable{ Name: "v1", Negated: true },
				Var2: &model.Variable{ Name: "v2", Negated: false },
				Var3: &model.Variable{ Name: "v3", Negated: true },
			},
		},
	}
}

func singleClauseSolution() *model.Solution {
	return &model.Solution{
		Score: 1.0,
		Variables: []*model.SolvedVariable{
			{ Name: "v1", Value: true },
			{ Name: "v2", Value: true },
			{ Name: "v3", Value: true },
		},
	}
}

func twoClauseJob() *model.Job {
	return &model.Job{
		Name: u.NewString(),
		Clauses: []*model.Clause{
			{
				Var1: &model.Variable{ Name: "v1", Negated: true },
				Var2: &model.Variable{ Name: "v2", Negated: true },
				Var3: &model.Variable{ Name: "v3", Negated: true },
			},
			{
				Var1: &model.Variable{ Name: "v1", Negated: false },
				Var2: &model.Variable{ Name: "v2", Negated: false },
				Var3: &model.Variable{ Name: "v3", Negated: false },
			},
		},
	}
}

func twoClauseSolution() *model.Solution {
	return &model.Solution{
		Score: 1.0,
		Variables: []*model.SolvedVariable{
			{ Name: "v1", Value: false },
			{ Name: "v2", Value: false },
			{ Name: "v3", Value: true },
		},
	}
}

func bigSolvableJob(random *rand.Rand) *model.Job {
	job := &model.Job{
		Name: u.NewString(),
		Clauses: []*model.Clause{},
	}
	for i := 0; i < 100; i++ {
		job.Clauses = append(job.Clauses, &model.Clause{
			Var1: &model.Variable{ Name: fmt.Sprintf("%d", i * 3 + 0), Negated: random.Intn(2) == 0 },
			Var2: &model.Variable{ Name: fmt.Sprintf("%d", i * 3 + 1), Negated: random.Intn(2) == 0 },
			Var3: &model.Variable{ Name: fmt.Sprintf("%d", i * 3 + 2), Negated: random.Intn(2) == 0 },
		})
	}
	return job
}

func bigSolvableSolution() *model.Solution {
	solution := &model.Solution{
		Score: 1.0,
		Variables: []*model.SolvedVariable{},
	}
	for i := 0; i < 300; i++ {
		solution.Variables = append(solution.Variables, &model.SolvedVariable{
			Name: fmt.Sprintf("%d", i),
			Value: true,
		})
	}
	return solution
}

func bigUnsolvableJob(random *rand.Rand) *model.Job {
	job := &model.Job{
		Name: u.NewString(),
		Clauses: []*model.Clause{{
				Var1: &model.Variable{ Name: "1000", Negated: true },
				Var2: &model.Variable{ Name: "1000", Negated: true },
				Var3: &model.Variable{ Name: "1000", Negated: true },
			},{
				Var1: &model.Variable{ Name: "1000", Negated: false },
				Var2: &model.Variable{ Name: "1000", Negated: false },
				Var3: &model.Variable{ Name: "1000", Negated: false },
			},
		},
	}
	for i := 0; i < 100; i++ {
		job.Clauses = append(job.Clauses, &model.Clause{
			Var1: &model.Variable{ Name: fmt.Sprintf("%d", i * 3 + 0), Negated: random.Intn(2) == 0 },
			Var2: &model.Variable{ Name: fmt.Sprintf("%d", i * 3 + 1), Negated: random.Intn(2) == 0 },
			Var3: &model.Variable{ Name: fmt.Sprintf("%d", i * 3 + 2), Negated: random.Intn(2) == 0 },
		})
	}
	return job
}

func bigUnsolvableSolution() *model.Solution {
	solution := &model.Solution{
		Score: 0.0,
		Variables: []*model.SolvedVariable{},
	}
	for i := 0; i < 301; i++ {
		solution.Variables = append(solution.Variables, &model.SolvedVariable{
			Name: fmt.Sprintf("%d", i),
			Value: false,
		})
	}
	return solution
}
