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
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"unicode"
	humanize "github.com/dustin/go-humanize"
)

type Container struct {
	Lowerbound int
	Upperbound int
}

const (
	NOT_SPECIFIED = 0
)

var filename string
var resistors_per_container int
var container_count int

func init() {
	flag.IntVar(&resistors_per_container, "resistors-per-container", NOT_SPECIFIED, "How many resistors should be in one container maximally?")
	flag.IntVar(&container_count, "containers", NOT_SPECIFIED, "Count of containers you can use")
	flag.StringVar(&filename, "file", "", "The file to read from")
}

func main() {
	flag.Parse()

	if xor(container_count > NOT_SPECIFIED, resistors_per_container > NOT_SPECIFIED) {
		resistors := ReadResistors(filename)
		var sorted []Container

		sorted, resistors_per_container, container_count = CalculateSorting(resistors, resistors_per_container,
			container_count)

		fmt.Printf("In order to sort %d resistors you can use %d container(s) with up to %d resistor(s) each!\n",
			len(resistors), container_count, resistors_per_container)

		for index, container := range sorted {
			fmt.Printf("%5s Container: %10s - %10s\n", humanize.Ordinal(index + 1),
				FormatResistorValue(container.Lowerbound), FormatResistorValue(container.Upperbound))
		}

	} else {
		fmt.Fprintf(os.Stderr, "%s: Either specify --containers or --resistors-per-container\n", os.Args[0])
		os.Exit(1)
	}
}

func ReadResistors(filename string) []int {
	var f *os.File
	var err error
	if filename == "" {
		f = os.Stdin
	} else {
		f, err = os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Could not open file: %s", os.Args[0], err)
		}
		defer f.Close()
	}

	resistors := make([]int, 0)
	ResistorScanner := bufio.NewScanner(f)

	for ResistorScanner.Scan() {
		resistor_value := ParseResistorValue(ResistorScanner.Text())
		resistors = append(resistors, resistor_value)
	}

	if err := ResistorScanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: Error whilst Reading: %s", os.Args[0], err)
	}

	sort.Ints(resistors)

	return resistors
}

func CalculateSorting(resistors []int, resistors_per_container int, container_count int) (sorted []Container, new_rpc int, new_cc int) {

	if resistors_per_container <= NOT_SPECIFIED {
		new_rpc = int(math.Ceil(float64(len(resistors)) / float64(container_count)))
	} else {
		new_rpc = resistors_per_container
	}

	// always (re)calculate the container count because
	// we can't use floating point resistors!
	new_cc = int(math.Ceil(float64(len(resistors)) / float64(new_rpc)))

	if container_count > len(resistors) {
		new_cc = len(resistors)
	}

	sorted = make([]Container, 0)

	for i := 0; i < new_cc; i++ {
		// calculate the indices of the smallest and the biggest resistor
		// in the group
		var container Container
		firstresistor := i * new_rpc
		lastresistor := firstresistor + new_rpc - 1

		if lastresistor >= len(resistors) {
			lastresistor = len(resistors) - 1
		}

		container.Lowerbound = resistors[firstresistor]
		container.Upperbound = resistors[lastresistor]

		sorted = append(sorted, container)
	}

	return sorted, new_rpc, new_cc
}

func ParseResistorValue(resistor string) int {
	convertable := ""
	multiplier := 1.0

	for _, char := range resistor {
		if unicode.IsDigit(char) {
			convertable += string(char)
		} else if char == 'K' || char == 'k' {
			multiplier = 1000.0
			convertable += "."
		} else if char == 'M' || char == 'm' {
			multiplier = 1000000.0
			convertable += "."
		} else if char == '.' || char == ',' {
			convertable += "."
		} else if char == 'Ω' || char == 'R' || char == 'r' {
			break
		}
	}

	if convertable[len(convertable)-1] == '.' {
		convertable = convertable[:len(convertable)-1]
	}

	resistor_value, err := strconv.ParseFloat(convertable, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: expected to calculate an float but calculated %s! %s\n", os.Args[0], convertable, err)
		os.Exit(1)
	}
	resistor_value = resistor_value * multiplier

	return int(resistor_value)
}

func FormatResistorValue(value int) string {
	if value >= 1000000 {
		// megaohm
		return fmt.Sprintf("%sMΩ", humanize.Ftoa(float64(value)/1000000.0))
	} else if value >= 1000 {
		// kiloohm
		return fmt.Sprintf("%sKΩ", humanize.Ftoa(float64(value)/1000.0))
	} else {
		return fmt.Sprintf("%dΩ", value)
	}
}

func xor(a bool, b bool) bool {
	return (a || b) && !(a && b)
}
