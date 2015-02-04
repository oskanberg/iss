package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var record [][]SimpleAgent
var stepsRecorded int

// keep allocated for smoother GC
var newStep []SimpleAgent

func Report(population Population) {
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

	record = append(record, newStep)

	stepsRecorded++
	if stepsRecorded%POSITION_LOG_INTERVAL == 0 {
		fmt.Println("Writing to", strconv.Itoa(stepsRecorded), ".json")
		jsonEnc, err := json.Marshal(record)
		if err != nil {
			fmt.Println(err)
		}
		f, err := os.Create("output/" + strconv.Itoa(stepsRecorded) + ".json")
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		_, err = f.Write(jsonEnc)
		if err != nil {
			fmt.Println(err)
		}

		record = nil
	}
}
