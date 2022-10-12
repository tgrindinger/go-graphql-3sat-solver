package model

import "github.com/google/uuid"

type Solution struct {
	Uuid      uuid.UUID   `json:"uuid"`
	Variables []*SolvedVariable `json:"variables"`
	Score     float64     `json:"score"`
}
