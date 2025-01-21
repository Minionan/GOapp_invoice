// code/invoice/invoice.go
package invoice

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"GOapp_invoice/code/database"

	"github.com/xuri/excelize/v2"
)

func GetMaxJobRows(f *excelize.File) (int, error) {
	sheets := f.GetSheetList()
	//fmt.Printf("Searching through sheets: %v\n", sheets) // Logging

	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			fmt.Printf("Error reading sheet %s: %v\n", sheet, err)
			continue
		}

		for _, row := range rows {
			for _, cell := range row {
				// Check if cell starts with "{{MAX_JOB_ROWS_"
				if strings.HasPrefix(cell, "{{MAX_JOB_ROWS_") && strings.HasSuffix(cell, "}}") {
					// Extract the number between "{{MAX_JOB_ROWS_" and "}}"
					numStr := strings.TrimPrefix(cell, "{{MAX_JOB_ROWS_")
					numStr = strings.TrimSuffix(numStr, "}}")

					// Convert the extracted string to an integer
					maxRows, err := strconv.Atoi(numStr)
					if err != nil {
						//fmt.Printf("Error converting max rows value: %v\n", err) // Logging
						return 3, nil // Default value if conversion fails
					}

					//fmt.Printf("Found MAX_JOB_ROWS value: %d\n", maxRows) // Logging
					return maxRows, nil
				}
			}
		}
	}
	fmt.Println("Placeholder not found, returning default value")
	return 1, nil // Default value if placeholder not found
}

// Get maxJobRow from template.xlsx file
func GetMaxJobRowsHandler(w http.ResponseWriter, r *http.Request) {
	// Attempt to open the "template.xlsx" file.
	f, err := excelize.OpenFile("template.xlsx")
	if err != nil {
		http.Error(w, "Failed to open template file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	maxRows, err := GetMaxJobRows(f)
	if err != nil {
		http.Error(w, "Failed to get max rows", http.StatusInternalServerError)
		return
	}
	// Return maxJobRows in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"maxRows": maxRows})
}

func GetJobsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT jobName, price FROM jobs")
		if err != nil {
			http.Error(w, "Failed to query jobs data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var jobs []database.Job
		for rows.Next() {
			var job database.Job
			if err := rows.Scan(&job.JobName, &job.Price); err != nil {
				http.Error(w, "Failed to scan job data", http.StatusInternalServerError)
				return
			}
			jobs = append(jobs, job)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jobs)
	}
}

func GetClientsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT clientName, parentName, address1, address2, phone, email, abbreviation FROM clients")
		if err != nil {
			http.Error(w, "Failed to query clients data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var clients []database.Client
		for rows.Next() {
			var client database.Client
			if err := rows.Scan(&client.ClientName, &client.ParentName, &client.Address1, &client.Address2, &client.Phone, &client.Email, &client.Abbreviation); err != nil {
				http.Error(w, "Failed to scan client data", http.StatusInternalServerError)
				return
			}
			clients = append(clients, client)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	}
}
