package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/seetharamu/Banking/models"
)

func HandleStatistics(w http.ResponseWriter, r *http.Request) {
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
