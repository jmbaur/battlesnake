package game

import (
	"errors"
	"fmt"
	"math"
)

func canMoveDirection(s *state, direction string) bool {
	x := s.You.Head.X
	y := s.You.Head.Y

	// check for collision with board walls
	if (direction == Up && y+1 >= s.Board.Height) ||
		(direction == Down && y-1 <= 0) ||
		(direction == Right && x+1 >= s.Board.Width) ||
		(direction == Left && x-1 <= 0) {
		return false
	}

	// check for collision with any snakes
	for _, snake := range s.Board.Snakes {
		for _, segment := range snake.Body {
			if (direction == Up && y+1 == segment.Y && x == segment.X) ||
				(direction == Down && y-1 == segment.Y && x == segment.X) ||
				(direction == Right && x+1 == segment.X && y == segment.Y) ||
				(direction == Left && x-1 == segment.X && y == segment.Y) {
				return false
			}
		}
	}

	return true
}

var ErrNoFood = errors.New("no food available")

func closestFood(s *state) (*coordinate, error) {
	x := s.You.Head.X
	y := s.You.Head.Y

	if len(s.Board.Food) == 0 {
		return nil, ErrNoFood
	}

	lowest := math.Sqrt(math.Pow(float64(s.Board.Height), 2) + math.Pow(float64(s.Board.Width), 2))
	var lowestCoords coordinate

	for _, food := range s.Board.Food {
		if distance := math.Sqrt(math.Pow(float64(x-food.X), 2) + math.Pow(float64(y-food.Y), 2)); distance < lowest {
			lowest = distance
			lowestCoords = food
		}
	}
	return &lowestCoords, nil
}
