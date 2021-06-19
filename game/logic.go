package game

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
