package main

import "github.com/oskanberg/go-vector"

var agentRepulsionVector *vector.Vector2D
var agentOrientationVector *vector.Vector2D
var agentAttractionVector *vector.Vector2D
var antiPredatorVector *vector.Vector2D

func updateVectorsForPrey(agent, other *SimpleAgent, repulsion, orientation, attraction *vector.Vector2D) {
	var behaviour *MovementParameters

	if agent.Family == other.Family {
		behaviour = &agent.genetics.SameSpecies
	} else {
		behaviour = &agent.genetics.OtherSpecies
	}

	differenceVector := agent.Position.WrappedDistanceVector(&other.Position, SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
	difference := differenceVector.MagnitudeSquared()
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
	differenceVector := agent.Position.WrappedDistanceVector(&other.Position, SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
	if differenceVector.MagnitudeSquared() < VIEW_DISTANCE_SQUARED {
		*antiPredator = *antiPredator.Add(differenceVector.Normalised().Multiplied(-1))
	}
}

func movePrey(agent *SimpleAgent, population Population) {
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
		agent.VelocityNext = *agentRepulsionVector.Normalised().Multiplied(PREY_SPEED)
	} else {
		antiPred := antiPredatorVector.Normalised().Multiplied(agent.genetics.PredatorRepulsion)
		updateVector := agentOrientationVector.Add(agentAttractionVector).Add(antiPred).Normalised().Multiplied(PREY_SPEED)
		agent.VelocityNext = *updateVector
	}
}

func movePredator(agent *SimpleAgent, population *Population) {
	var differenceVector *vector.Vector2D
	var updateVector *vector.Vector2D = vector.NewVector2d(1000, 0)
	var nearest *SimpleAgent
	var shortestDistance float64 = 10000000

	for _, prey := range population.TypeA {
		differenceVector = agent.Position.WrappedDistanceVector(&prey.Position, SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
		difference := differenceVector.MagnitudeSquared()
		if difference < shortestDistance {
			shortestDistance = difference
			updateVector = differenceVector
			nearest = prey
		}
	}

	for _, prey := range population.TypeB {
		differenceVector = agent.Position.WrappedDistanceVector(&prey.Position, SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
		difference := differenceVector.MagnitudeSquared()
		if difference < shortestDistance {
			shortestDistance = difference
			updateVector = differenceVector
			nearest = prey
		}
	}

	if updateVector.Magnitude() < PREDATOR_SPEED {
		Kill(population, nearest)
	}
	agent.VelocityNext = *updateVector.Normalised().Multiplied(PREDATOR_SPEED)
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
		agent.Position = *agent.Position.Add(&agent.Velocity).Wrap(SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
	}
	for _, agent := range population.TypeB {
		agent.Velocity = agent.VelocityNext
		agent.Position = *agent.Position.Add(&agent.Velocity).Wrap(SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
	}
	for _, agent := range population.Predators {
		agent.Velocity = agent.VelocityNext
		agent.Position = *agent.Position.Add(&agent.Velocity).Wrap(SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
	}
}
