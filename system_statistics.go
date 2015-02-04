package main

import "fmt"

var typeAMeanPredatorRepulsion float64
var typeBMeanPredatorRepulsion float64

var typeAFitness float64
var typeBFitness float64

var typeAAgentsRecorded int
var typeBAgentsRecorded int
var iterationsRecorded int

func Statistics(population Population) {
	iterationsRecorded++

	for _, agent := range population.TypeA {
		typeAMeanPredatorRepulsion += agent.Genetics.PredatorRepulsion
		typeAFitness++
	}
	typeAAgentsRecorded += len(population.TypeA)

	for _, agent := range population.TypeB {
		typeBMeanPredatorRepulsion += agent.Genetics.PredatorRepulsion
		typeBFitness++
	}
	typeBAgentsRecorded += len(population.TypeB)

	if iterationsRecorded%STATISTICS_LOG_INTERVAL == 0 {
		typeAMeanPredatorRepulsion /= float64(typeAAgentsRecorded)
		typeBMeanPredatorRepulsion /= float64(typeBAgentsRecorded)

		fmt.Println("Type A repulsion:", typeAMeanPredatorRepulsion)
		fmt.Println("Type B repulsion:", typeBMeanPredatorRepulsion)

		fmt.Println("Type A fitness:", typeAFitness/SUBPOPULATION_SIZE)
		fmt.Println("Type B fitness:", typeBFitness/SUBPOPULATION_SIZE)

		typeAMeanPredatorRepulsion = 0
		typeBMeanPredatorRepulsion = 0
		typeAAgentsRecorded = 0
		typeBAgentsRecorded = 0
		iterationsRecorded = 0
		typeAFitness = 0
		typeBFitness = 0
	}
}
