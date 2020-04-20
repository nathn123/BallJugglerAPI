package main

import (
	"fmt"
	"time"
)

type Schedule struct {
	startedGames chan Game
}

func (s *Schedule) Init() {
	fmt.Println("Schedule init")
	s.startedGames = make(chan Game)
	for i := 0; i < 5; i++ {
		go worker(i, s.startedGames)
	}

}

// get games
// then start the games on time

func worker(id int, games <-chan Game) {
	fmt.Printf("Worker initialised ID:%v\n", id)
	for g := range games {
		fmt.Printf("Worker running ID:%v\n", id)
		playGame(g)
		fmt.Printf("Worker complete ID:%v\n", id)
	}
}

func playGame(g Game) {
	lineDone := false
	go g.chatHub.run()
	fmt.Println("FIRST BROADCAST")
	g.chatHub.broadcast <- []byte("Game is starting")
	fmt.Println("FIRST BROADCAST END")
	for i, d := range g.draw {
		g.chatHub.broadcast <- []byte(fmt.Sprintf("The next number in the draw is: %v", d))
		fmt.Printf("GameID: %v Draw num:%v  The next number is %v\n", g.ID, i, d)
		if !lineDone && i == g.result.line.turn {
			lineDone = true
			g.chatHub.broadcast <- []byte(fmt.Sprintf("The first line has been won by %v", g.result.line.playerID))
			fmt.Printf("GameID: %v The first line has been won by %v\n", g.ID, g.result.line.playerID)
		}
		if lineDone && i == g.result.house.turn {
			g.chatHub.broadcast <- []byte(fmt.Sprintf("The house has been won by %v", g.result.house.playerID))
			fmt.Printf("GameID: %v The House has been won by %v\n", g.ID, g.result.house.playerID)
			g.chatHub.broadcast <- []byte(fmt.Sprintf("The Game is over Congratulations to %v", g.result.house.playerID))
			fmt.Printf("The Game is over Congratulations to %v\n", g.result.house.playerID)
			// TODO close chat window
			return
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
