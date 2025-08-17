package main

import (
	"fmt"
	"log"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	height = 20
	width  = 50
)

type Point struct {
	x, y int
}

func interuptFunc() {
	//Function for getting input from user
	for {
		char, _, err := keyboard.GetKey()
		if err != nil {
			close(keys)
			return
		}
		keys <- char
	}
}

var keys chan rune

/*
to add :
- deleting racer tail ->done
- borders ->done
- map of points where user was
- random bariers
- final destination
- game go to 0 when user go into wall or previous point
*/

func main() {
	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	keys = make(chan rune)

	//base location of racer
	racer := Point{5, 5}
	oldRacer := Point{5, 5}

	//base direction of racer
	dir := Point{0, 0}

	//base value of stop timer
	timeValue := 50

	//clearing temirnal
	fmt.Print("\033[H\033[2J")
	//clear after the game is over
	defer fmt.Print("\033[H\033[2J")

	//launching interputing func
	go interuptFunc()

	//Hidding the cursor
	fmt.Print("\033[?25l")
	//Restor the cursor after
	defer fmt.Print("\033[?25h")

	//Making the borders
	for i := range height {
		for j := range width {
			fmt.Printf("\033[%d;%dH#", i+2, j+2)
		}
	}
	for i := 1; i < height-1; i++ {
		for j := 1; j < width-1; j++ {
			fmt.Printf("\033[%d;%dH ", i+2, j+2)
		}
	}

	for {
		//Clean old racer
		fmt.Printf("\033[%d;%dH ", oldRacer.y, oldRacer.x)

		// Move racer
		racer = Point{racer.x + dir.x, racer.y + dir.y}
		oldRacer = racer

		// Draw racer
		fmt.Printf("\033[%d;%dH`", racer.y, racer.x)

		// Handle input
		select {
		case char := <-keys:
			switch char {
			case 'w':
				timeValue = 70
				dir = Point{0, -1}
			case 's':
				timeValue = 70
				dir = Point{0, 1}
			case 'a':
				timeValue = 50
				dir = Point{-1, 0}
			case 'd':
				timeValue = 50
				dir = Point{1, 0}
			case 'p':
				return
			}
		default:
			// if no input keep doing the same
		}

		//different sleep time depending on direction
		time.Sleep(time.Duration(timeValue) * time.Millisecond)

	}
}
