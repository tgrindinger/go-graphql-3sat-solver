package solvers

import (
	"math/rand"
	"time"

	"github.com/tgrindinger/go-graphql-3sat-solver/graph/factories"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

type geneticSolver struct {
	maxPopulation int
	maxTime time.Duration
	solutionFactory *factories.SolutionFactory
	populationGenerator *PopulationGenerator
	randomFactory factories.RandomFactory
}

func NewGeneticSolver(
	maxPopulation int,
	maxTime time.Duration,
	solutionFactory *factories.SolutionFactory,
	populationGenerator *PopulationGenerator,
	randomFactory factories.RandomFactory,
) *geneticSolver {
	return &geneticSolver{
		maxPopulation: maxPopulation,
		maxTime: maxTime,
		solutionFactory: solutionFactory,
		populationGenerator: populationGenerator,
		randomFactory: randomFactory,
	}
}

func (s *geneticSolver) Solve(job *model.Job) *model.Solution {
	start := time.Now()
	random := s.randomFactory.Build()
	variables := job.Variables()
	population := s.populationGenerator.generatePopulation(s.maxPopulation, variables)
	bestMember, cycles := s.start(job, population, random)
	elapsed := time.Since(start)
	return s.solutionFactory.ConstructSolution(bestMember, job, cycles, elapsed)
}

func (s *geneticSolver) start(job *model.Job, population population, random *rand.Rand) (member, int) {
	cycles := 0
	var bestMember member
	if len(population) < s.maxPopulation {
		bestMember, _ = population.best(job)
	} else {
		bestMember, cycles = s.evolve(job, population, random)
	}
	return bestMember, cycles
}

func (s *geneticSolver) evolve(job *model.Job, population population, random *rand.Rand) (member, int) {
	start := time.Now()
	cycles := 0
	bestMember, bestScore := population.best(job)
	for time.Since(start) < s.maxTime && bestScore < 1.0 {
		population = s.reproduce(job, population, random)
		bestMember, bestScore = population.best(job)
		cycles++
	}
	return bestMember, cycles
}

func (s *geneticSolver) reproduce(job *model.Job, _population population, random *rand.Rand) population {
	newPop := population{}
	rankedPop := _population.rank(job)
	for i := 0; i < len(_population); i++ {
		child := s.breed(rankedPop, _population, random)
		for newPop.memberExists(child) {
			child = s.breed(rankedPop, _population, random)
		}
		newPop = append(newPop, child)
	}
	return newPop
}

func (s *geneticSolver) breed(memberRanks memberRanks, population population, random *rand.Rand) member {
	parent1 := memberRanks.selectMember(random)
	parent2 := memberRanks.selectMember(random)
	for parent1 == parent2 {
		parent2 = memberRanks.selectMember(random)
	}
	return population[parent1].crossOver(population[parent2], random)
}
