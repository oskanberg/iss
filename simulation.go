package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/oskanberg/go-vector"
)

const (
	PREYA = iota
	PREYB
	PRED
)

type MovementParameters struct {
	OrientationRadiusSq float64
	AttractionRadiusSq  float64
}

var upperLimit float64 = VIEW_DISTANCE_SQUARED - STATIC_REPULSION_RADIUS_SQUARED

func (s *MovementParameters) Mutated() *MovementParameters {
	dOrientation := rand.Float64()*MUTATION_MOVEMENT - (MUTATION_MOVEMENT / 2)
	dAttraction := rand.Float64()*MUTATION_MOVEMENT - (MUTATION_MOVEMENT / 2)

	newOrientation := math.Mod(s.OrientationRadiusSq+dOrientation, upperLimit) + STATIC_REPULSION_RADIUS_SQUARED
	newAttraction := math.Mod(s.AttractionRadiusSq+dAttraction, upperLimit) + STATIC_REPULSION_RADIUS_SQUARED
	return &MovementParameters{
		OrientationRadiusSq: newOrientation,
		AttractionRadiusSq:  newAttraction,
	}
}

type BehaviourParameters struct {
	SameSpecies       MovementParameters
	OtherSpecies      MovementParameters
	PredatorRepulsion float64
}

func (s *BehaviourParameters) Mutated() *BehaviourParameters {
	dRepulsion := (rand.Float64() * MUTATION_PREDATOR_REPULSION) - (MUTATION_PREDATOR_REPULSION / 2)
	newRepulsion := math.Mod(s.PredatorRepulsion+dRepulsion, 1.0)
	return &BehaviourParameters{
		SameSpecies:       *s.SameSpecies.Mutated(),
		OtherSpecies:      *s.OtherSpecies.Mutated(),
		PredatorRepulsion: newRepulsion,
	}
}

type SimpleAgent struct {
	Position     vector.Vector2D `json:"p"`
	Velocity     vector.Vector2D `json:"v"`
	VelocityNext vector.Vector2D `json:"vn"`
	Fitness      int             `json:"f"`
	Family       int             `json:"t"`

	alive    bool
	genetics *BehaviourParameters
}

type Population struct {
	TypeA     []*SimpleAgent
	TypeB     []*SimpleAgent
	Predators []*SimpleAgent
	TypeADead []*SimpleAgent
	TypeBDead []*SimpleAgent
}

func RandomMovementParameters() MovementParameters {
	a := rand.Float64()*(VIEW_DISTANCE-STATIC_REPULSION_RADIUS) + STATIC_REPULSION_RADIUS
	o := rand.Float64()*(a-STATIC_REPULSION_RADIUS) + STATIC_REPULSION_RADIUS
	return MovementParameters{
		OrientationRadiusSq: o * o,
		AttractionRadiusSq:  a * a,
	}
}

func RandomBehaviours() *BehaviourParameters {
	return &BehaviourParameters{
		SameSpecies:       RandomMovementParameters(),
		OtherSpecies:      RandomMovementParameters(),
		PredatorRepulsion: rand.Float64(),
	}
}

func NewRandomSimpleAgent(family int) *SimpleAgent {
	randomLocation := *vector.NewVector2d(rand.Float64()*SIMULATION_SPACE_SIZE, rand.Float64()*SIMULATION_SPACE_SIZE)
	randomVelocity := *vector.NewRandomUnitVector()
	return &SimpleAgent{
		Position:     randomLocation,
		Velocity:     randomVelocity,
		VelocityNext: randomVelocity,
		Fitness:      0,
		Family:       family,
		alive:        true,
		genetics:     RandomBehaviours(),
	}
}

func main() {
	// defer profile.Start(profile.CPUProfile).Stop()
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())
	population := Population{
		TypeA:     make([]*SimpleAgent, SUBPOPULATION_SIZE),
		TypeB:     make([]*SimpleAgent, SUBPOPULATION_SIZE),
		Predators: make([]*SimpleAgent, PREDATOR_POPULATION_SIZE),
	}
	for i, _ := range population.TypeA {
		population.TypeA[i] = NewRandomSimpleAgent(PREYA)
	}
	for i, _ := range population.TypeB {
		population.TypeB[i] = NewRandomSimpleAgent(PREYB)
	}
	for i, _ := range population.Predators {
		population.Predators[i] = NewRandomSimpleAgent(PRED)
	}

	for i := 0; i < SIMULATION_STEPS; i++ {
		if i > 0 && i%1000 == 0 {
			fmt.Println("step", i)
			Evolve(&population)
		}
		Move(&population)
		UpdatePosition(population)
		UpdateFitness(population)
		Report(population)
	}
}
