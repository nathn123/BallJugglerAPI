package main

import (
	"errors"
	"fmt"
)

//Game stores data about the game
type Game struct {
	ID int
	// start Time.Time
	draw        []int
	ticketCount uint
	ticketMin   int
	ticketMax   int
	result      ResultStruct
	chatHub     *Chat
}

//Player stores data about a player
type Player struct {
	ID      int
	Name    string
	tickets []TicketPack
}

// GenTicketPack Generate Tickets for a game
func (g *Game) GenTicketPack(num int, pID int) TicketPack {
	tp := TicketPack{playerID: pID, gameID: g.ID}
	tp.Gen(num, g.ticketMin, g.ticketMax, g.draw)
	fmt.Println("GENERATED TICKETS :", tp)
	return tp
}

func (g *Game) genDraw() {
	g.draw = draw(g.ticketMin, g.ticketMax)
}

//Init to initialise the game of bingo
func (g *Game) Init(players []Player) {
	fmt.Printf("Game ID: %v Initialising\n", g.ID)
	fmt.Printf("Game ID: %v Num Players: %v\n", g.ID, len(players))
	// check winner
	// gen a win map of players and their win times
	var winMap = make(map[int]ticketWin)
	for _, p := range players {
		tw, err := findTicketWin(p, g.ID)
		if err != nil {
			fmt.Println("Player has not tickets in this game", err)
		} else {
			fmt.Println(tw)
			winMap[p.ID] = tw
		}
	}
	fmt.Printf("Game ID: %v WinMap\n", winMap)
	g.result = FindWinner(winMap)
	fmt.Printf("Game ID: %v Initialised\n", g.ID)
	fmt.Printf("Game ID: %v Results\n", g.result)

}

func findTicketWin(p Player, ID int) (ticketWin, error) {
	var t ticketWin
	for _, tp := range p.tickets {
		fmt.Println("TICKET PACK:", tp)
		if tp.gameID == ID {
			// add this ticket pack to the decision
			return tp.win, nil
		}
	}
	return t, errors.New("No matching ticket pack")
}
