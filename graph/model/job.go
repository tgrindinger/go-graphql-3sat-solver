package model

import (
	"sort"

	"github.com/google/uuid"
)

type Job struct {
	Name    string    `json:"name"`
	Clauses []*Clause `json:"clauses"`
	Done    bool      `json:"done"`
	Uuid    uuid.UUID `json:"uuid"`
}

func (j *Job) Variables() []string {
	variables := map[string]bool{}
	for _, c := range j.Clauses {
		variables[c.Var1.Name] = true
		variables[c.Var2.Name] = true
		variables[c.Var3.Name] = true
	}
	return j.keys(variables)
}

func (j *Job) keys(variables map[string]bool) []string {
	keys := make([]string, 0, len(variables))
	for k := range variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (j *Job) Score(variables map[string]bool) float64 {
	correct := 0
	for _, clause := range j.Clauses {
		if clause.satisfied(variables) {
			correct++
		}
	}
	return float64(correct) / float64(len(j.Clauses))
}

func (c *Clause) satisfied(variables map[string]bool) bool {
	return (variables[c.Var1.Name] != c.Var1.Negated) ||
		(variables[c.Var2.Name] != c.Var2.Negated) ||
		(variables[c.Var3.Name] != c.Var3.Negated)
}
