package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
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

type Garden struct {
	Size     Pos
	RobotPos Pos
	// -1 == empty
	RobotCarryMass int
	// -1 == empty
	Objects         map[rune]int
	ObjectLocations map[Pos]rune
}

var (
	out        = colorable.NewColorableStdout()
	noDelay    = false
	step       = false
	jsonOutput = false
)

func (w *Garden) print(clear bool) {
	if jsonOutput {
		return
	}
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

func readGarden(path string) *Garden {
	gardenFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open garden file: %s", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(gardenFile))

	garden := &Garden{
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
			fmt.Fprintf(os.Stderr, "Invalid garden file: sytax error line %d", line)
			os.Exit(1)
		}
		r, _ := utf8.DecodeRuneInString(strings.TrimSpace(parts[0]))
		if r == utf8.RuneError {
			fmt.Fprintf(os.Stderr, "Invalid garden file: sytax error line %d", line)
			os.Exit(1)
		}
		parts[1] = strings.TrimSpace(parts[1])
		if parts[1] == "" {
			garden.Objects[r] = -1
		} else {
			mass, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid garden file: sytax error line %d", line)
				os.Exit(1)
			}
			garden.Objects[r] = mass
		}
	}

	garden.ObjectLocations = make(map[Pos]rune, len(garden.Objects))
	row := 0
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if garden.Size.X != 0 && len(text) != garden.Size.X {
			break
		}
		column := 0
		for _, r := range text {
			if r == '$' {
				garden.RobotPos = Pos{
					X: column,
					Y: row,
				}
			} else if r >= 'A' && r <= 'Z' {
				garden.ObjectLocations[Pos{
					X: column,
					Y: row,
				}] = r
			}
			column++
		}
		garden.Size.X = column
		row++
	}
	garden.Size.Y = row
	return garden
}

func (w *Garden) moveRobot(dx, dy int) {
	if w.RobotPos.X+dx >= w.Size.X || w.RobotPos.X+dx < 0 || w.RobotPos.Y+dy >= w.Size.Y || w.RobotPos.Y+dy < 0 {
		fmt.Fprintln(os.Stderr, "Cannot move out of map")
		os.Exit(2)
	}
	w.RobotPos.X += dx
	w.RobotPos.Y += dy
}

func (w *Garden) pickup() {
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

func (w *Garden) put() {
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

func (w *Garden) run(input io.Reader) (commandCount int) {
	if !noDelay {
		w.print(true)
	}
	delay := (-0.065*float64(len(w.Objects)) + 0.5) * 1000
	if delay <= 10 {
		delay = 10
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
				fmt.Fprintln(os.Stderr, "Invalid input (command repeat argument is not a valid number)")
				os.Exit(1)
			}
		}
		if step {
			fmt.Print("Press enter to continue...")
			bufio.NewScanner(os.Stdin).Scan()
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
			commandCount++
			if !noDelay {
				if delay > 50 || i%2 == 0 {
					w.print(true)
				}
				time.Sleep(delayDuration)
			}
		}
		if !noDelay && (delay <= 50 || count%2 == 0) {
			w.print(true)
			fmt.Println(text)
		}
		if w.check() {
			return
		}
	}
	return
}

func (w *Garden) check() bool {
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

func inputInt(prompt string, min, max int) int {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		scanner.Scan()
		text := scanner.Text()
		num, err := strconv.Atoi(strings.TrimSpace(text))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Value must be a number")
			continue
		}
		if num < min || num > max {
			fmt.Fprintln(os.Stderr, "Value must be between", min, "and", max)
			continue
		}
		return num
	}
}

