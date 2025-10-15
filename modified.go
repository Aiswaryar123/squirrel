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
	n00 uint
	n01 uint
	n10 uint
	n11 uint
}
type MaxMin struct {
	max      float64
	min      float64
	maxEvent string
	minEvent string
}

func preprocess(journalEntries []Entry) []Entry {
	var journal []Entry
	for _, entry := range journalEntries {
		hasPeanuts := slices.Contains(entry.Events, "peanuts")
		notBrushedTeeth := !slices.Contains(entry.Events, "brushed teeth")
		if hasPeanuts && notBrushedTeeth {
			entry.Events = append(entry.Events, "dirty teeth")
			journal = append(journal, entry)
		} else {
			journal = append(journal, entry)
		}
	}
	return journal
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

	return Counts{n00: n00, n01: n01, n10: n10, n11: n11}
}
func phi(counts Counts) float64 {
	n00 := float64(counts.n00)
	n01 := float64(counts.n01)
	n10 := float64(counts.n10)
	n11 := float64(counts.n11)
	n1_ := n10 + n11
	n0_ := n00 + n01
	n_1 := n01 + n11
	n_0 := n10 + n00
	num := (n11 * n00) - (n10 * n01)
	den := math.Sqrt(n1_ * n0_ * n_1 * n_0)
	return num / den
}
func getCorrelations(journalEntries []Entry) map[string]float64 {
	corrValues := make(map[string]float64)
	for _, entry := range journalEntries {
		for _, e := range entry.Events {
			counts := getCounts(journalEntries, e)
			corr := phi(counts)
			corrValues[e] = corr

		}
	}
	return corrValues
}
func getMaxMin(corrValues map[string]float64) MaxMin {
	var results MaxMin
	results.max = -1.0
	results.min = 1.0
	for key, value := range corrValues {
		if value > results.max {
			results.max = value
			results.maxEvent = key
		}
		if value < results.min {
			results.min = value
			results.minEvent = key
		}

	}
	return results
}

func main() {
	data, err := os.ReadFile("sample.json")
	if err != nil {
		fmt.Println("Could not open json file: ", err)
	}

	var journalEntries []Entry
	err = json.Unmarshal(data, &journalEntries)
	if err != nil {
		fmt.Println("Could not unmarshal json file: ", err)
	}
	//fmt.Println(journalEntries)
	journalEntries = preprocess(journalEntries)
	
	corrValues := getCorrelations(journalEntries)
	results := getMaxMin(corrValues)

	fmt.Printf("Most Positively Correlated event: %s (%f)\n", results.maxEvent, results.max)
	fmt.Printf("Most Negatively Correlated event: %s (%f)\n", results.minEvent, results.min)

}