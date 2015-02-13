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
	orientationRadius   float64
	attractionRadius    float64
	OrientationRadiusSq float64 `json:"o"`
	AttractionRadiusSq  float64 `json:"a"`
}

func (s *MovementParameters) Mutated() *MovementParameters {
	delta := randFloat(-MUTATION_MOVEMENT, MUTATION_MOVEMENT)
	newAttraction := math.Min(s.attractionRadius+delta, PREY_VIEW_DISTANCE)
	newAttraction = math.Max(newAttraction, STATIC_REPULSION_RADIUS)

	delta = randFloat(-MUTATION_MOVEMENT, MUTATION_MOVEMENT)
	newOrientation := math.Min(s.orientationRadius+delta, newAttraction)
	newOrientation = math.Max(newOrientation, STATIC_REPULSION_RADIUS)
	return &MovementParameters{
		orientationRadius:   newOrientation,
		attractionRadius:    newAttraction,
		OrientationRadiusSq: newOrientation * newOrientation,
		AttractionRadiusSq:  newAttraction * newAttraction,
	}
}

type BehaviourParameters struct {
	SameSpecies       MovementParameters `json:"ss"`
	OtherSpecies      MovementParameters `json:"os"`
	PredatorRepulsion float64            `json:"pr"`
}

func (s *BehaviourParameters) Mutated() *BehaviourParameters {
	dRepulsion := randFloat(-MUTATION_PREDATOR_REPULSION, MUTATION_PREDATOR_REPULSION)
	newRepulsion := math.Mod(s.PredatorRepulsion+dRepulsion, 1.0)
	return &BehaviourParameters{
		SameSpecies:       *s.SameSpecies.Mutated(),
		OtherSpecies:      *s.OtherSpecies.Mutated(),
		PredatorRepulsion: newRepulsion,
	}
}

type SimpleAgent struct {
	Position     vector.Vector2D      `json:"p"`
	Velocity     vector.Vector2D      `json:"v"`
	VelocityNext vector.Vector2D      `json:"vn"`
	Fitness      int                  `json:"f"`
	Family       int                  `json:"t"`
	Genetics     *BehaviourParameters `json:"g"`

	visibleSame  int
	visibleOther int
}

type Population struct {
	TypeA     []*SimpleAgent
	TypeB     []*SimpleAgent
	Predators []*SimpleAgent
	TypeADead []*SimpleAgent
	TypeBDead []*SimpleAgent
}

func RandomMovementParameters() MovementParameters {
	a := rand.Float64() * PREY_VIEW_DISTANCE
	o := rand.Float64() * a
	return MovementParameters{
		orientationRadius:   o,
		attractionRadius:    a,
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
		Genetics:     RandomBehaviours(),
		visibleSame:  0,
		visibleOther: 0,
	}
}

func main() {
	// defer profile.Start(profile.CPUProfile).Stop()
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())

	population := Population{
		TypeA: make([]*SimpleAgent, SUBPOPULATION_SIZE),
		TypeB: make([]*SimpleAgent, SUBPOPULATION_SIZE),
	}
	for i, _ := range population.TypeA {
		population.TypeA[i] = NewRandomSimpleAgent(PREYA)
	}
	for i, _ := range population.TypeB {
		population.TypeB[i] = NewRandomSimpleAgent(PREYB)
	}

	for generation := 0; generation < GENERATIONS; generation++ {
		fmt.Println("Generation", generation)
		RecordGenetics(population)

		population.Predators = nil
		for i := 0; i < WARM_UP_PERIOD; i++ {
			RecordPositions(population)
			// don't record mmnp durinng warmup?
			// RecordNearby(population)

			Move(&population)
			UpdatePosition(population)
		}

		fmt.Println("Predation")
		population.Predators = make([]*SimpleAgent, PREDATOR_POPULATION_SIZE)
		for i, _ := range population.Predators {
			population.Predators[i] = NewRandomSimpleAgent(PRED)
		}

		for i := 0; i < (EVOLUTION_INTERVAL - WARM_UP_PERIOD); i++ {
			RecordPositions(population)
			RecordNearby(population)

			Move(&population)
			UpdatePosition(population)
			UpdateFitness(population)
		}

		if generation > 0 && generation%POSITION_LOG_INTERVAL == 0 {
			WritePositions()
		}

		if generation > 0 && generation%STATISTICS_LOG_INTERVAL == 0 {
			WriteStatistics()
		}

		fmt.Println("End of generation")
		Evolve(&population)
	}
}

func randFloat(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
