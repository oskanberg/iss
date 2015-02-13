package main

import (
	"math"
	"math/rand"

	"github.com/oskanberg/go-vector"
)

var agentRepulsionVector *vector.Vector2D
var agentOrientationVector *vector.Vector2D
var agentAttractionVector *vector.Vector2D
var antiPredatorVector *vector.Vector2D

func updateVectorsForPrey(agent, other *SimpleAgent, repulsion, orientation, attraction *vector.Vector2D) {
	var behaviour *MovementParameters

	differenceVector := agent.Position.WrappedDistanceVector(&other.Position, SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
	difference := differenceVector.MagnitudeSquared()

	angle := agent.Velocity.Dot(differenceVector.Normalised())
	if difference < PREY_VIEW_DISTANCE_SQUARED && angle < BLIND_ANGLE {
		if agent.Family == other.Family {
			behaviour = &agent.Genetics.SameSpecies
			agent.visibleSame++
		} else {
			behaviour = &agent.Genetics.OtherSpecies
			agent.visibleOther++
		}
	} else {
		return
	}

	if difference < STATIC_REPULSION_RADIUS_SQUARED {
		// repulsion
		*repulsion = *repulsion.Add(differenceVector.Normalised())
	} else if difference > STATIC_REPULSION_RADIUS_SQUARED && difference < behaviour.OrientationRadiusSq {
		// orientation
		*orientation = *orientation.Add(&other.Velocity)
	} else if difference > behaviour.OrientationRadiusSq && difference < behaviour.AttractionRadiusSq {
		// attraction
		*attraction = *attraction.Subtract(differenceVector.Normalised())
	}
}

func updateVectorsForPredator(agent, other *SimpleAgent, antiPredator *vector.Vector2D) {
	differenceVector := other.Position.WrappedDistanceVector(&agent.Position, SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
	angle := agent.Velocity.Dot(differenceVector.Normalised())
	if differenceVector.MagnitudeSquared() < PREY_VIEW_DISTANCE_SQUARED && angle < BLIND_ANGLE {
		*antiPredator = *antiPredator.Add(differenceVector.Normalised())
	}
}

func movePrey(agent *SimpleAgent, population Population) {
	agent.visibleOther = 0
	agent.visibleSame = 0

	agentRepulsionVector = vector.NewVector2d(0, 0)
	agentOrientationVector = vector.NewVector2d(0, 0)
	agentAttractionVector = vector.NewVector2d(0, 0)
	antiPredatorVector = vector.NewVector2d(0, 0)

	for _, other := range population.TypeA {
		if agent == other {
			continue
		}
		updateVectorsForPrey(agent, other, agentRepulsionVector, agentOrientationVector, agentAttractionVector)
	}

	for _, other := range population.TypeB {
		if agent == other {
			continue
		}
		updateVectorsForPrey(agent, other, agentRepulsionVector, agentOrientationVector, agentAttractionVector)
	}

	for _, other := range population.Predators {
		updateVectorsForPredator(agent, other, antiPredatorVector)
	}

	if agentRepulsionVector.MagnitudeSquared() > 0 {
		agent.VelocityNext = *agentRepulsionVector.Normalised()
	} else {
		antiPred := antiPredatorVector.Normalised().Multiplied(agent.Genetics.PredatorRepulsion)
		updateVector := agentOrientationVector.Normalised().Add(agentAttractionVector.Normalised()).Add(antiPred).Normalised()

		angle := math.Atan2(agent.Velocity.X*updateVector.Y-agent.Velocity.Y*updateVector.X, agent.Velocity.X*updateVector.X+agent.Velocity.Y*updateVector.Y)
		if math.Abs(angle) > MAX_TURN_ANGLE_PREY {
			if angle > 0 {
				updateVector = updateVector.Rotated(MAX_TURN_ANGLE_PREY - angle)
			} else {
				updateVector = updateVector.Rotated(-MAX_TURN_ANGLE_PREY - angle)
			}
		}

		if updateVector.Magnitude() > 0 {
			agent.VelocityNext = *updateVector.Normalised()
		} else {
			agent.VelocityNext = agent.Velocity
		}
	}
}

func movePredator(agent *SimpleAgent, population *Population) {
	var differenceVector *vector.Vector2D
	var updateVector *vector.Vector2D = vector.NewVector2d(1000, 0)
	var nearest *SimpleAgent
	var shortestDistance float64 = 10000000

	var inView float64 = 0

	for _, prey := range population.TypeA {
		differenceVector = agent.Position.WrappedDistanceVector(&prey.Position, SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
		difference := differenceVector.MagnitudeSquared()
		if difference < PREDATOR_VIEW_DISTANCE_SQUARED {
			inView++
		}
		if difference < shortestDistance {
			shortestDistance = difference
			updateVector = differenceVector
			nearest = prey
		}
	}

	for _, prey := range population.TypeB {
		differenceVector = agent.Position.WrappedDistanceVector(&prey.Position, SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
		difference := differenceVector.MagnitudeSquared()
		if difference < PREDATOR_VIEW_DISTANCE_SQUARED {
			inView++
		}
		if difference < shortestDistance {
			shortestDistance = difference
			updateVector = differenceVector
			nearest = prey
		}
	}

	if updateVector.Magnitude() < PREDATOR_SPEED {
		Kill(population, nearest)
		agent.Position = *vector.NewVector2d(rand.Float64()*SIMULATION_SPACE_SIZE, rand.Float64()*SIMULATION_SPACE_SIZE)
	}

	angle := math.Atan2(agent.Velocity.X*updateVector.Y-agent.Velocity.Y*updateVector.X, agent.Velocity.X*updateVector.X+agent.Velocity.Y*updateVector.Y)
	if math.Abs(angle) > MAX_TURN_ANGLE_PREDATOR {
		if angle > 0 {
			updateVector = updateVector.Rotated(MAX_TURN_ANGLE_PREDATOR - angle)
		} else {
			updateVector = updateVector.Rotated(-MAX_TURN_ANGLE_PREDATOR - angle)
		}
	}

	agent.VelocityNext = *updateVector.Normalised()
}

func Move(population *Population) {
	for _, agent := range population.TypeA {
		movePrey(agent, *population)
	}
	for _, agent := range population.TypeB {
		movePrey(agent, *population)
	}
	for _, agent := range population.Predators {
		movePredator(agent, population)
	}
}

func UpdatePosition(population Population) {
	for _, agent := range population.TypeA {
		agent.Velocity = agent.VelocityNext
		agent.Position = *agent.Position.Add(agent.Velocity.Multiplied(PREY_SPEED)).Wrap(SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
	}
	for _, agent := range population.TypeB {
		agent.Velocity = agent.VelocityNext
		agent.Position = *agent.Position.Add(agent.Velocity.Multiplied(PREY_SPEED)).Wrap(SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
	}
	for _, agent := range population.Predators {
		agent.Velocity = agent.VelocityNext
		agent.Position = *agent.Position.Add(agent.Velocity.Multiplied(PREDATOR_SPEED)).Wrap(SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
	}
}
