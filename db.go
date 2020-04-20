package main

import (
	"errors"
	"fmt"
	"time"
)

type Db struct {
	games   []dbGameEntry
	players []dbPlayerEntry
	tickets []dbTicketPackEntry
}

type dbGameEntry struct {
	ID        int
	gameData  Game
	startTime time.Time
}

type dbPlayerEntry struct {
	ID         int
	playerData Player
}

type dbTicketPackEntry struct {
	ID         int
	ticketData TicketPack
}

func (d *Db) Init(s *Schedule) {
	fmt.Println("DB init")
	// createGameData(d, s)
}

func createGameData(db *Db, s *Schedule) {
	var numP, numG = 5, 5
	for i := 0; i < numP; i++ {
		p := Player{ID: i, Name: fmt.Sprintf("Player %v", i)}
		db.AddPlayer(i, p)
	}
	startTime := time.Now().Local()
	for i := 0; i < numG; i++ {
		// every three mins
		startTime = startTime.Add(10 * time.Second)
		var game = Game{ID: i, ticketMin: 1, ticketMax: 90, chatHub: NewHub()}
		game.genDraw()
		db.AddGame(i, game, startTime, s)
		for p := 0; p < numP; p++ {
			db.AddTicketPack((i*numP)+p, i, p, 20)
		}

	}
	// take now  create one game every 3 mins for 20 games

}

func (d *Db) AddGame(ID int, gameData Game, startTime time.Time, s *Schedule) {
	var gameEntry = dbGameEntry{ID: ID, gameData: gameData, startTime: startTime}
	go d.scheduleGame(s, gameEntry)
	d.games = append(d.games, gameEntry)
}

func (d *Db) AddPlayer(ID int, player Player) {
	var playerEntry = dbPlayerEntry{ID: ID, playerData: player}
	d.players = append(d.players, playerEntry)
}

func (d *Db) AddTicketPack(ID int, gameID int, playerID int, num int) {

	for _, g := range d.games {
		if g.ID == gameID {
			tpe := dbTicketPackEntry{ID: ID, ticketData: g.gameData.GenTicketPack(num, playerID)}
			d.tickets = append(d.tickets, tpe)
			break
		}
	}

}

func (d *Db) GetGame(ID int) (Game, error) {
	var game Game
	for _, g := range d.games {
		if g.ID == ID {
			return g.gameData, nil
		}
	}
	return game, errors.New("Game not found")
}

func (d *Db) GetGames() []Game {
	var games []Game
	for _, g := range d.games {
		games = append(games, g.gameData)
	}
	return games
}

func (d *Db) GetPlayer(ID int) (Player, []TicketPack, error) {
	var p Player
	var tp []TicketPack
	for _, player := range d.players {
		if player.playerData.ID == ID {
			for _, t := range d.tickets {
				if t.ticketData.playerID == ID {
					tp = append(tp, t.ticketData)
				}

			}
			return player.playerData, tp, nil
		}

	}
	return p, tp, errors.New("Player Not Found")
}

func (d *Db) GetPlayerInGame(pID int, gID int) (Player, error) {
	var player Player
	ps := d.getPlayers(gID)
	for _, p := range ps {
		if p.ID == pID {
			return p, nil
		}
	}
	return player, errors.New("Player not in game")
}

func (d *Db) scheduleGame(s *Schedule, g dbGameEntry) {
	fmt.Printf("Game ID: %v Start time: %v\n", g.ID, g.startTime)
	t := time.NewTimer(time.Until(g.startTime))
	<-t.C
	g.gameData.Init(d.getPlayers(g.ID))
	s.startedGames <- g.gameData

}

// Get all players in a game
func (d *Db) getPlayers(ID int) []Player {
	pIDs := make(map[int]TicketPack)
	var ps []Player
	tp := d.getTicketPacks(ID)
	for _, t := range tp {
		pIDs[t.playerID] = t
	}
	for pID, tp := range pIDs {
		for _, p := range d.players {
			if p.playerData.ID == pID {
				p.playerData.tickets = append(p.playerData.tickets, tp)
				ps = append(ps, p.playerData)
			}
		}

	}
	return ps
}

func (d *Db) getTicketPacks(ID int) []TicketPack {
	var rTP []TicketPack
	for _, tp := range d.tickets {
		if tp.ticketData.gameID == ID {
			rTP = append(rTP, tp.ticketData)
		}

	}
	return rTP
}
