// code/pages/invoice/invoice-save.go
package invoice

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"GOapp_invoice/code/database"
)

// Store invoice record in the database
func SaveInvoiceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var invoiceData database.InvoiceData
		err := json.NewDecoder(r.Body).Decode(&invoiceData)
		if err != nil {
			http.Error(w, "Failed to decode invoice data", http.StatusBadRequest)
			return
		}

		// Insert the invoice data into the invoices table
		_, err = db.Exec(`
            INSERT INTO invoices (invoiceNumber, invoiceDate, clientName, parentName, address1, address2, phone, email, cost, VAT, total)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `, invoiceData.InvoiceNumber, invoiceData.InvoiceDate, invoiceData.ClientName, invoiceData.ParentName, invoiceData.Address1, invoiceData.Address2, invoiceData.Phone, invoiceData.Email, invoiceData.Cost, invoiceData.VAT, invoiceData.Total)
		if err != nil {
			http.Error(w, "Failed to save invoice to database", http.StatusInternalServerError)
			return
		}

		// Insert job-row data into the invoices_job_row table
		for _, job := range invoiceData.Jobs {
			_, err = db.Exec(`
                INSERT INTO invoices_job_row (invoiceNumber, jobName, quantity, price, fullPrice)
                VALUES (?, ?, ?, ?, ?)
            `, invoiceData.InvoiceNumber, job.JobName, job.Quantity, job.Price, job.FullPrice)
			if err != nil {
				http.Error(w, "Failed to save job-row data to database", http.StatusInternalServerError)
				return
			}
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}
