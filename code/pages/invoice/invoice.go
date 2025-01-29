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

func CheckInvoiceNumberExists(db *sql.DB, invoiceNumber string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM invoices WHERE invoiceNumber = ?)", invoiceNumber).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Check if the chosen invoiceNumber is already stored in the database
func CheckInvoiceNumberExistsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invoiceNumber := r.URL.Query().Get("invoiceNumber")
		if invoiceNumber == "" {
			http.Error(w, "Invoice number is required", http.StatusBadRequest)
			return
		}

		exists, err := CheckInvoiceNumberExists(db, invoiceNumber)
		if err != nil {
			http.Error(w, "Failed to check invoice number", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"exists": exists})
	}
}

// List all jobs in the database
func JobGetHandler(db *sql.DB) http.HandlerFunc {
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
		rows, err := db.Query("SELECT id, clientName, parentName, address1, address2, phone, email, abbreviation FROM clients")
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
				&client.ID,
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

// List invoices stored in the database (used for displaying invoices on invoice.html and invoiceEdit.html pages)
func GetInvoicesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invoiceNumber := r.URL.Query().Get("invoiceNumber")

		if invoiceNumber != "" {
			// Fetch detailed invoice data for a specific invoice
			var invoice database.InvoiceData
			err := db.QueryRow(`
                SELECT invoiceNumber, invoiceDate, clientName, parentName, address1, address2, phone, email, cost, VAT, total
                FROM invoices
                WHERE invoiceNumber = ?
            `, invoiceNumber).Scan(
				&invoice.InvoiceNumber,
				&invoice.InvoiceDate,
				&invoice.ClientName,
				&invoice.ParentName,
				&invoice.Address1,
				&invoice.Address2,
				&invoice.Phone,
				&invoice.Email,
				&invoice.Cost,
				&invoice.VAT,
				&invoice.Total,
			)
			if err != nil {
				http.Error(w, "Failed to fetch invoice details", http.StatusInternalServerError)
				return
			}

			// Fetch the jobs associated with the invoice
			rows, err := db.Query(`
                SELECT jobName, quantity, price, fullPrice
                FROM invoices_job_row
                WHERE invoiceNumber = ?
            `, invoiceNumber)
			if err != nil {
				http.Error(w, "Failed to fetch job details", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var jobs []database.Job
			for rows.Next() {
				var job database.Job
				if err := rows.Scan(&job.JobName, &job.Quantity, &job.Price, &job.FullPrice); err != nil {
					http.Error(w, "Failed to scan job details", http.StatusInternalServerError)
					return
				}
				jobs = append(jobs, job)
			}

			invoice.Jobs = jobs

			// Return the detailed invoice data in JSON format
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(invoice)
		} else {
			// Fetch a list of all invoices (for invoice.html)
			rows, err := db.Query("SELECT invoiceNumber, parentName, email, cost, total, invoiceDate FROM invoices ORDER BY invoiceDate DESC LIMIT 10")
			if err != nil {
				http.Error(w, "Failed to query invoices data", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var invoices []database.InvoiceData
			for rows.Next() {
				var invoice database.InvoiceData
				if err := rows.Scan(&invoice.InvoiceNumber, &invoice.ParentName, &invoice.Email, &invoice.Cost, &invoice.Total, &invoice.InvoiceDate); err != nil {
					http.Error(w, "Failed to scan invoice data", http.StatusInternalServerError)
					return
				}
				invoices = append(invoices, invoice)
			}

			// Return the list of invoices in JSON format
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(invoices)
		}
	}
}
