package factories

import "github.com/tgrindinger/go-graphql-3sat-solver/graph/model"

type SolutionFactory struct {
}

func (f *SolutionFactory) ConstructSolution(variables map[string]bool, job *model.Job) *model.Solution {
	return &model.Solution{
		Uuid: job.Uuid,
		Variables: f.packageSolvedVariables(variables),
		Score: f.score(variables, job.Clauses),
	}
}

func (f *SolutionFactory) score(variables map[string]bool, clauses []*model.Clause) float64 {
	correct := 0
	for _, c := range clauses {
		value1 := variables[c.Var1.Name] != c.Var1.Negated
		value2 := variables[c.Var2.Name] != c.Var2.Negated
		value3 := variables[c.Var3.Name] != c.Var3.Negated
		if value1 || value2 || value3 {
			correct++
		}
	}
	return float64(correct) / float64(len(clauses))
}

func (f *SolutionFactory) packageSolvedVariables(variables map[string]bool) []*model.SolvedVariable {
	solvedVariables := []*model.SolvedVariable{}
	for key, value := range variables {
		solvedVariables = append(solvedVariables, &model.SolvedVariable{
			Name: key,
			Value: value,
		})
	}
	return solvedVariables
}
