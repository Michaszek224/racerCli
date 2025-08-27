/*to add :
- deleting racer tail ->done
- borders ->done
- handling collision ->done
- final destination->done
- some kind of points ->done
- game restart after lose or win ->done
- table score instead of admin panel ->done
- change speed according to level diff -> done
- local db to store best scores -> done
*/

/*Bug to fix:
- after colission with border old destination do not remove itself and barrier do not regenerate -> fixed
*/

package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/eiannone/keyboard"
	_ "github.com/mattn/go-sqlite3"
)

const (
	height = 20
	width  = 50
	offset = 2
)

type Point struct {
	x, y int
}

var keys chan rune
var barriers []Point
var racer Point
var destination Point
var dir Point
var score int
var oldRacer Point
var timeValue int

func interuptFunc() {
	//Getting input from user
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

func generateMap() {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if i == 0 || i == height-1 || j == 0 || j == width-1 {
				// Border
				fmt.Printf("\033[%d;%dH#", i+offset, j+offset)
				barriers = append(barriers, Point{j + offset, i + offset})
			}
		}
	}
}

func generateDestination() {
	destination = Point{
		rand.IntN(width-2) + offset + 1,
		rand.IntN(height-2) + offset + 1,
	}
	fmt.Printf("\033[%d;%dH?", destination.y, destination.x)
}

func cleanMap() {
	for i := 1; i < height-1; i++ {
		for j := 1; j < width-1; j++ {
			fmt.Printf("\033[%d;%dH ", i+offset, j+offset)
		}
	}
}

func newRound(currentScore int) {
	generateMap()
	cleanMap()
	racer = Point{5, 5}
	oldRacer = Point{5, 5}
	dir = Point{0, 0}
	generateDestination()
	score = currentScore
}

func initDb() *sql.DB {
	db, err := sql.Open("sqlite3", "scores.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bestScore(
			id INTEGER PRIMARY KEY CHECK (id=1),
			score INTEGER NOT NULL,
			date TEXT NOT NULL
			);
	`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		INSERT OR IGNORE INTO bestScore (id, score, date)
		VALUES (1, 0, datetime('now'));
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getDb(db *sql.DB) (int, string) {
	var score int
	var date string
	err := db.QueryRow("SELECT score, date FROM bestScore where id = 1").Scan(&score, &date)
	if err != nil {
		return 0, time.Now().Format("2006-01-02")
	}
	return score, date
}

func updateDb(db *sql.DB, newScore int) {
	_, err := db.Exec("UPDATE bestScore SET score = ?, date = datetime('now') WHERE id = 1", newScore)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db := initDb()
	defer db.Close()

	bestScore, bestDate := getDb(db)

	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	keys = make(chan rune)

	//base location of racer
	racer = Point{5, 5}
	oldRacer = Point{5, 5}

	//base direction of racer
	dir = Point{0, 0}

	//base value of stop timer
	timeValue = 70

	//clearing temirnal
	fmt.Print("\033[H\033[2J")

	//launching interputing func
	go interuptFunc()

	//Hidding the cursor
	fmt.Print("\033[?25l")
	//Restor the cursor after
	defer fmt.Print("\033[?25h")

	//Making the borders
	generateMap()

	//Radnom destination
	generateDestination()

	score = 0

	for {
		//Handling collision
		if slices.Contains(barriers, racer) {
			newRound(0)
			continue
		}

		//Handling win condition
		if racer.x == destination.x && racer.y == destination.y {
			if score+1 > bestScore {
				updateDb(db, score+1)
				bestScore, bestDate = getDb(db)
			}
			newRound(score + 1)
			continue
		}

		//Clean old racer
		fmt.Printf("\033[%d;%dH ", oldRacer.y, oldRacer.x)

		// Move racer
		racer = Point{racer.x + dir.x, racer.y + dir.y}

		oldRacer = racer

		// Draw racer
		fmt.Printf("\033[%d;%dH@", racer.y, racer.x)

		//Clearing old logs
		fmt.Printf("\033[%d;%dH            ", 9, 100)
		fmt.Printf("\033[%d;%dHA                                ", 10, 100)
		fmt.Printf("\033[%d;%dHA                                ", 11, 100)

		// Logging
		fmt.Printf("\033[%d;%dHScore Table", 9, 100)
		fmt.Printf("\033[%d;%dHHighest Score: %d (%s)", 10, 100, bestScore, bestDate)
		fmt.Printf("\033[%d;%dHCurrent Score: %d", 11, 100, score)

		// Handle input
		select {
		case char := <-keys:
			switch char {
			case 'w':
				timeValue = 70 - score //using + score to increase diff
				dir = Point{0, -1}
			case 's':
				timeValue = 70 - score
				dir = Point{0, 1}
			case 'a':
				timeValue = 50 - score
				dir = Point{-1, 0}
			case 'd':
				timeValue = 50 - score
				dir = Point{1, 0}
			case 'x':
				dir = Point{0, 0}
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
