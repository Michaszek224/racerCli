package main

import (
	"fmt"
	"log"
	"time"

	"github.com/eiannone/keyboard"
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

func main() {
	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	keys = make(chan rune)

	//base location of racer
	racer := Point{5, 5}

	//base direction of racer
	dir := Point{1, 0}

	//base value of stop timer
	timeValue := 50

	//clearing temirnal
	fmt.Print("\033[H\033[2J")

	//launching interputing func
	go interuptFunc()

	for {
		// Draw racer
		fmt.Printf("\033[%d;%dH#", racer.y, racer.x)

		// Move racer
		racer = Point{racer.x + dir.x, racer.y + dir.y}

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
