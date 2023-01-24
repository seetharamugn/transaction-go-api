package views

import (
	"encoding/json"
	"net/http"

	"github.com/seetharamu/Banking/controllers"
	"github.com/seetharamu/Banking/models"
)

func AddTransaction(w http.ResponseWriter, r *http.Request) {
	var t models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	controllers.AddTransaction(w, r, t)
}

func GetStatistics(w http.ResponseWriter, r *http.Request) {
	controllers.GetStatistics(w, r)
}

func SetLocation(w http.ResponseWriter, r *http.Request) {
	var l struct {
		Location string `json:"location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	controllers.SetLocation(w, r, l.Location)
}
func ResetLocation(w http.ResponseWriter, r *http.Request) {
	controllers.ResetLocation(w)
}
