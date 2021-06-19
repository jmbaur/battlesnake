package game

import "fmt"

func canMoveDirection(s *state, direction string) bool {
	x := s.You.Head.X
	y := s.You.Head.Y

	b := true
	for _, snake := range s.Board.Snakes {
		for _, segment := range snake.Body {
			// Check whether the next move made would place us either...
			//   1. Out of bounds.
			//   2. In conflict with another snake (or ourselves).
			switch direction {
			case Up:
				fmt.Println("Up", y, s.Board.Height, snake.Name, segment.Y)
				if y+1 == s.Board.Height || (y+1 == segment.Y && x == segment.X) {
					b = false
				}
			case Down:
				fmt.Println("Down", y, 0, snake.Name, segment.Y)
				if y-1 == 0 || (y-1 == segment.Y && x == segment.X) {
					b = false
				}
			case Right:
				fmt.Println("Right", x, s.Board.Width, snake.Name, segment.X)
				if x+1 == s.Board.Width || (x+1 == segment.X && y == segment.Y) {
					b = false
				}
			case Left:
				fmt.Println("Left", x, 0, snake.Name, segment.X)
				if x-1 == 0 || (x-1 == segment.X && y == segment.Y) {
					b = false
				}
			}
		}
	}

	return b
}
