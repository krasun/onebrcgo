package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type StationAverage struct {
	sum   float64
	count int
	min   float64
	max   float64
}

func main() {
	file, err := os.Open("./measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	measurements := make(map[string]*StationAverage)
	stations := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		parts := strings.Split(row, ";")
		station := parts[0]
		temperature, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		_, exists := measurements[station]
		if exists {
			measurements[station].max = math.Max(measurements[station].max, temperature)
			measurements[station].min = math.Min(measurements[station].min, temperature)
			measurements[station].sum += temperature
			measurements[station].count += 1
		} else {
			measurements[station] = &StationAverage{
				max:   temperature,
				min:   temperature,
				sum:   temperature,
				count: 1,
			}
			stations = append(stations, station)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Strings(stations)

	fmt.Print("{")
	for i, station := range stations {
		average := measurements[station]
		fmt.Printf("%s=%.1f/%.1f/%.1f", station, average.min, average.sum/float64(average.count), average.max)
		if i != (len(stations) - 1) {
			fmt.Print(", ")
		}
	}
	fmt.Print("}")
}
