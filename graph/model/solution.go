package model

import (
	"time"

	"github.com/google/uuid"
)

type Solution struct {
	Uuid      uuid.UUID         `json:"uuid"`
	Variables []*SolvedVariable `json:"variables"`
	Score     float64           `json:"score"`
	Cycles    int               `json:"cycles"`
	Elapsed   time.Duration     `json:"elapsed"`
}
