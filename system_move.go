package main

import "github.com/oskanberg/go-vector"

var agentRepulsionVector vector.Vector2D
var agentOrientationVector vector.Vector2D
var agentAttractionVector vector.Vector2D

func Move(agents []*SimpleAgent) {
	for _, agent := range agents {

		agentRepulsionVector := vector.NewVector2d(0, 0)
		agentOrientationVector := vector.NewVector2d(0, 0)
		agentAttractionVector := vector.NewVector2d(0, 0)

		for _, other := range agents {
			if agent == other {
				continue
			}

			differenceVector := agent.Position.Subtract(&other.Position)
			difference := differenceVector.Magnitude()
			if difference < STATIC_REPULSION_RADIUS {
				// repulsion
				agentRepulsionVector = agentRepulsionVector.Add(differenceVector.Normalised())
			} else if difference > STATIC_REPULSION_RADIUS && difference < STATIC_ORIENTATION_RADIUS {
				// orientation
				agentOrientationVector = agentOrientationVector.Add(&other.Velocity)
			} else if difference > STATIC_ORIENTATION_RADIUS && difference < STATIC_ATTRACTION_RADIUS {
				// attraction
				agentAttractionVector = agentAttractionVector.Subtract(differenceVector.Normalised())
			}
		}

		if agentRepulsionVector.Magnitude() > 0 {
			agent.VelocityNext = *agentRepulsionVector.Normalised()
		} else {
			agent.VelocityNext = *agent.Velocity.Add(agentOrientationVector).Add(agentAttractionVector).Normalised()
		}
	}
}

func UpdatePosition(agents []*SimpleAgent) {
	for _, agent := range agents {
		agent.Velocity = agent.VelocityNext
		agent.Position = *agent.Position.Add(&agent.Velocity).Wrap(SIMULATION_SPACE_SIZE, SIMULATION_SPACE_SIZE)
	}
}
