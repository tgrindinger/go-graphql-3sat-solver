package model

import (
	"github.com/google/uuid"
)

type Job struct {
	Name    string    `json:"name"`
	Clauses []*Clause `json:"clauses"`
	Done    bool      `json:"done"`
	Uuid    uuid.UUID `json:"uuid"`
}
