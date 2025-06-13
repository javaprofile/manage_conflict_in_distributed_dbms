package main

import (
	"fmt"
)

type Transaction struct {
	ID      int
	Timestamp int
	Reads   map[string]bool
	Writes  map[string]bool
}

func ConflictRate(transactions []Transaction) float64 {
	var conflicts int
	var totalTransactions int = len(transactions)

	for i := 0; i < totalTransactions; i++ {
		for j := i + 1; j < totalTransactions; j++ {
			if isConflict(transactions[i], transactions[j]) {
				conflicts++
			}
		}
	}

	conflictRate := float64(conflicts) / float64(totalTransactions*(totalTransactions-1)/2) * 100
	return conflictRate
}

func isConflict(t1, t2 Transaction) bool {
	for item := range t1.Writes {
		if _, exists := t2.Writes[item]; exists {
			return true
		}
	}
	for item := range t1.Writes {
		if _, exists := t2.Reads[item]; exists {
			return true
		}
	}
	for item := range t2.Writes {
		if _, exists := t1.Reads[item]; exists {
			return true
		}
	}
	return false
}

func main() {
	transactions := []Transaction{
		{ID: 1, Timestamp: 1, Reads: map[string]bool{"A": true}, Writes: map[string]bool{"B": true}},
		{ID: 2, Timestamp: 2, Reads: map[string]bool{"B": true}, Writes: map[string]bool{"C": true}},
		{ID: 3, Timestamp: 3, Reads: map[string]bool{"A": true}, Writes: map[string]bool{"C": true}},
		{ID: 4, Timestamp: 4, Reads: map[string]bool{"C": true}, Writes: map[string]bool{"D": true}},
	}

	conflictRate := ConflictRate(transactions)
	fmt.Printf("Conflict Rate: %.2f%%\n", conflictRate)
}
