package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mattn/go-colorable"
)

type Pos struct {
	X int
	Y int
}

type World struct {
	Size     Pos
	RobotPos Pos
	// -1 == empty
	RobotCarryMass int
	// -1 == empty
	Objects         map[rune]int
	ObjectLocations map[Pos]rune
}

var (
	out     = colorable.NewColorableStdout()
	noDelay = false
)

func (w *World) print(clear bool) {
	if clear {
		fmt.Fprintf(out, "\033[H\033[2J")
	}
	for letter := 'A'; letter < 'A'+rune(len(w.Objects)); letter++ {
		mass := w.Objects[letter]
		if mass == -1 {
			fmt.Printf("%s: empty\n", string(letter))
		} else {
			fmt.Printf("%s: %d\n", string(letter), mass)
		}
	}
	fmt.Println()
	for y := 0; y < w.Size.Y; y++ {
		for x := 0; x < w.Size.X; x++ {
			if w.RobotPos.X == x && w.RobotPos.Y == y {
				fmt.Print("$")
			} else if r, ok := w.ObjectLocations[Pos{
				X: x,
				Y: y,
			}]; ok {
				fmt.Print(string(r))
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	if w.RobotCarryMass == -1 {
		fmt.Println("\nRobot:", "nothing")
	} else {
		fmt.Println("\nRobot:", w.RobotCarryMass)
	}
}

func readWorld(path string) *World {
	worldFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open world file: %s", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(worldFile))

	world := &World{
		Objects:        make(map[rune]int),
		RobotCarryMass: -1,
	}

	var line int
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			// end of object list
			break
		}
		parts := strings.Split(text, "=")
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "Invalid world file: sytax error line %d", line)
			os.Exit(1)
		}
		r, _ := utf8.DecodeRuneInString(strings.TrimSpace(parts[0]))
		if r == utf8.RuneError {
			fmt.Fprintf(os.Stderr, "Invalid world file: sytax error line %d", line)
			os.Exit(1)
		}
		parts[1] = strings.TrimSpace(parts[1])
		if parts[1] == "" {
			world.Objects[r] = -1
		} else {
			mass, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid world file: sytax error line %d", line)
				os.Exit(1)
			}
			world.Objects[r] = mass
		}
	}

	world.ObjectLocations = make(map[Pos]rune, len(world.Objects))
	row := 0
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if world.Size.X != 0 && len(text) != world.Size.X {
			break
		}
		column := 0
		for _, r := range text {
			if r == '$' {
				world.RobotPos = Pos{
					X: column,
					Y: row,
				}
			} else if r >= 'A' && r <= 'Z' {
				world.ObjectLocations[Pos{
					X: column,
					Y: row,
				}] = r
			}
			column++
		}
		world.Size.X = column
		row++
	}
	world.Size.Y = row
	return world
}

func (w *World) moveRobot(dx, dy int) {
	if w.RobotPos.X+dx >= w.Size.X || w.RobotPos.X+dx < 0 || w.RobotPos.Y+dy >= w.Size.Y || w.RobotPos.Y+dy < 0 {
		fmt.Fprintln(os.Stderr, "Cannot move out of map")
		os.Exit(2)
	}
	w.RobotPos.X += dx
	w.RobotPos.Y += dy
}

func (w *World) pickup() {
	if w.RobotCarryMass != -1 {
		fmt.Fprintln(os.Stderr, "Cannot pick up object: robot already carries an object")
		os.Exit(2)
	}
	if obj, ok := w.ObjectLocations[w.RobotPos]; ok {
		mass := w.Objects[obj]
		if mass == -1 {
			fmt.Fprintln(os.Stderr, "Cannot pick up object from empty area")
			os.Exit(2)
		}
		w.Objects[obj] = -1
		w.RobotCarryMass = mass
	} else {
		fmt.Fprintln(os.Stderr, "Cannot pick up object form grass area")
		os.Exit(2)
	}
}

func (w *World) put() {
	if w.RobotCarryMass == -1 {
		fmt.Fprintln(os.Stderr, "Cannot place object: robot doesn't carry any object")
		os.Exit(2)
	}
	if obj, ok := w.ObjectLocations[w.RobotPos]; ok {
		mass := w.Objects[obj]
		if mass != -1 {
			fmt.Fprintln(os.Stderr, "Cannot place object on occupied area")
			os.Exit(2)
		}
		w.Objects[obj] = w.RobotCarryMass
		w.RobotCarryMass = -1
	} else {
		fmt.Fprintln(os.Stderr, "Cannot place object on grass area")
		os.Exit(2)
	}
}

func (w *World) run(input io.Reader) {
	if !noDelay {
		w.print(true)
	}
	delay := (-0.065*float64(len(w.Objects)) + 0.5) * 1000
	if delay <= 0.01 {
		delay = 0.01
	}
	delayDuration := time.Duration(delay) * time.Millisecond

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		parts := strings.Split(text, " ")
		count := 1
		if len(parts) > 1 {
			var err error
			count, err = strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid input")
				os.Exit(1)
			}
		}
		for i := 0; i < count; i++ {
			switch parts[0] {
			case "fahre_norden":
				w.moveRobot(0, -1)
			case "fahre_osten":
				w.moveRobot(1, 0)
			case "fahre_sueden":
				w.moveRobot(0, 1)
			case "fahre_westen":
				w.moveRobot(-1, 0)
			case "gegenstand_aufheben":
				w.pickup()
			case "gegenstand_absetzen":
				w.put()
			default:
				fmt.Fprintln(os.Stderr, "Unknown command:", parts[0])
				os.Exit(1)
			}
			if !noDelay {
				w.print(true)
				time.Sleep(delayDuration)
			}
		}
		if w.check() {
			if noDelay {
				w.print(false)
			}
			fmt.Println("Success! The garden is tidy.")
			os.Exit(0)
		}
	}
}

func (w *World) check() bool {
	if w.RobotCarryMass != -1 {
		return false
	}
	previous := math.MaxInt
	for letter := 'A'; letter < 'A'+rune(len(w.Objects)); letter++ {
		if w.Objects[letter] > previous {
			return false
		}
		previous = w.Objects[letter]
	}
	return true
}

func main() {
	flag.BoolVar(&noDelay, "no-delay", false, "Disable delay between steps")
	var input string
	flag.StringVar(&input, "input", "", "File path to file containing commands")
	flag.Parse()

	worldFile := flag.Arg(0)
	if worldFile == "" {
		fmt.Fprintf(os.Stderr, "USAGE: %s [OPTIONS] <world_file>\n", os.Args[0])
		os.Exit(1)
	}

	world := readWorld(worldFile)

	file := os.Stdin
	if input != "" {
		var err error
		file, err = os.Open(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to open input file:", err)
			os.Exit(1)
		}
		defer file.Close()
	}

	world.run(file)
}
