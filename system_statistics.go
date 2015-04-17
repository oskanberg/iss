package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GeneticsRecord struct {
	Attraction        FloatRecord
	AttractionStdDev  FloatRecord
	Orientation       FloatRecord
	OrientationStdDev FloatRecord
}

type FloatRecord struct {
	Other float64 `json:"o"`
	Same  float64 `json:"s"`
}

var geneticsA []GeneticsRecord = make([]GeneticsRecord, 0)
var geneticsB []GeneticsRecord = make([]GeneticsRecord, 0)

var nearbyA []FloatRecord = make([]FloatRecord, 0)
var nearbyB []FloatRecord = make([]FloatRecord, 0)

func meanNearby(subpopulation []*SimpleAgent) (float64, float64) {
	if len(subpopulation) == 0 {
		return 0, 0
	}

	var same, other float64
	for _, agent := range subpopulation {
		same += float64(agent.nearbySame)
		other += float64(agent.nearbyOther)
	}

	same /= float64(len(subpopulation))
	other /= float64(len(subpopulation))
	return same, other
}

func totalOrientation(subpopulation []*SimpleAgent) (float64, float64) {
	if len(subpopulation) == 0 {
		return 0, 0
	}

	var same, other float64
	for _, agent := range subpopulation {
		same += float64(agent.Genetics.SameSpecies.orientationRadius)
		other += float64(agent.Genetics.OtherSpecies.orientationRadius)
	}

	return same, other
}

func totalAttraction(subpopulation []*SimpleAgent) (float64, float64) {
	if len(subpopulation) == 0 {
		return 0, 0
	}

	var same, other float64
	for _, agent := range subpopulation {
		same += float64(agent.Genetics.SameSpecies.attractionRadius)
		other += float64(agent.Genetics.OtherSpecies.attractionRadius)
	}

	return same, other
}

func getAverageGenetics(subpopulation []*SimpleAgent, dead []*SimpleAgent) GeneticsRecord {
	var same, dSame, other, dOther, total float64
	same, other = totalAttraction(subpopulation)
	dSame, dOther = totalAttraction(dead)
	total = float64(len(subpopulation) + len(dead))
	attraction := FloatRecord{
		Other: (other + dOther) / total,
		Same:  (same + dSame) / total,
	}

	same, other = totalOrientation(subpopulation)
	dSame, dOther = totalOrientation(dead)
	orientation := FloatRecord{
		Other: (other + dOther) / total,
		Same:  (same + dSame) / total,
	}
	r := GeneticsRecord{
		Attraction:  attraction,
		Orientation: orientation,
	}
	return r
}

var iterationsRecorded int

func RecordGenetics(population Population) {
	record := getAverageGenetics(population.TypeA, population.TypeADead)
	geneticsA = append(geneticsA, record)
	record = getAverageGenetics(population.TypeB, population.TypeBDead)
	geneticsB = append(geneticsB, record)
}

func RecordNearby(population Population) {
	nearbySame, nearbyOther := meanNearby(population.TypeA)
	nearbyA = append(nearbyA, FloatRecord{
		Same:  nearbySame,
		Other: nearbyOther,
	})

	nearbySame, nearbyOther = meanNearby(population.TypeB)
	nearbyB = append(nearbyB, FloatRecord{
		Same:  nearbySame,
		Other: nearbyOther,
	})
}

type Stat struct {
	Mean     float64
	Variance float64
	Size     float64
}

var typeAFitness, typeBFitness []Stat

func RecordFitness(population Population) {
	// allFitness := make([]float64, len(population.TypeA)+len(population.TypeADead))
	// i := 0
	var total int
	for _, agent := range population.TypeA {
		// allFitness[i] = agent.Fitness
		// i++
		total += agent.Fitness
	}
	for _, agent := range population.TypeADead {
		// allFitness[i] = agent.Fitness
		// i++
		total += agent.Fitness
	}
	size := float64(len(population.TypeA) + len(population.TypeADead))
	s := Stat{
		Mean:     float64(total) / size,
		Size:     size,
		Variance: 0,
	}
	typeAFitness = append(typeAFitness, s)

	// allFitness := make([]float64, len(population.TypeB)+len(population.TypeBDead))
	// i := 0
	total = 0
	for _, agent := range population.TypeB {
		// allFitness[i] = agent.Fitness
		// i++
		total += agent.Fitness
	}
	for _, agent := range population.TypeBDead {
		// allFitness[i] = agent.Fitness
		// i++
		total += agent.Fitness
	}
	size = float64(len(population.TypeB) + len(population.TypeBDead))
	s = Stat{
		Mean:     float64(total) / size,
		Size:     size,
		Variance: 0,
	}
	typeBFitness = append(typeBFitness, s)
}

var record [][]SimpleAgent
var stepsRecorded int

// keep allocated for smoother GC
var newStep []SimpleAgent

func RecordPositions(population Population) {
	numDead := len(population.TypeADead) + len(population.TypeBDead)
	newStep = make([]SimpleAgent, POPULATION_SIZE-numDead)
	for i, agent := range population.TypeA {
		newStep[i] = *agent
	}
	for i, agent := range population.TypeB {
		newStep[len(population.TypeA)+i] = *agent
	}
	for i, agent := range population.Predators {
		newStep[len(population.TypeA)+len(population.TypeB)+i] = *agent
	}
	stepsRecorded++
	record = append(record, newStep)
}

func writeJSON(obj interface{}, filename string) {
	jsonEnc, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	_, err = f.Write(jsonEnc)
	if err != nil {
		fmt.Println(err)
	}
}

func WritePositions() {
	writeJSON(record, "output/Positions.json")
	record = nil
}

func WriteStatistics() {
	fmt.Println("Writing stats")

	writeJSON(nearbyA, "output/NearbyA.json")
	writeJSON(nearbyB, "output/NearbyB.json")

	writeJSON(geneticsA, "output/GeneticsA.json")
	writeJSON(geneticsB, "output/GeneticsB.json")

	// writeJSON(typeADead, "output/DeadA.json")
	// writeJSON(typeBDead, "output/DeadB.json")

	writeJSON(typeAFitness, "output/FitnessA.json")
	writeJSON(typeBFitness, "output/FitnessB.json")
}
