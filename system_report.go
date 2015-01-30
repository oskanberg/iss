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

func Report(Population Population) {
	newStep = make([]SimpleAgent, POPULATION_SIZE)
	for i, agent := range Population.TypeA {
		newStep[i] = *agent
	}
	for i, agent := range Population.TypeB {
		newStep[len(Population.TypeA)+i] = *agent
	}
	for i, agent := range Population.Predators {
		newStep[len(Population.TypeA)+len(Population.TypeB)+i] = *agent
	}
	record = append(record, newStep)

	stepsRecorded++
	if stepsRecorded%LOG_WRITE_FREQ == 0 {
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
