// code/pages/invoice/client-edit.go
package invoice

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
)

type AddClientRequest struct {
	ClientName   string `json:"clientName"`
	ParentName   string `json:"parentName"`
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Abbreviation string `json:"abbreviation"`
	Status       bool   `json:"status"`
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

		_, err := db.Exec("INSERT INTO clients (clientName, parentName, address1, address2, phone, email, abbreviation, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			req.ClientName, req.ParentName, req.Address1, req.Address2, req.Phone, req.Email, req.Abbreviation, req.Status)
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

		_, err := db.Exec("UPDATE clients SET clientName = ?, parentName = ?, address1 = ?, address2 = ?, phone = ?, email = ?, abbreviation = ?, status = ? WHERE id = ?",
			req.ClientName, req.ParentName, req.Address1, req.Address2, req.Phone, req.Email, req.Abbreviation, req.Status, id)
		if err != nil {
			http.Error(w, "Failed to update client", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// Clinet status update (used only for client.js to update client status on checkmark changes)
func ClientStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		status := r.URL.Query().Get("status")

		// Convert status to boolean
		statusBool := status == "true"

		_, err := db.Exec("UPDATE clients SET status = ? WHERE id = ?", statusBool, id)
		if err != nil {
			http.Error(w, "Failed to update client status", http.StatusInternalServerError)
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

// Export clients to CSV
func ClientExportCSVHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT clientName, parentName, address1, address2, phone, email, abbreviation, status FROM clients")
		if err != nil {
			http.Error(w, "Failed to fetch clients", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=clients.csv")
		csvWriter := csv.NewWriter(w)
		defer csvWriter.Flush()

		// Write CSV header
		if err := csvWriter.Write([]string{"Client Name", "Payee Name", "Address 1", "Address 2", "Phone", "Email", "Abbreviation", "Status"}); err != nil {
			http.Error(w, "Failed to write CSV header", http.StatusInternalServerError)
			return
		}

		// Write rows
		for rows.Next() {
			var clientName, parentName, address1, address2, phone, email, abbreviation string
			var status bool
			if err := rows.Scan(&clientName, &parentName, &address1, &address2, &phone, &email, &abbreviation, &status); err != nil {
				http.Error(w, "Failed to read client data", http.StatusInternalServerError)
				return
			}
			if err := csvWriter.Write([]string{clientName, parentName, address1, address2, phone, email, abbreviation, strconv.FormatBool(status)}); err != nil {
				http.Error(w, "Failed to write CSV row", http.StatusInternalServerError)
				return
			}
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating over client data", http.StatusInternalServerError)
			return
		}
	}
}

// Import clients from CSV
func ClientImportHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			http.Error(w, "Failed to parse CSV", http.StatusInternalServerError)
			return
		}

		importedClients := []string{} // Ensure this is always an empty array, not nil
		skippedClients := []string{}  // Ensure this is always an empty array, not nil

		for i, record := range records {
			if i == 0 {
				continue // Skip header
			}

			clientName := record[0]
			parentName := record[1]
			address1 := record[2]
			address2 := record[3]
			phone := record[4]
			email := record[5]
			abbreviation := record[6]
			statusStr := record[7]

			// Convert status string to bool
			statusBool, err := strconv.ParseBool(statusStr)
			if err != nil {
				http.Error(w, "Invalid status value in CSV", http.StatusBadRequest)
				return
			}

			// Check if client already exists
			var exists bool
			err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM clients WHERE clientName = ? AND abbreviation = ?)", clientName, abbreviation).Scan(&exists)
			if err != nil {
				http.Error(w, "Failed to check client existence", http.StatusInternalServerError)
				return
			}

			if exists {
				skippedClients = append(skippedClients, clientName)
				continue
			}

			// Insert new client
			_, err = db.Exec("INSERT INTO clients (clientName, parentName, address1, address2, phone, email, abbreviation, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
				clientName, parentName, address1, address2, phone, email, abbreviation, statusBool)
			if err != nil {
				http.Error(w, "Failed to insert client", http.StatusInternalServerError)
				return
			}

			importedClients = append(importedClients, clientName)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"imported": importedClients,
			"skipped":  skippedClients,
		})
	}
}
