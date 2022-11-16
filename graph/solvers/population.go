package solvers

import (
	"math/rand"
	"sort"

	"github.com/samber/lo"
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
	for index, rank := range ranks {
		normTotal += rank.fitness / total
		ranks[index].fitness = normTotal
	}
	return ranks
}

func (m member) crossOver(_member member, random *rand.Rand) member {
	child := member{}
	keys := lo.Keys(m)
	sort.Strings(keys)
	for _, k := range keys {
		if random.Intn(2) == 1 {
			child[k] = m[k]
		} else {
			child[k] = _member[k]
		}
	}
	return child
}

func (m member) matches(member member) bool {
	return lo.EveryBy(lo.Keys(m), func(k string) bool {
		return m[k] == member[k]
	})
}

func (m member) mutate(random *rand.Rand) member {
	rate := 0.01
	mutated := member{}
	keys := lo.Keys(m)
	sort.Strings(keys)
	for _, k := range keys {
		if random.Float64() < rate {
			mutated[k] = !m[k]
		} else {
			mutated[k] = m[k]
		}
	}
	return mutated
}
