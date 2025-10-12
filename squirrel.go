package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type Entry struct {
	Events   []string `json:"events"`
	Squirrel bool     `json:"squirrel"`
}

func contains(events []string, event string) bool {
	for _, e := range events {
		if e == event {
			return true
		}
	}
	return false
}

func phi(array [4]float64) float64 {
	num := array[0]*array[3] - array[1]*array[2]
	sqr := math.Sqrt((array[0] + array[1]) * (array[2] + array[3]) *
		(array[0] + array[2]) * (array[1] + array[3]))
	if sqr == 0 {
		return 0
	}
	return num / sqr
}

func analyzeJournal(sample string) (string, float64, string, float64, error) {
	data, err := os.ReadFile(sample)
	if err != nil {
		return "", 0, "", 0, err
	}

	var journal []Entry
	if err := json.Unmarshal(data, &journal); err != nil {
		return "", 0, "", 0, err
	}

	eventSet := make(map[string]bool)
	for _, entry := range journal {
		for _, event := range entry.Events {
			eventSet[event] = true
		}
	}

	var maxEvent, minEvent string
	maxPhi := -1.0
	minPhi := 1.0

	for event := range eventSet {
		var array [4]float64
		for _, entry := range journal {
			has := contains(entry.Events, event)
			squirrel := entry.Squirrel

			switch {
			case has && squirrel:
				array[0]++
			case has && !squirrel:
				array[1]++
			case !has && squirrel:
				array[2]++
			case !has && !squirrel:
				array[3]++
			}
		}

		corr := phi(array)
		if corr > maxPhi {
			maxPhi = corr
			maxEvent = event
		}
		if corr < minPhi {
			minPhi = corr
			minEvent = event
		}
	}

	return maxEvent, maxPhi, minEvent, minPhi, nil
}

func main() {
	maxEvent, maxPhi, minEvent, minPhi, err := analyzeJournal("sample.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Most positively correlated event: %s (%.4f)\n", maxEvent, maxPhi)
	fmt.Printf("Most negatively correlated event: %s (%.4f)\n", minEvent, minPhi)
}

