package main

import (
	"fmt"
	"math/rand"
	"time"
)

// WinStruct used to show the winners for line / house
type winStruct struct {
	id       int
	playerID int
	turn     int
}

// ResultStruct used to show the results for a draw
type ResultStruct struct {
	line  winStruct
	house winStruct
}

func draw(drawMin int, drawMax int) []int {
	rand.Seed(time.Now().UnixNano())
	rand.Seed(rand.Int63())
	var results []int
	for i := drawMin; i <= drawMax; i++ {
		fmt.Println(i)
		results = append(results, i)
	}
	rand.Shuffle(len(results), func(i int, j int) {
		results[i], results[j] = results[j], results[i]
	})
	return results
}

//FindWinner returns the ticket id for the first line and first house
func FindWinner(t map[int]ticketWin) ResultStruct {
	return ResultStruct{findFirst(t, true), findFirst(t, false)}
}

// findFirst used to find the first line or house
func findFirst(ts map[int]ticketWin, line bool) winStruct {
	var id int
	var turn = 91
	var playerID = 0
	for i, t := range ts {
		if line {
			if t.line < turn {
				turn = t.line
				id = t.lineID
				playerID = i
			}
		} else {
			if t.house < turn {
				turn = t.house
				id = t.houseID
				playerID = i
			}
		}
	}
	return winStruct{id, playerID, turn}
}
