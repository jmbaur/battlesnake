package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	lineDelim = "--------------------------------------------------------------------------------"
)

// Server config struct
type Server struct {
	// Host is the IP to use. If unspecified, the default 0.0.0.0 will be used.
	Host string
	// Port is the port this server will try to bind to.
	Port int
}

func (s Server) String() string {
	return fmt.Sprintf("%s\nServer config:\nHost: %s\nPort: %d\n%s\n", lineDelim, s.Host, s.Port, lineDelim)
}

// Run the server
func (s Server) Run() error {
	http.HandleFunc("/", getBattlesnakeHandler)
	http.HandleFunc("/start", startGameHandler)
	http.HandleFunc("/move", moveHandler)
	http.HandleFunc("/end", endGameHandler)
	fmt.Println(s)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.Host, s.Port), nil)
}

type battlesnakeInfo struct {
	APIVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
	Version    string `json:"version"`
}

func getBattlesnakeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(battlesnakeInfo{
		APIVersion: "1",
		Author:     "jmbaur",
		Color:      "#0000ff",
		Head:       "default",
		Tail:       "default",
		Version:    "0.0.1-beta",
	})
	if err != nil {
		log.Printf("Failed to marshal BattlesnakeInfo JSON: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Snake is a representation of a battlesnake mid-game.
type Snake struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Health int    `json:"health"`
	Body   []struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"body"`
	Latency string `json:"latency"`
	Head    struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"head"`
	Length int    `json:"length"`
	Shout  string `json:"shout"`
	Squad  string `json:"squad"`
}

// GameState is received on the start of a new game and for each move request.
// TODO: potentially split out this struct?
type GameState struct {
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
		Height int `json:"height"`
		Width  int `json:"width"`
		Food   []struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"food"`
		Hazards []struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"hazards"`
		Snakes []Snake `json:"snakes"`
	} `json:"board"`
	You Snake `json:"you"`
}

func startGameHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getGameState(r.Body)
	if err != nil {
		log.Println(err)
	}
	r.Body.Close()

	w.WriteHeader(http.StatusOK)
	// TODO: initialize game state
}

// MoveResponse is sent on each move request.
type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout"`
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getGameState(r.Body)
	if err != nil {
		log.Println(err)
	}
	r.Body.Close()

	// TODO: decide move

	data, err := json.Marshal(MoveResponse{
		Move:  "up",
		Shout: "hello, world",
	})
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func endGameHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getGameState(r.Body)
	if err != nil {
		log.Println(err)
	}
	r.Body.Close()

	// TODO: deallocate resources
	w.WriteHeader(http.StatusOK)
}

func getGameState(reader io.Reader) (*GameState, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from request body: %v", err)
	}
	var state GameState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal start game data: %v", err)
	}
	return &state, nil
}
