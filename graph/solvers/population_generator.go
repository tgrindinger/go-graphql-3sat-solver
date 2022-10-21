package solvers

import (
	"math"
	"math/rand"

	"github.com/tgrindinger/go-graphql-3sat-solver/graph/factories"
)

type PopulationGenerator struct {
	randomFactory factories.RandomFactory
}

func NewPopulationGenerator(randomFactory factories.RandomFactory) *PopulationGenerator {
	return &PopulationGenerator{randomFactory: randomFactory}
}

func (g *PopulationGenerator) generatePopulation(maxPopulation int, names []string) population {
	population := g.generateBaseMembers(names)
	target := int(math.Min(float64(maxPopulation), math.Pow(2, float64(len(names)))))
	random := g.randomFactory.Build()
	for i := len(population); i < target; i++ {
		member := g.generateMember(names, random)
		for population.memberExists(member) {
			member = g.generateMember(names, random)
		}
		population = append(population, member)
	}
	return population
}

func (g *PopulationGenerator) generateMember(names []string, random *rand.Rand) member {
	member := member{}
	for _, name := range names {
		member[name] = random.Intn(2) == 1
	}
	return member
}

func (g *PopulationGenerator) generateBaseMembers(names []string) population {
	positiveMember := member{}
	negativeMember := member{}
	for _, name := range names {
		positiveMember[name] = true
		negativeMember[name] = false
	}
	return population{positiveMember, negativeMember}
}
