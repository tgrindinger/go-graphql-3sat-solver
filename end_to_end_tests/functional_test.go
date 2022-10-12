package end_to_end_tests

import (
	"context"
	"testing"
	"time"

	u "github.com/google/uuid"
	"github.com/shurcooL/graphql"

	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

func TestHappyPath(t *testing.T) {
	cases := []struct {
		desc string
		input *model.NewJob
		want *model.Solution
	}{
		{ "returns a valid solution", simpleJob(), simpleSolution() },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			uuid := createJob(t, tc.input)
			solution := findSolution(t, uuid)
			assertSolutionsAreEqual(t, solution, tc.want, uuid)
		})
	}
}

func assertSolutionsAreEqual(t testing.TB, got *FindSolutionResponse, want *model.Solution, uuid u.UUID) {
	if got.Uuid.(string) != uuid.String() {
		t.Errorf("failed to match uuid on job and solution: job '%s' solution '%s'", uuid.String(), got.Uuid.(string))
	}
	if float64(got.Score) != want.Score {
		t.Errorf("failed to match score: got %f want %f", got.Score, want.Score)
	}
	if len(got.Variables) != len(want.Variables) {
		t.Fatalf("failed to match number of variables: got %d want %d", len(got.Variables), len(want.Variables))
	}
	assertVariablesAreEqual(t, got, want)
}

func assertVariablesAreEqual(t testing.TB, got *FindSolutionResponse, want *model.Solution) {
	for _, wantVar := range want.Variables {
		gotVar := findMatchingVariable(t, got, wantVar)
		if gotVar.Value != graphql.Boolean(wantVar.Value) {
			t.Errorf("failed to match value of variables: got %t want %t", gotVar.Value, wantVar.Value)
		}
	}
}

func findMatchingVariable(t testing.TB, got *FindSolutionResponse, wantVar *model.SolvedVariable) FindSolutionResponseVariable {
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
		return FindSolutionResponseVariable{}
}

func findSolution(t testing.TB, uuid u.UUID) *FindSolutionResponse {
	client := graphql.NewClient("http://localhost:8080/query", nil)
	var query struct {
		FindSolution FindSolutionResponse `graphql:"solution(uuid: $uuid)"`
	}
	variables := map[string]interface{}{
		"uuid": uuid.String(),
	}
	err := retry(5 * time.Second, func() error {
		return client.Query(context.Background(), &query, variables)
	})
	if err != nil {
		t.Fatalf("failed to execute findSolution request: %v", err)
	}
	return &query.FindSolution
}

func retry(timeout time.Duration, f func() error) error {
	start := time.Now()
	err := f()
	for err != nil && time.Since(start) < timeout {
		err = f()
	}
	return err
}

type FindSolutionResponse struct {
	Uuid graphql.ID
	Score graphql.Float
	Variables []FindSolutionResponseVariable
}

type FindSolutionResponseVariable struct {
	Name graphql.String
	Value graphql.Boolean
}

func createJob(t testing.TB, job *model.NewJob) u.UUID {
	client := graphql.NewClient("http://localhost:8080/query", nil)
	var query struct {
		CreateJob struct {
			Uuid graphql.ID
		} `graphql:"createJob(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": *job,
	}
	err := client.Mutate(context.Background(), &query, variables)
	if err != nil {
		t.Fatalf("failed to execute createJob request: %v", err)
	}
	uuid, err := u.Parse(query.CreateJob.Uuid.(string))
	if err != nil {
		t.Fatalf("failed to parse uuid: %v", err)
	}
	return uuid
}

func simpleJob() *model.NewJob {
	return &model.NewJob{
		Name: "simplejob",
		Clauses: []*model.NewClause{
			{
				Var1: &model.NewVariable{Name: "var1", Negated: false},
				Var2: &model.NewVariable{Name: "var2", Negated: false},
				Var3: &model.NewVariable{Name: "var3", Negated: false},
			},
		},
	}
}

func simpleSolution() *model.Solution {
	return &model.Solution{
		Score: 1.0,
		Variables: []*model.SolvedVariable{
			{Name: "var1", Value: true},
			{Name: "var2", Value: true},
			{Name: "var3", Value: true},
		},
	}
}
