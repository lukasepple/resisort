/*
	resisort - a resistor sorting helper tool
	(c) 2014 Lukas Epple <sternenseemann @ lukasepple.de>

	This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"unicode"
)

const (
	TO_BE_DETERMINATED              = -1
	DEFAULT_RESISTORS_PER_CONTAINER = 10
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Too few arguments, try %s --help\n", os.Args[0])
		os.Exit(1)
	}

	// arg parsing, not too elegant
	// but simple
	var filename string
	resistors_per_container := 0
	container_count := 0

	read_rpc := false
	read_container_count := false

	for _, Arg := range os.Args {
		switch Arg {
		case "--help":
			PrintHelp()
			break
		case "--resistors-per-container":
			read_rpc = true
			break
		case "--containers":
			read_container_count = true
			break
		default:
			if read_rpc {
				var err error
				resistors_per_container, err = strconv.Atoi(Arg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: expected integer, got %s\n", os.Args[0], Arg)
				}
				read_rpc = false
			} else if read_container_count {
				resistors_per_container = TO_BE_DETERMINATED
				var err error
				container_count, err = strconv.Atoi(Arg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: expected integer, got %s\n", os.Args[0], Arg)
				}
				read_container_count = false
			} else {
				filename = Arg
			}
		}
	}

	resistors := ReadResistors(filename)


	if resistors_per_container == 0 {
		resistors_per_container = DEFAULT_RESISTORS_PER_CONTAINER
	} else if resistors_per_container == TO_BE_DETERMINATED {
		resistors_per_container = int(math.Ceil(float64(len(resistors)) / float64(container_count)))
	}

	if container_count == 0 {
		container_count = int(math.Ceil(float64(len(resistors)) / float64(resistors_per_container)))
	}

	if container_count > len(resistors) {
		container_count = len(resistors)
	}

	fmt.Printf("In order to sort %d resistors you can use %d container(s) with up to %d resistor(s) each!\n", len(resistors), container_count, resistors_per_container)

	for i := 0; i < container_count; i++ {
		// calculate the indices of the smallest and the biggest resistor
		// in the group
		firstresistor := i * resistors_per_container
		lastresistor  := firstresistor + resistors_per_container - 1

		if lastresistor >= len(resistors) {
			lastresistor = len(resistors) - 1
		}

		fmt.Printf("%dth Container: %10s - %10s\n", i + 1, FormatResistorValue(resistors[firstresistor]), FormatResistorValue(resistors[lastresistor]))
	}
}

func PrintHelp() {
	fmt.Fprintf(os.Stderr, "%s [(--containers | --resistors-per-container) int] filename\n", os.Args[0])
	os.Exit(0)
}

func ReadResistors(filename string) []int {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Could not open file: %s", os.Args[0], err)
	}
	defer f.Close()

	resistors := make([]int, 0) // enough to hold a E12 set
	ResistorScanner := bufio.NewScanner(f)

	for ResistorScanner.Scan() {
		resistor := ResistorScanner.Text()
		convertable := ""
		multiplier := 1.0

		for _, char := range resistor {
			if unicode.IsDigit(char) {
				convertable += string(char)
			} else if char == 'K' {
				multiplier = 1000.0
				convertable += "."
			} else if char == 'M' {
				multiplier = 1000000.0
				convertable += "."
			} else if char == '.' || char == ',' {
				convertable += "."
			} else if char == 'Ω' || char == 'R' {
				convertable += "."
			}
		}

		if convertable[len(convertable)-1] == '.' {
			convertable = convertable[:len(convertable)-1]
		}

		resistor_value, err := strconv.ParseFloat(convertable, 64)
		resistor_value = resistor_value * multiplier

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: expected to calculate an float but calculated %s! %s\n", os.Args[0], convertable, err)
		}

		resistors = append(resistors, int(resistor_value))
	}

	if err := ResistorScanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: Error whilst Reading: %s", os.Args[0], err)
	}

	sort.Ints(resistors)

	return resistors
}

func FormatResistorValue(value int) string {
	if value >= 1000000 {
		// megaohm
		return fmt.Sprintf("%.1fMΩ", float64(value) / 1000000.0)
	} else if value >= 1000 {
		// kiloohm
		return fmt.Sprintf("%.1fKΩ", float64(value) / 1000.0)
	} else {
		return fmt.Sprintf("%dΩ", value)
	}
}
