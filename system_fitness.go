package main

func UpdateFitness(agents Population) {
	for _, agent := range agents.TypeA {
		agent.Fitness++
	}
	for _, agent := range agents.TypeB {
		agent.Fitness++
	}
}
