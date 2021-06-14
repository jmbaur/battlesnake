package game

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	Up    = "up"
	Down  = "down"
	Left  = "left"
	Right = "right"
)

var runningGames = make(map[string]chan []byte)

type coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type snake struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Health  int          `json:"health"`
	Body    []coordinate `json:"body"`
	Latency int          `json:"latency"`
	Head    coordinate   `json:"head"`
	Length  int          `json:"length"`
	Shout   string       `json:"shout"`
	Squad   string       `json:"squad"`
}

type state struct {
	Game struct {
		ID      string `json:"id"`
		RuleSet struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"ruleset"`
		Timeout int `json:"timeout"` // in milliseconds
	} `json:"game"`
	Turn  int `json:"turn"`
	Board struct {
		Height  int          `json:"height"`
		Width   int          `json:"width"`
		Food    []coordinate `json:"food"`
		Hazards []coordinate `json:"hazards"`
		Snakes  []snake      `json:"snakes"`
	} `json:"board"`
	You snake `json:"you"`
}

func getGameState(data []byte) (*state, error) {
	var s state
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal start game data: %v", err)
	}
	return &s, nil
}

func Start(data []byte) {
	s, err := getGameState(data)
	if err != nil {
		log.Println(err)
	}

	snakeChan := make(chan []byte)
	runningGames[s.Game.ID] = snakeChan

	go runSnake(snakeChan, s)
}

func Decide(data []byte) string {
	s, err := getGameState(data)
	if err != nil {
		log.Println(err)
	}
	snakeChan, ok := runningGames[s.Game.ID]
	if !ok {
		log.Printf("Could not find game with ID '%s'\n", s.Game.ID)
	}

	snakeChan <- data
	return string(<-snakeChan)
}

func End(data []byte) {
	s, err := getGameState(data)
	if err != nil {
		log.Println(err)
		log.Printf("%s\n", data)
	}
	snakeChan, ok := runningGames[s.Game.ID]
	if !ok {
		log.Printf("Could not find game with ID '%s'\n", s.Game.ID)
	}
	snakeChan <- nil
	defer close(snakeChan)
	log.Printf("deleting game %s\n", s.Game.ID)
	delete(runningGames, s.Game.ID)
	log.Printf("snakes left: %+v\n", runningGames)
}

func runSnake(c chan []byte, initialState *state) {
loop:
	for {
		switch msg := <-c; msg {
		case nil:
			log.Println("Snake ended")
			break loop
		default:
			s, err := getGameState(msg)
			printBoard(*s)
			if err != nil {
				log.Println(err)
			}
			log.Printf("%s got request to make move decision\n", s.You.Name)

			var direction string
			switch {
			case canMoveDirection(s, Up):
				direction = Up
			case canMoveDirection(s, Down):
				direction = Down
			case canMoveDirection(s, Right):
				direction = Right
			case canMoveDirection(s, Left):
				direction = Left
			default:
				log.Println("could not make decision on where to move")
				// TODO: does the game engine have some default value? (e.g. move the
				// same direction as last move)
			}
			c <- []byte(direction)

			log.Printf("%s is going %s\n", s.You.Name, direction)
		}
	}
}

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
