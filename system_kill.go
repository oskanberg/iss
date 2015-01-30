package main

func DeleteFromSlice(slice []*SimpleAgent, index int) []*SimpleAgent {
	copy(slice[index:], slice[index+1:])
	slice[len(slice)-1] = nil
	slice = slice[:len(slice)-1]
	return slice
}

func Kill(population *Population, prey *SimpleAgent) {
	// messy and inefficient, but will have to do for now
	if prey.Family == PREYA {
		for i, agent := range population.TypeA {
			if agent == prey {
				population.TypeADead = append(population.TypeADead, agent)
				population.TypeA = DeleteFromSlice(population.TypeA, i)
			}
		}
	} else {
		for i, agent := range population.TypeB {
			if agent == prey {
				population.TypeBDead = append(population.TypeBDead, agent)
				population.TypeB = DeleteFromSlice(population.TypeB, i)
			}
		}
	}
}
