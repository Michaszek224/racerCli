/*
to add :
- deleting racer tail ->done
- borders ->done
- handling collision ->done
- random bariers
- final destination
*/

package main

import (
	"fmt"
	"log"
	"slices"
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

var keys chan rune
var barriers []Point

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

func endProgram(message string) {
	fmt.Print("\033[H\033[2J")
	fmt.Print("\033[?25h")
	fmt.Println(message)
}

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

	//launching interputing func
	go interuptFunc()

	//Hidding the cursor
	fmt.Print("\033[?25l")
	//Restor the cursor after
	defer fmt.Print("\033[?25h")

	//Making the borders
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if i == 0 || i == height-1 || j == 0 || j == width-1 {
				// Border
				fmt.Printf("\033[%d;%dH#", i+2, j+2)
				barriers = append(barriers, Point{j + 2, i + 2})
			}
		}
	}

	for {
		//Handling collision
		if slices.Contains(barriers, racer) {
			endProgram("You lose")
			break
		}

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
				endProgram("Endend using command")
				return
			}
		default:
			// if no input keep doing the same
		}

		//different sleep time depending on direction
		time.Sleep(time.Duration(timeValue) * time.Millisecond)

	}
}