func generateGarden(path string, width, height, areaCount int) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create garden file:", err)
		os.Exit(1)
	}
	defer file.Close()

	areas := make(map[rune]int, areaCount)
	areaLocations := make(map[int]rune, areaCount)
	for i := 0; i < areaCount; i++ {
		areas[rune('A'+i)] = random.Intn(areaCount*3) * 10
		for {
			location := random.Intn(width * height)
			if _, ok := areaLocations[location]; !ok {
				areaLocations[location] = rune('A' + i)
				break
			}
		}
	}

	emptyCount := random.Intn(areaCount/2) + 1
	for i := 0; i < emptyCount; i++ {
		areas[rune('A'+random.Intn(areaCount-1))] = -1
	}

	var robotPos int
	for {
		robotPos = random.Intn(width * height)
		if _, ok := areaLocations[robotPos]; !ok {
			break
		}
	}

	for r, weight := range areas {
		if weight >= 0 {
			fmt.Fprintf(file, "%s=%d\n", string(r), weight)
		} else {
			fmt.Fprintf(file, "%s=\n", string(r))
		}
	}
	fmt.Fprintln(file)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if r, ok := areaLocations[y*width+x]; ok {
				fmt.Fprint(file, string(r))
			} else if y*width+x == robotPos {
				fmt.Fprint(file, "$")
			} else {
				fmt.Fprint(file, ".")
			}
		}
		if y < height-1 {
			fmt.Fprintln(file)
		}
	}
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	var generate bool
	flag.BoolVar(&generate, "generate", false, "Generate a random garden")
	var size string
	flag.StringVar(&size, "size", "", "The size of the garden to generate (format: widthxheight (e.g. 64x32 => min: 4, max: 128) or random), default: prompt user")
	var areaCountStr string
	flag.StringVar(&areaCountStr, "area-count", "", "The number of areas to generate (2-26 or random), default: prompt user")
	flag.BoolVar(&noDelay, "no-delay", false, "Disable delay between steps")
	flag.BoolVar(&step, "step", false, "Prompt to press enter before every step")
	var input string
	flag.StringVar(&input, "input", "", "File path to file containing commands")
	var seedStr string
	flag.StringVar(&seedStr, "seed", "", "PRNG seed for generator (integer)")
	flag.BoolVar(&jsonOutput, "json", false, "Output stats in JSON format")
	flag.Parse()

	if jsonOutput {
		noDelay = jsonOutput
	}

	if seedStr != "" {
		seed, err := strconv.Atoi(seedStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Seed must be an integer")
			os.Exit(1)
		}
		random = rand.New(rand.NewSource(int64(seed)))
	}

	if flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "The garden file must be the last argument.")
		fmt.Fprintf(os.Stderr, "USAGE: %s [OPTIONS] <garden_file>\n", os.Args[0])
		os.Exit(1)
	}

	gardenFile := flag.Arg(0)
	if gardenFile == "" {
		fmt.Fprintf(os.Stderr, "USAGE: %s [OPTIONS] <garden_file>\n", os.Args[0])
		os.Exit(1)
	}

	if generate {
		var width, height, areaCount int
		if size != "" {
			if size == "random" {
				width = random.Intn(128-4) + 4
				height = random.Intn(128-4) + 4
			} else {
				parts := strings.Split(size, "x")
				if len(parts) != 2 {
					fmt.Fprintln(os.Stderr, "Invalid size format")
					os.Exit(1)
				}
				var err error
				width, err = strconv.Atoi(parts[0])
				if err != nil {
					fmt.Fprintln(os.Stderr, "Invalid size format")
					os.Exit(1)
				}
				height, err = strconv.Atoi(parts[1])
				if err != nil {
					fmt.Fprintln(os.Stderr, "Invalid size format")
					os.Exit(1)
				}
				if width < 4 || width > 128 {
					fmt.Fprintln(os.Stderr, "Invalid width: min: 4, max: 128")
					os.Exit(1)
				}
				if height < 4 || height > 128 {
					fmt.Fprintln(os.Stderr, "Invalid height: min: 4, max: 128")
					os.Exit(1)
				}
			}
		} else {
			width = inputInt("Width (4-128): ", 4, 128)
			height = inputInt("Height (4-128): ", 4, 128)
		}

		maxAreaCount := width*height - 1
		if maxAreaCount > 26 {
			maxAreaCount = 26
		}

		if areaCountStr != "" {
			if areaCountStr == "random" {
				areaCount = random.Intn(maxAreaCount-2) + 2
			} else {
				var err error
				areaCount, err = strconv.Atoi(areaCountStr)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Invalid area count: not a number")
					os.Exit(1)
				}
				if areaCount < 2 || areaCount > 26 {
					fmt.Fprintln(os.Stderr, "Invalid area count: min: 2, max: 26")
					os.Exit(1)
				}
			}
		} else {
			areaCount = inputInt(fmt.Sprintf("Area count (2-%d): ", maxAreaCount), 2, maxAreaCount)
		}

		generateGarden(gardenFile, width, height, areaCount)
		return
	}

	if noDelay && step {
		fmt.Fprintln(os.Stderr, "Cannot enable -no-delay and -step at the same time")
		os.Exit(1)
	}
	if step && input == "" {
		fmt.Fprintln(os.Stderr, "Cannot enable -step if input is set to STDIN (use -input to specify an input file)")
		os.Exit(1)
	}

	garden := readGarden(gardenFile)

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

	commandCount := garden.run(file)
	if jsonOutput {
		type result struct {
			Success      bool `json:"success"`
			CommandCount int  `json:"commandCount"`
		}
		json.NewEncoder(os.Stdout).Encode(&result{
			Success:      garden.check(),
			CommandCount: commandCount,
		})
		return
	}
	garden.print(!noDelay)
	if garden.check() {
		fmt.Printf("Success! The garden is tidy. The robot executed %d commands.\n", commandCount)
	} else {
		fmt.Println("Failure! The objects are not sorted.")
	}
}
