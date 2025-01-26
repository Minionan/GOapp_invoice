// code/pages/invoice/invoice-edit.go
package invoice

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"GOapp_invoice/code/database"
)

// update existing invoice in the database
func UpdateInvoiceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var invoiceData database.InvoiceData
		err := json.NewDecoder(r.Body).Decode(&invoiceData)
		if err != nil {
			http.Error(w, "Failed to decode invoice data", http.StatusBadRequest)
			return
		}

		// Update the invoice in the database
		_, err = db.Exec(`
            UPDATE invoices
            SET invoiceDate = ?, clientName = ?, parentName = ?, address1 = ?, address2 = ?, phone = ?, email = ?, cost = ?, VAT = ?, total = ?
            WHERE invoiceNumber = ?
        `, invoiceData.InvoiceDate, invoiceData.ClientName, invoiceData.ParentName, invoiceData.Address1, invoiceData.Address2, invoiceData.Phone, invoiceData.Email, invoiceData.Cost, invoiceData.VAT, invoiceData.Total, invoiceData.InvoiceNumber)
		if err != nil {
			http.Error(w, "Failed to update invoice", http.StatusInternalServerError)
			return
		}

		// Delete existing job rows for the invoice
		_, err = db.Exec("DELETE FROM invoices_job_row WHERE invoiceNumber = ?", invoiceData.InvoiceNumber)
		if err != nil {
			http.Error(w, "Failed to delete existing job rows", http.StatusInternalServerError)
			return
		}

		// Insert updated job rows
		for _, job := range invoiceData.Jobs {
			_, err = db.Exec(`
                INSERT INTO invoices_job_row (invoiceNumber, jobName, quantity, price, fullPrice)
                VALUES (?, ?, ?, ?, ?)
            `, invoiceData.InvoiceNumber, job.JobName, job.Quantity, job.Price, job.FullPrice)
			if err != nil {
				http.Error(w, "Failed to insert job row", http.StatusInternalServerError)
				return
			}
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// delete invoice from the database
func DeleteInvoiceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invoiceNumber := r.URL.Query().Get("invoiceNumber")
		if invoiceNumber == "" {
			http.Error(w, "Invoice number is required", http.StatusBadRequest)
			return
		}

		// Delete the invoice from the database
		_, err := db.Exec("DELETE FROM invoices WHERE invoiceNumber = ?", invoiceNumber)
		if err != nil {
			http.Error(w, "Failed to delete invoice", http.StatusInternalServerError)
			return
		}

		// Delete associated job rows
		_, err = db.Exec("DELETE FROM invoices_job_row WHERE invoiceNumber = ?", invoiceNumber)
		if err != nil {
			http.Error(w, "Failed to delete associated job rows", http.StatusInternalServerError)
			return
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}
