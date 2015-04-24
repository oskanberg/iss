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

var thresholdA []FloatRecord
var thresholdB []FloatRecord

var geneticsA []GeneticsRecord
var geneticsB []GeneticsRecord

var nearbyA []FloatRecord
var nearbyB []FloatRecord

var iterationsRecorded int

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

func RecordGeneticsAverage(population Population) {
	geneRecord := getAverageGenetics(population.TypeA, population.TypeADead)
	geneticsA = append(geneticsA, geneRecord)
	geneRecord = getAverageGenetics(population.TypeB, population.TypeBDead)
	geneticsB = append(geneticsB, geneRecord)
}

func RecordThresholdGenetics(population Population) {
	newCount := FloatRecord{Same: 0, Other: 0}
	for _, agent := range population.TypeA {
		if agent.Genetics.SameSpecies.attractionRadius > GENETIC_THRESHOLD {
			newCount.Other++
		}
		if agent.Genetics.OtherSpecies.attractionRadius > GENETIC_THRESHOLD {
			newCount.Other++
		}
	}
	thresholdA = append(thresholdA, newCount)

	newCount = FloatRecord{Same: 0, Other: 0}
	for _, agent := range population.TypeA {
		if agent.Genetics.SameSpecies.attractionRadius > GENETIC_THRESHOLD {
			newCount.Other++
		}
		if agent.Genetics.OtherSpecies.attractionRadius > GENETIC_THRESHOLD {
			newCount.Other++
		}
	}
	thresholdB = append(thresholdB, newCount)
}

var detailGeneticsRecordA [][]BehaviourParameters
var detailGeneticsRecordB [][]BehaviourParameters

func RecordIndividualGenetics(population Population) {
	newStep := make([]BehaviourParameters, len(population.TypeA))
	lenA := len(population.TypeA)
	for i, agent := range population.TypeA {
		newStep[i] = *agent.Genetics
	}
	detailGeneticsRecordA = append(detailGeneticsRecordA, newStep)

	newStep = make([]BehaviourParameters, len(population.TypeB))
	for i, agent := range population.TypeB {
		newStep[lenA+i] = *agent.Genetics
	}
	detailGeneticsRecordB = append(detailGeneticsRecordB, newStep)
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
	// numDead := len(population.TypeADead) + len(population.TypeBDead)
	newStep = make([]SimpleAgent, len(population.TypeA)+len(population.TypeB)+len(population.Predators))
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

	writeJSON(thresholdA, "output/ThresholdGeneticsA.json")
	writeJSON(thresholdB, "output/ThresholdGeneticsB.json")

	// writeJSON(detailGeneticsRecordA, "output/DetailGeneticsA.json")
	// writeJSON(detailGeneticsRecordB, "output/DetailGeneticsB.json")

	writeJSON(geneticsA, "output/GeneticsA.json")
	writeJSON(geneticsB, "output/GeneticsB.json")

	// writeJSON(typeADead, "output/DeadA.json")
	// writeJSON(typeBDead, "output/DeadB.json")

	// writeJSON(typeAFitness, "output/FitnessA.json")
	// writeJSON(typeBFitness, "output/FitnessB.json")
}
