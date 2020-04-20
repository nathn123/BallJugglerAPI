package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type gameJSON struct {
	StartDate time.Time `json:"startDate"`
	TicketMin int       `json:"ticketMin"`
	TicketMax int       `json:"ticketMax"`
}

type playerJSON struct {
	Name string `json:"name"`
}

type ticketJSON struct {
	GameID   int `json:"gameID"`
	PlayerID int `json:"playerID"`
	Num      int `json:"num"`
}

type chatJSON struct {
	GameID   int `json:"gameID"`
	PlayerID int `json:"playerID"`
}

// func main() {

// 	fmt.Print("To add a player Enter P:<Name>\n")
// 	fmt.Print("To Start the game Enter Start\n")
// 	scanner := bufio.NewScanner(os.Stdin)
// 	var players []Player

// 	for scanner.Scan() {
// 		if scanner.Text() == "Start" {
// 			var drawLength = 90
// 			draw := draw(drawLength)
// 			var game = Game{ID: 1, draw: draw, ticketMin: 1, ticketMax: 90}
// 			game.Start(players)
// 		} else if strings.Contains(scanner.Text(), "P:") {
// 			v := strings.Split(scanner.Text(), ":")
// 			players = append(players, Player{ID: 1, Name: v[1]})
// 			fmt.Printf("Player Added: %v\n", v[1])

// 		} else if strings.Contains(scanner.Text(), "TEST") {
// 			var drawLength = 90
// 			draw := draw(drawLength)
// 			for i := 0; i < 100; i++ {
// 				var players []Player
// 				var game = Game{ID: 1, draw: draw, ticketMin: 1, ticketMax: 90}
// 				for i := 0; i < 20; i++ {
// 					players = append(players, Player{ID: i})
// 				}
// 				game.Start(players)
// 			}
// 		}
// 	}
// }
type mainObj struct {
	s  *Schedule
	db *Db
}

func main() {
	var schedule Schedule
	var db Db
	o := mainObj{s: &schedule, db: &db}
	o.Init()
	router := mux.NewRouter()
	// base
	router.HandleFunc("/", home)
	// see games avaliable
	router.HandleFunc("/game", game).Methods("GET")
	// add game
	router.HandleFunc("/game", o.addGame).Methods("POST")
	// join a game
	router.HandleFunc("/game/join/{playerID}/{gameID}", o.joinGame)
	// see details about game on get
	router.HandleFunc("/game/:id", game)
	// join game via websocket
	router.HandleFunc("/game/:id", game)
	// create player on post
	// get player on get
	router.HandleFunc("/player", o.addPlayer).Methods("POST")
	// buy tickets on post
	router.HandleFunc("/tickets", o.buyTickets).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "THIS IS THE HOMEPAGE")
	// json.Encoder("THIS IS THE HOMEPAGE")
}
func game(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "THIS IS THE HOMEPAGE")
	// json.Encoder("THIS IS THE HOMEPAGE")
}

func (m *mainObj) Init() {
	// create scheduler
	m.s.Init()
	// create db
	m.db.Init(m.s)
}

func (m *mainObj) addGame(w http.ResponseWriter, r *http.Request) {
	var gj gameJSON
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &gj); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		fmt.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	gameID := len(m.db.games)
	if gj.TicketMin < 1 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		json.NewEncoder(w).Encode("Ticket minimum must be 1 or greater")
		return
	}
	if gj.TicketMax < gj.TicketMin {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		json.NewEncoder(w).Encode("Ticket max must be greater than Ticket min")
		return
	}
	t, _ := time.ParseDuration("1s")
	if time.Until(gj.StartDate) < t {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		json.NewEncoder(w).Encode("Start date must be in the future")
		return
	}
	gd := Game{ID: gameID, ticketMin: gj.TicketMin, ticketMax: gj.TicketMax, chatHub: NewHub()}
	gd.genDraw()

	m.db.AddGame(gameID, gd, gj.StartDate, m.s)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(gd); err != nil {
		panic(err)
	}

}
func (m *mainObj) addPlayer(w http.ResponseWriter, r *http.Request) {
	var pj playerJSON
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &pj); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		fmt.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	playerID := len(m.db.players)
	pd := Player{ID: playerID, Name: pj.Name}

	m.db.AddPlayer(playerID, pd)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(pd); err != nil {
		panic(err)
	}

}
func (m *mainObj) buyTickets(w http.ResponseWriter, r *http.Request) {
	var tj ticketJSON
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &tj); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		fmt.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	ticketID := len(m.db.tickets)

	m.db.AddTicketPack(ticketID, tj.GameID, tj.PlayerID, tj.Num)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ticketID); err != nil {
		panic(err)
	}

}
func (m *mainObj) joinGame(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	gameID, err := strconv.Atoi(params["gameID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	playerID, err := strconv.Atoi(params["playerID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	g, err := m.db.GetGame(gameID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	_, err = m.db.GetPlayerInGame(playerID, gameID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	ServeWs(g.chatHub, w, r)
}
