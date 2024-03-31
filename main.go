package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"sync"
)

type StationMetrics struct {
	sum   float64
	count int
	min   float64
	max   float64
}

func parse(row string) (string, float64, error) {
	for p, r := range row {
		if r == ';' {
			station, data := row[:p], row[p+len(";"):]
			temperature, err := strconv.ParseFloat(data, 64)
			if err != nil {
				return "", 0, fmt.Errorf("failed to parse the temperature \"%s\" as a number", data)
			}

			return station, temperature, nil
		}
	}

	return "", 0, fmt.Errorf("failed to locate \";\" in \"%s\"", row)
}

func compute(lines chan string) (map[string]*StationMetrics, []string) {
	measurements := make(map[string]*StationMetrics)
	stations := make([]string, 0)

	for line := range lines {
		station, temperature, err := parse(line)
		if err != nil {
			log.Fatal(err)
		}

		s, exists := measurements[station]
		if exists {
			s.max = math.Max(s.max, temperature)
			s.min = math.Min(s.min, temperature)
			s.sum += temperature
			s.count += 1
		} else {
			measurements[station] = &StationMetrics{
				max:   temperature,
				min:   temperature,
				sum:   temperature,
				count: 1,
			}
			stations = append(stations, station)
		}
	}

	sort.Strings(stations)

	return measurements, stations
}

func findNextNewLinePosition(file *os.File, startPosition int64) (int64, error) {
	_, err := file.Seek(startPosition, io.SeekStart)
	if err != nil {
		return 0, fmt.Errorf("failed to set the offset to %d: %w", startPosition, err)
	}

	var buf [1]byte
	for {
		n, err := file.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				return 0, io.EOF
			}

			return 0, err
		}

		startPosition += int64(n)

		if buf[0] == '\n' {
			return startPosition, nil
		}
	}
}

func readFile(scanner *bufio.Scanner, lines chan string) {
	for scanner.Scan() {
		lines <- scanner.Text()
	}
}

func createScanners(filePath string, chunkNumber int) ([]*bufio.Scanner, func() error, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get the file information: %w", err)
	}

	files := make([]*os.File, chunkNumber)
	closeFiles := func() error {
		errs := make([]error, 0)
		for _, file := range files {
			if file == nil {
				continue
			}

			err := file.Close()
			if err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return fmt.Errorf("failed to close files: %w", errors.Join(errs...))
		}

		return nil
	}

	fileSize := fileInfo.Size()
	chunkSize := fileSize / int64(chunkNumber)
	var startPosition int64
	scanners := make([]*bufio.Scanner, chunkNumber)
	for i := 0; i < chunkNumber; i++ {
		fileName := filePath

		file, err := os.Open(fileName)
		if err != nil {
			closeFiles()

			return nil, nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}

		nextPosition, err := findNextNewLinePosition(file, startPosition+chunkSize)
		if i == chunkNumber-1 && err == io.EOF {
			nextPosition = fileSize
		} else if err != nil {
			closeFiles()

			return nil, nil, fmt.Errorf("failed to find the closest new line for position %d: %w", startPosition, err)
		}

		files[i] = file
		scanners[i] = bufio.NewScanner(io.NewSectionReader(file, startPosition, nextPosition-startPosition))

		startPosition = nextPosition
	}

	return scanners, closeFiles, nil
}

func main() {
	// filePath := os.Args[1:][0]
	filePath := "./measurements.txt"

	chunkNumber := 2
	scanners, closeFiles, err := createScanners(filePath, chunkNumber)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFiles()

	var readWaitGroup sync.WaitGroup
	lines := make(chan string, 1000000)

	for _, scanner := range scanners {
		readWaitGroup.Add(1)
		go func(scanner *bufio.Scanner) {
			defer readWaitGroup.Done()

			readFile(scanner, lines)

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}(scanner)
	}

	var measurements map[string]*StationMetrics
	var stations []string
	var computeWaitGroup sync.WaitGroup
	computeWaitGroup.Add(1)
	go func() {
		defer computeWaitGroup.Done()
		measurements, stations = compute(lines)
	}()

	readWaitGroup.Wait()
	close(lines)

	computeWaitGroup.Wait()

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
