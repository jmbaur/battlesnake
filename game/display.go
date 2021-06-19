package game

import "fmt"

// The Battlesnake board is oriented as quadrant I of the Cartesian 2D plane.
func printBoard(s state) {
	board := [][]string{}

	// Set empty board.
	h := s.Board.Height
	w := s.Board.Width
	for i := 0; i < h; i++ {
		board = append(board, []string{})
		for j := 0; j < w; j++ {
			board[i] = append(board[i], "-")
		}
	}

	// Place snakes on board
	for i, snake := range s.Board.Snakes {
		snakeChar := string(rune(i + 65))
		for _, coord := range snake.Body {
			board[coord.Y][coord.X] = snakeChar
		}
	}

	// Place food on board
	foodChar := "f"
	for _, coord := range s.Board.Food {
		board[coord.Y][coord.X] = foodChar
	}

	// Place hazards on board
	hazardChar := "h"
	for _, coord := range s.Board.Hazards {
		board[coord.Y][coord.X] = hazardChar
	}

	// Print board contents.
	fmt.Printf("\n")
	fmt.Printf("Turn: %d\n", s.Turn)
	for i := h - 1; i > 0; i-- {
		fmt.Printf(" ")
		for j := 0; j < w; j++ {
			fmt.Printf("%s ", board[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}
