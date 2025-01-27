// code/pages/invoice/client-edit.go
package invoice

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type AddClientRequest struct {
	ClientName   string `json:"clientName"`
	ParentName   string `json:"parentName"`
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Abbreviation string `json:"abbreviation"`
}

// Fetching client details is handled by GetClientsHandler function in code/pages/invoice/invoice.go

// Add new client
func ClientAddHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddClientRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO clients (clientName, parentName, address1, address2, phone, email, abbreviation) VALUES (?, ?, ?, ?, ?, ?, ?)",
			req.ClientName, req.ParentName, req.Address1, req.Address2, req.Phone, req.Email, req.Abbreviation)
		if err != nil {
			http.Error(w, "Failed to insert client", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// Update client details
func ClientUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		var req AddClientRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("UPDATE clients SET clientName = ?, parentName = ?, address1 = ?, address2 = ?, phone = ?, email = ?, abbreviation = ? WHERE id = ?",
			req.ClientName, req.ParentName, req.Address1, req.Address2, req.Phone, req.Email, req.Abbreviation, id)
		if err != nil {
			http.Error(w, "Failed to update client", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// Delete client
func ClientDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		_, err := db.Exec("DELETE FROM clients WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Failed to delete client", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}
