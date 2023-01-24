package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/seetharamu/Banking/models"
)

const rollingTimeWindow = 60 // seconds

var transactions = make([]models.Transaction, 0, rollingTimeWindow)
var location string
var mu sync.Mutex

func AddTransaction(w http.ResponseWriter, r *http.Request, t models.Transaction) {
	// Check timestamp
	if t.Timestamp.After(time.Now()) {
		http.Error(w, "Invalid Timestamp", http.StatusUnprocessableEntity)
		return
	}
	fmt.Println(t.Timestamp)
	fmt.Println(t.Amount)
	fmt.Println(t.Location)

	// Check location
	if location != "" && location != t.Location {
		http.Error(w, "Unauthorized Location", http.StatusUnauthorized)
		return
	}

	// Check if transaction is older than 60 seconds
	if time.Since(t.Timestamp) > time.Second*60 {
		http.Error(w, "Transaction is older than 60 seconds", http.StatusNoContent)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Add transaction to the rolling time window
	if len(transactions) == cap(transactions) {
		transactions = append(transactions[1:], t)
	} else {
		transactions = append(transactions, t)
	}
	fmt.Println(transactions)
	w.WriteHeader(http.StatusCreated)
}

func GetStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Filter transactions that happened in the last 60 seconds
	now := time.Now()
	recentTransactions := make([]models.Transaction, 0)
	for _, t := range transactions {
		if now.Sub(t.Timestamp) <= time.Second*rollingTimeWindow {
			recentTransactions = append(recentTransactions, t)
		}
	}

	// Calculate statistics
	var stats models.Statistics
	for _, t := range recentTransactions {
		stats.Sum += t.Amount
		stats.Count++
		if t.Amount > stats.Max {
			stats.Max = t.Amount
		}
		if t.Amount < stats.Min {
			stats.Min = t.Amount
		}
	}
	if stats.Count > 0 {
		stats.Avg = stats.Sum / float64(stats.Count)
	}
	// Return statistics in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func SetLocation(w http.ResponseWriter, r *http.Request, l string) {
	mu.Lock()
	location = l
	mu.Unlock()
}

func ResetLocation(w http.ResponseWriter) {
	mu.Lock()
	location = ""
	mu.Unlock()
}
