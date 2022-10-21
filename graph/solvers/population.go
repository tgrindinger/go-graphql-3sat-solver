package solvers

import (
	"math/rand"
	"sort"

	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

type member map[string]bool
type population []member

func (p population) memberExists(newMember member) bool {
	for _, member := range p {
		if newMember.matches(member) {
			return true
		}
	}
	return false
}

func (p population) best(job *model.Job) (member, float64) {
	bestScore := 0.0
	bestMember := p[0]
	for _, member := range p {
		score := job.Score(member)
		if score > bestScore {
			bestScore = score
			bestMember = member
		}
	}
	return bestMember, bestScore
}

func (p population) rank(job *model.Job) memberRanks {
	ranks := memberRanks{}
	total := 0.0
	for index, member := range p {
		score := job.Score(member)
		ranks = append(ranks, memberRank{index: index, fitness: score})
		total += score
	}
	sort.Sort(sort.Reverse(ranks))
	normTotal := 0.0
	for _, rank := range ranks {
		normTotal += rank.fitness / total
		rank.fitness = normTotal
	}
	return ranks
}

func (m member) crossOver(_member member, random *rand.Rand) member {
	child := member{}
	for k, v := range m {
		if random.Intn(2) == 1 {
			child[k] = v
		} else {
			child[k] = _member[k]
		}
	}
	return child
}

func (m member) matches(member member) bool {
	for k, v := range m {
		if v != member[k] {
			return false
		}
	}
	return true
}
