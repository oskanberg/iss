package main

import (
	"math/rand"
	"time"

	"github.com/oskanberg/go-vector"
)

const (
	PREYA = iota
	PREYB
	PRED
)

type SimpleAgent struct {
	Position     vector.Vector2D `json:"p"`
	Velocity     vector.Vector2D `json:"v"`
	VelocityNext vector.Vector2D `json:"vn"`
	Fitness      int             `json:"f"`
	AgentType    int             `json:"t"`
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	population := make([]*SimpleAgent, POPULATION_SIZE)
	for i, _ := range population {
		randomLocation := *vector.NewVector2d(rand.Float64()*SIMULATION_SPACE_SIZE, rand.Float64()*SIMULATION_SPACE_SIZE)
		randomVelocity := *vector.NewRandomUnitVector()
		population[i] = &SimpleAgent{
			Position:     randomLocation,
			Velocity:     randomVelocity,
			VelocityNext: randomVelocity,
			Fitness:      0,
			AgentType:    PREYA,
		}
	}
	for i := 0; i < SIMULATION_STEPS; i++ {
		Move(population)
		UpdatePosition(population)
		Report(population)
	}
}
