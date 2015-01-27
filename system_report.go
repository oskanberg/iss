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

func Report(agents []*SimpleAgent) {
	newStep = make([]SimpleAgent, POPULATION_SIZE)
	for i, agent := range agents {
		newStep[i] = *agent
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
	}
}
