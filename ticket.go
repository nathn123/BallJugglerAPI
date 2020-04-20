package main

import (
	"fmt"
	"math/rand"
	"time"
)

//Ticket used to store ticket data
type ticket struct {
	id    int
	rows  [3]map[int]int
	line  int
	house int
}

// TicketPack to store multiple tickets for a player per game
type TicketPack struct {
	playerID int
	gameID   int
	tickets  []ticket
	win      ticketWin
}

type ticketWin struct {
	line    int
	lineID  int
	house   int
	houseID int
}

//Gen used to generate a ticket pack
func (tp *TicketPack) Gen(count int, ticketMin int, ticketMax int, draw []int) {
	var tickets []ticket
	tp.win = ticketWin{100, -1, 100, -1}
	fmt.Println("TICKET COUNT:", count)
	fmt.Println("TICKET DRAW:", draw)
	for i := 0; i < count; i++ {

		t := ticket{id: i}
		t.gen(ticketMin, ticketMax)
		t.checkTicket(draw)
		if t.line < tp.win.line {
			tp.win.line = t.line
			tp.win.lineID = t.id
		}
		if t.house < tp.win.house {
			tp.win.house = t.house
			tp.win.houseID = t.id
		}
		tickets = append(tickets, t)
	}

}

//Gen Used to generate a ticket
func (t *ticket) gen(min int, max int) {
	for i := 0; i < 3; i++ {
		// generate 5 nums in array -- loop 5 times
		nums := make(map[int]int)
		for len(nums) < 5 {
			// gen a number between 1 - 90
			rand.Seed(time.Now().UnixNano())
			newNum := rand.Intn(max-min) + min
			// test against exist list
			if CheckNum(nums, newNum) {
				// if passed store in list -- if failed try again
				nums[newNum] = newNum

			}
		}
		t.rows[i] = nums
	}
}

func (t *ticket) checkTicket(n []int) {
	row1, row2, row3 := 0, 0, 0
	lineDone := false
	for i, v := range n {
		if _, ok := t.rows[0][v]; ok {
			row1++
		}
		if _, ok := t.rows[1][v]; ok {
			row2++
		}
		if _, ok := t.rows[2][v]; ok {
			row3++
		}

		if (row1 == 5 || row2 == 5 || row3 == 5) && !lineDone {
			t.line = i + 1
			lineDone = true
		}

		if row1 == 5 && row2 == 5 && row3 == 5 {
			t.house = i + 1
			return
		}
	}
	return
}

// CheckNum ...
// check if number range exists
// to test divide by 10 cast to int  compare against list i.e  70 in list 71 new gen  7 == 7 failure
func CheckNum(l map[int]int, n int) bool {
	for _, v := range l {
		if int(v/10) == int(n/10) {
			return false
		}
	}
	return true
}

// PrintTicket used to print the ticket for debug
func PrintTicket(t ticket) {
	for i, v := range t.rows {
		fmt.Printf("Row: %v\n%v\n", i+1, v)
	}
}
