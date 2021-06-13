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

type coords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type snake struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Health  int      `json:"health"`
	Body    []coords `json:"body"`
	Latency int      `json:"latency"`
	Head    coords   `json:"head"`
	Length  int      `json:"length"`
	Shout   string   `json:"shout"`
	Squad   string   `json:"squad"`
}

type state struct {
	Game struct {
		ID      string `json:"id"`
		RuleSet struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"rule_set"`
		Timeout int `json:"timeout"` // in milliseconds
	} `json:"game"`
	Turn  int `json:"turn"`
	Board struct {
		Height  int      `json:"height"`
		Width   int      `json:"width"`
		Food    []coords `json:"food"`
		Hazards []coords `json:"hazards"`
		Snakes  []snake  `json:"snakes"`
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
	fmt.Printf("deleting game %s\n", s.Game.ID)
	delete(runningGames, s.Game.ID)
	fmt.Printf("snakes left: %+v\n", runningGames)
}

func runSnake(c chan []byte, initialState *state) {
	printBoard(*initialState)

loop:
	for {
		switch msg := <-c; msg {
		case nil:
			fmt.Println("Snake ended")
			break loop
		default:
			s, err := getGameState(msg)
			printBoard(*s)
			if err != nil {
				log.Println(err)
			}
			// spew.Dump(s)
			fmt.Printf("%s got request to make move decision\n", s.You.Name)
			c <- []byte(Up)
			fmt.Printf("%s made decision\n", s.You.Name)
		}
	}
}

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
		fmt.Printf("%s: %+v\n", snakeChar, snake.Body)
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
	for i := h - 1; i > 0; i-- {
		fmt.Printf(" ")
		for j := 0; j < w; j++ {
			fmt.Printf("%s ", board[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}
