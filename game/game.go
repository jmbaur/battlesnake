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

type cell struct {
	visited   bool
	visitable bool
	kind      int
}

// Kinds of cells.
const (
	emptyCell = iota
	headCell
	bodyCell
	foodCell
)

type graph struct {
	cells  [][]cell
	height int
	width  int
}

func stateToGraph(s *state) *graph {
	var g graph

	g.height = s.Board.Height
	g.width = s.Board.Width

	for h := 0; h <= s.Board.Height; h++ {
		g.cells = append(g.cells, []cell{})
		for w := 0; w <= s.Board.Width; w++ {
			g.cells[h] = append(g.cells[h], cell{visitable: true})
		}
	}

	for _, snake := range s.Board.Snakes {
		for _, segment := range snake.Body {
			g.cells[segment.Y][segment.X].visitable = false
			g.cells[segment.Y][segment.X].kind = bodyCell
		}
		g.cells[snake.Head.Y][snake.Head.X].visitable = false
		g.cells[snake.Head.Y][snake.Head.X].kind = headCell
	}

	for _, food := range s.Board.Food {
		g.cells[food.Y][food.X].visitable = true
		g.cells[food.Y][food.X].kind = foodCell
	}

	return &g
}

func (g *graph) String() string {
	var pp string
	pp += "\n"

	for i := g.height - 1; i >= 0; i-- {
		pp += " "
		for j := 0; j < g.width; j++ {
			var char string
			switch g.cells[i][j].kind {
			case headCell:
				char = "S"
			case bodyCell:
				char = "s"
			case foodCell:
				char = "f"
			default:
				if g.cells[i][j].visitable {
					char = "+"
				} else {
					char = "-"
				}
			}
			pp += char + " "
		}
		pp += "\n"
	}
	pp += "\n"

	return pp
}
