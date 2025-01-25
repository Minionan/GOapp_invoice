// code/pages/invoice/invoice.go
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

// List all jobs in the database
func GetJobsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query the database for job entries
		rows, err := db.Query("SELECT id, jobName, price FROM jobs") // Include 'id' in the query
		if err != nil {
			http.Error(w, "Failed to query jobs data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Iterate over each row
		var jobs []database.Job
		for rows.Next() {
			var job database.Job
			// Scan the columns (include 'id' here)
			if err := rows.Scan(&job.ID, &job.JobName, &job.Price); err != nil {
				http.Error(w, "Failed to scan job data", http.StatusInternalServerError)
				return
			}
			// Accumulate each job into the jobs slice.
			jobs = append(jobs, job)
		}

		// Send the list of jobs in JSON format.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jobs)
	}
}

// List client details from the database
func GetClientsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query the database for client information.
		rows, err := db.Query("SELECT clientName, parentName, address1, address2, phone, email, abbreviation FROM clients")
		if err != nil {
			http.Error(w, "Failed to query clients data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var clients []database.Client

		// Iterate over the results
		for rows.Next() {
			var client database.Client
			// Scan each column into the respective Client struct fields.
			if err := rows.Scan(
				&client.ClientName,
				&client.ParentName,
				&client.Address1,
				&client.Address2,
				&client.Phone,
				&client.Email,
				&client.Abbreviation,
			); err != nil {
				http.Error(w, "Failed to scan client data", http.StatusInternalServerError)
				return
			}
			// Accumulate client info into the clients slice.
			clients = append(clients, client)
		}

		// Send the client list in JSON format.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	}
}

// List invoices stored in the database
func GetInvoicesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query the database for invoice entries
		rows, err := db.Query("SELECT invoiceNumber, clientName, parentName, phone, email, cost, total FROM invoices")
		if err != nil {
			http.Error(w, "Failed to query invoices data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Iterate over each row
		var invoices []database.InvoiceData
		for rows.Next() {
			var invoice database.InvoiceData
			// Scan the columns
			if err := rows.Scan(&invoice.InvoiceNumber, &invoice.ClientName, &invoice.ParentName, &invoice.Phone, &invoice.Email, &invoice.Cost, &invoice.Total); err != nil {
				http.Error(w, "Failed to scan invoice data", http.StatusInternalServerError)
				return
			}
			// Accumulate each invoice into the invoices slice
			invoices = append(invoices, invoice)
		}

		// Send the list of invoices in JSON format
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invoices)
	}
}
