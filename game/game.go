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

type ruleset struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type game struct {
	ID      string  `json:"id"`
	RuleSet ruleset `json:"ruleset"`
	Timeout int     `json:"timeout"` // in milliseconds
}

type board struct {
	Height  int          `json:"height"`
	Width   int          `json:"width"`
	Food    []coordinate `json:"food"`
	Hazards []coordinate `json:"hazards"`
	Snakes  []snake      `json:"snakes"`
}

type state struct {
	Game  game  `json:"game"`
	Turn  int   `json:"turn"`
	Board board `json:"board"`
	You   snake `json:"you"`
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

			var directions []string
			for _, d := range []string{Up, Down, Right, Left} {
				if canMoveDirection(s, d) {
					directions = append(directions, d)
				}
			}

			var direction string

			if s.You.Health < (s.Board.Height+s.Board.Width)/2 {
				food, err := closestFood(s)
				if err != ErrNoFood {
					direction = desiredDirection(s.You.Head, *food)
				}
			}

			c <- []byte(direction)

			log.Printf("%s is going %s\n", s.You.Name, direction)
		}
	}
}
