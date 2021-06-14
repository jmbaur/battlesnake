package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/jmbaur/battlesnake/game"
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
	fmt.Print(s)
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

func startGameHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Send()
	}

	w.WriteHeader(http.StatusOK)
	game.Start(data)
}

// MoveResponse is sent on each move request.
type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout"`
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Send()
		// TODO: what should we do here?
	}

	data, err = json.Marshal(MoveResponse{
		Move:  game.Decide(data),
		Shout: "hello, world",
	})
	if err != nil {
		log.Error().Err(err).Send()
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func endGameHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Send()
	}

	w.WriteHeader(http.StatusOK)
	game.End(data)
}
