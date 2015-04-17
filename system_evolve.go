package main

import (
	"math/rand"

	"github.com/oskanberg/go-vector"
)

func selectFitnessProportionate(sample []*SimpleAgent, highestFitness int) *SimpleAgent {
	if highestFitness == 0 {
		return sample[rand.Intn(len(sample))]
	}
	var index int
	for {
		index = rand.Intn(len(sample))
		if rand.Float64() < float64(sample[index].Fitness)/float64(highestFitness) {
			return sample[index]
		}
	}
}

func replaceSubspecies(subspecies []*SimpleAgent, dead []*SimpleAgent) []*SimpleAgent {

	var highestFitness int
	for _, agent := range subspecies {
		if agent.Fitness > highestFitness {
			highestFitness = agent.Fitness
		}
	}

	replacements := make([]*SimpleAgent, len(dead))
	for i := range dead {
		replacements[i] = NewRandomSimpleAgent(dead[i].Family)
		parent := selectFitnessProportionate(append(subspecies, dead...), highestFitness)
		*replacements[i].Genetics = *parent.Genetics.Mutated()
	}

	for _, r := range replacements {
		subspecies = append(subspecies, r)
	}

	return subspecies
}

var typeADead, typeBDead []int

func Evolve(population *Population) {
	// randomise location
	for _, agent := range population.TypeA {
		agent.Position = *vector.NewVector2d(rand.Float64()*SIMULATION_SPACE_SIZE, rand.Float64()*SIMULATION_SPACE_SIZE)
		agent.Velocity = *vector.NewRandomUnitVector()
		agent.VelocityNext = agent.Velocity
	}
	for _, agent := range population.TypeB {
		agent.Position = *vector.NewVector2d(rand.Float64()*SIMULATION_SPACE_SIZE, rand.Float64()*SIMULATION_SPACE_SIZE)
		agent.Velocity = *vector.NewRandomUnitVector()
		agent.VelocityNext = agent.Velocity
	}

	typeADead = append(typeADead, len(population.TypeADead))
	typeBDead = append(typeBDead, len(population.TypeBDead))

	population.TypeA = replaceSubspecies(population.TypeA, population.TypeADead)
	population.TypeB = replaceSubspecies(population.TypeB, population.TypeBDead)

	//clear dead
	population.TypeADead = nil
	population.TypeBDead = nil
}
