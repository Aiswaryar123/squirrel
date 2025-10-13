package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
)

type Entry struct {
	Events   []string `json:"events"`
	Squirrel bool     `json:"squirrel"`
}

type Counts struct {
	n00, n01, n10, n11 uint
}

func main() {
	data, err := os.ReadFile("sample.json")
	if err != nil {
		fmt.Println("Could not open the file:", err)
		return
	}

	var journalEntry []Entry
	err = json.Unmarshal(data, &journalEntry)
	if err != nil {
		fmt.Println("File could not unmarshal:", err)
		return
	}

	//fmt.Println("Journal data loaded", journalEntry)

	corrValues := make(map[string]float64)

	for _, entry := range journalEntry {
		for _, e := range entry.Events {
			if _, exists := corrValues[e]; !exists {
				counts := getCounts(journalEntry, e)
				corr := phi(counts)
				corrValues[e] = corr
			}
		}
	}

	max := -1.0
	min := 1.0
	var maxEvent, minEvent string
	for key, value := range corrValues {
		if value > max {
			max = value
			maxEvent = key
		}
		if value < min {
			min = value
			minEvent = key
		}
	}

	fmt.Printf("Most Positively Correlated event: %s (%f)\n", maxEvent, max)
	fmt.Printf("Most Negatively Correlated event: %s (%f)\n", minEvent, min)
}

func getCounts(journalEntries []Entry, event string) Counts {
	var n00, n01, n10, n11 uint
	for _, entry := range journalEntries {
		if slices.Contains(entry.Events, event) {
			if entry.Squirrel {
				n11++
			} else {
				n10++
			}
		} else {
			if entry.Squirrel {
				n01++
			} else {
				n00++
			}
		}
	}
	return Counts{n00, n01, n10, n11}
}

func phi(count Counts) float64 {
	n00 := float64(count.n00)
	n01 := float64(count.n01)
	n10 := float64(count.n10)
	n11 := float64(count.n11)

	n1_ := n10 + n11
	n0_ := n00 + n01
	n_1 := n01 + n11
	n_0 := n00 + n10

	num := (n11 * n00) - (n10 * n01)
	den := math.Sqrt(n1_ * n0_ * n_1 * n_0)

	if den == 0 {
		return 0
	}
	return num / den
}
