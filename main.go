package main

import (
	"net/http"

	"github.com/seetharamu/Banking/views"
)

func main() {
	http.HandleFunc("/transactions", views.AddTransaction)
	http.HandleFunc("/statistics", views.GetStatistics)
	http.HandleFunc("/location", views.SetLocation)
	http.HandleFunc("/location/reset", views.ResetLocation)
	http.ListenAndServe(":8000", nil)
}
