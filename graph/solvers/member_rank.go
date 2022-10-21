package solvers

import "math/rand"

type memberRank struct {
	index   int
	fitness float64
}

type memberRanks []memberRank

func (m memberRanks) Len() int {
	return len(m)
}

func (m memberRanks) Less(i, j int) bool {
	return m[i].fitness < m[j].fitness
}

func (m memberRanks) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m memberRanks) selectMember(random *rand.Rand) int {
	target := random.Float64()
	for _, rank := range m {
		if rank.fitness > target {
			return rank.index
		}
	}
	return len(m) - 1
}
