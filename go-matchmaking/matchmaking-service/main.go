package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"
	"github.com/gorilla/websocket"
)

// Upgrader to upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketClient struct {
	conn *websocket.Conn
	id string
}

// List of active WebSocket clients
var clients = make(map[*WebSocketClient]bool)
var clientsMu sync.Mutex

func RegisterWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to Websocket", http.StatusInternalServerError)
		return
	}

	// Read player name from query params
	playerId := r.URL.Query().Get("id")
	client := &WebSocketClient{conn: conn, id: playerId}

	// Add client to the list of active clients
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()

	// Handle WebSocket connection close
	defer func() {
		clientsMu.Lock()
		delete(clients, client)
		clientsMu.Unlock()
		client.conn.Close()
	}()

	// Keep the connection open
	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func NotifyMatchReady(match *Match) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		for _, player := range append(match.Team1, match.Team2...) {
			if client.id == player.Id {
				// Send match info to this player
				client.conn.WriteJSON(match)
				fmt.Printf("Notified player %s about the match.\n", client.id)
				break
			}
		}
	}
}

// Player represents a player info sent by client
type Player struct {
	Id string    `json:"id"`
	Rating int    `json:"rating"`
	Role   string `json:"role"`
}

// Match contains two balanced teams
type Match struct {
	Team1 []Player `json:"team1"`
	Team2 []Player `json:"team2"`
}

// Matchmaker manages the pool of players and creates matches
type Matchmaker struct {
	PlayerPool []Player
	mu         sync.Mutex
}

// AddPlayer adds a player to the pool and triggers matchmaking if pool size reaches 10
func (m *Matchmaker) AddPlayer(player Player) (*Match, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Add player to the pool
	m.PlayerPool = append(m.PlayerPool, player)

	// If we have 10 players, create a match
	if len(m.PlayerPool) >= 10 {
		return m.createMatch(), true
	}
	return nil, false
}

// createMatch balances teams and creates a new match
func (m *Matchmaker) createMatch() *Match {
	// Sort players by rating (highest to lowest)
	sort.Slice(m.PlayerPool, func(i, j int) bool {
		return m.PlayerPool[i].Rating > m.PlayerPool[j].Rating
	})

	// Split players into two teams, balancing by rating
	var team1, team2 []Player
	team1Rating, team2Rating := 0, 0

	for i, player := range m.PlayerPool[:10] {
		// Add players alternatively to teams to balance ratings
		if team1Rating <= team2Rating {
			team1 = append(team1, player)
			team1Rating += player.Rating
		} else {
			team2 = append(team2, player)
			team2Rating += player.Rating
		}
		fmt.Printf("Assigning player %s to team %d\n", player.Id, (i%2)+1)
	}

	// Remove the players used to create this match from the pool
	m.PlayerPool = m.PlayerPool[10:]

	return &Match{
		Team1: team1,
		Team2: team2,
	}
}

var matchmaker = &Matchmaker{PlayerPool: []Player{}}

// handleMatchmaking handles player info submissions and matchmaking
func handleMatchmaking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var player Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	match, matched := matchmaker.AddPlayer(player)
	if matched {
		// Return the match as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(match)
		fmt.Println("Match created!")
	} else {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "Player %s added to the pool. Waiting for more players...", player.Id)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/player", handleMatchmaking)

	fmt.Println("Matchmaking service started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
