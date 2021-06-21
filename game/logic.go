package game

import (
	"errors"
	"fmt"
	"math"
)

func canMoveDirection(s *state, direction string) bool {
	x := s.You.Head.X
	y := s.You.Head.Y
	health := s.You.Health

	// check for collision with board walls
	if (direction == Up && y+1 >= s.Board.Height) ||
		(direction == Down && y-1 <= 0) ||
		(direction == Right && x+1 >= s.Board.Width) ||
		(direction == Left && x-1 <= 0) {
		return false
	}

	// check for collision with any snakes
	for _, snake := range s.Board.Snakes {
		canHeadButt := health > snake.Health
		for _, segment := range snake.Body {
			switch direction {
			case Up:
				if canHeadButt && x == snake.Head.X && y+1 == snake.Head.Y {
					return true
				}
				if y+1 == segment.Y && x == segment.X {
					return false
				}
			case Down:
				if canHeadButt && x == snake.Head.X && y-1 == snake.Head.Y {
					return true
				}
				if y-1 == segment.Y && x == segment.X {
					return false
				}
			case Right:
				if canHeadButt && x+1 == snake.Head.X && y == snake.Head.Y {
					return true
				}
				if x+1 == segment.X && y == segment.Y {
					return false
				}
			case Left:
				if canHeadButt && x-1 == snake.Head.X && y == snake.Head.Y {
					return true
				}
				if x-1 == segment.X && y == segment.Y {
					return false
				}
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

func desiredDirection(a, b coordinate) string {
	return "TODO"
}

func (g *graph) traversableNeighbors(coord coordinate) []coordinate {
	var coords []coordinate
	if g.height > coord.Y+1 && g.cells[coord.Y+1][coord.X].visitable {
		coords = append(coords, coordinate{coord.X, coord.Y + 1})
	}
	if g.width > coord.X+1 && g.cells[coord.Y][coord.X+1].visitable {
		coords = append(coords, coordinate{coord.X + 1, coord.Y})
	}
	if coord.Y-1 >= 0 && g.cells[coord.Y-1][coord.X].visitable {
		coords = append(coords, coordinate{coord.X, coord.Y - 1})
	}
	if coord.X-1 >= 0 && g.cells[coord.Y][coord.X-1].visitable {
		coords = append(coords, coordinate{coord.X - 1, coord.Y})
	}
	fmt.Println(coords)
	return coords
}

func dfs(g *graph, start, end coordinate) bool {
	x := start.X
	y := start.Y
	fmt.Println(x, y)

	if x == end.X && y == end.Y {
		return true
	}

	g.cells[y][x].visited = true

	for _, coord := range g.traversableNeighbors(start) {
		if !g.cells[coord.Y][coord.X].visited && dfs(g, coord, end) {
			return true
		}
	}

	return false
}
