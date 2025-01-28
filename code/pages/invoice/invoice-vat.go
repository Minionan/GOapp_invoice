// code/pages/invoice/invoiceVat.go
package invoice

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"GOapp_invoice/code/database"
)

// Vat represents the VAT rate structure
type Vat struct {
	Rate float64 `json:"rate"`
}

// GetVatHandler retrieves the current VAT rate from the database
func GetVatHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	var vat Vat
	err := db.QueryRow("SELECT rate FROM vat WHERE id = 1").Scan(&vat.Rate)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "VAT rate not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vat)
}

// UpdateVatHandler updates the VAT rate in the database
func UpdateVatHandler(w http.ResponseWriter, r *http.Request) {
	var vat Vat
	err := json.NewDecoder(r.Body).Decode(&vat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload"}`))
		return
	}

	if vat.Rate < 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "VAT rate cannot be negative"}`))
		return
	}

	db := database.GetDB()
	_, err = db.Exec("UPDATE vat SET rate = ? WHERE id = 1", vat.Rate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to update VAT rate"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "success"}`))
}
