package main

import "math/rand"

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

func replaceSubspecies(subspecies []*SimpleAgent, dead []*SimpleAgent) {

	var highestFitness int
	for _, agent := range subspecies {
		if agent.Fitness > highestFitness {
			highestFitness = agent.Fitness
		}
	}

	replacements := make([]*SimpleAgent, len(dead))
	for i, _ := range dead {
		replacements[i] = NewRandomSimpleAgent(dead[i].Family)
		parent := selectFitnessProportionate(append(subspecies, dead...), highestFitness)
		*replacements[i].genetics = *parent.genetics
	}

	for _, r := range replacements {
		subspecies = append(subspecies, r)
	}
}

func Evolve(population Population) {
	replaceSubspecies(population.TypeA, population.TypeADead)
	replaceSubspecies(population.TypeB, population.TypeBDead)

	//clear dead
	population.TypeADead = nil
	population.TypeADead = nil
}
