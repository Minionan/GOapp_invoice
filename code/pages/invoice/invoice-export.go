// code/pages/invoice/invoice-export.go
package invoice

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"

	"GOapp_invoice/code/database"
)

func InvoicesExportCSVHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query all invoices
		invoiceRows, err := db.Query("SELECT invoiceNumber, invoiceDate, clientName, parentName, address1, address2, phone, email, cost, VAT, total FROM invoices")
		if err != nil {
			log.Printf("Failed to query invoices: %v", err)
			http.Error(w, "Failed to query invoices", http.StatusInternalServerError)
			return
		}
		defer invoiceRows.Close()

		// Query all invoice jobs
		jobRows, err := db.Query("SELECT invoiceNumber, jobName, quantity, price, fullPrice FROM invoices_job_row")
		if err != nil {
			log.Printf("Failed to query invoice jobs: %v", err)
			http.Error(w, "Failed to query invoice jobs", http.StatusInternalServerError)
			return
		}
		defer jobRows.Close()

		// Create a CSV file with tab as the delimiter
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=invoices.csv")
		writer := csv.NewWriter(w)
		writer.Comma = '\t' // Set the delimiter to tab
		defer writer.Flush()

		// Write delimiter for the beginning of the invoices table
		writer.Write([]string{"---BEGINNING OF INVOICES TABLE---"})

		// Write invoices header
		writer.Write([]string{"invoiceNumber", "invoiceDate", "clientName", "parentName", "address1", "address2", "phone", "email", "cost", "VAT", "total"})

		// Write invoice rows
		for invoiceRows.Next() {
			var invoice database.InvoiceData
			err := invoiceRows.Scan(
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
				log.Printf("Failed to scan invoice data: %v", err)
				http.Error(w, "Failed to scan invoice data", http.StatusInternalServerError)
				return
			}
			writer.Write([]string{
				invoice.InvoiceNumber,
				invoice.InvoiceDate,
				invoice.ClientName,
				invoice.ParentName,
				invoice.Address1,
				invoice.Address2,
				invoice.Phone,
				invoice.Email,
				fmt.Sprintf("%f", invoice.Cost),
				fmt.Sprintf("%f", invoice.VAT),
				fmt.Sprintf("%f", invoice.Total),
			})
		}

		// Write delimiter for the end of the invoices table
		writer.Write([]string{"---END OF INVOICES TABLE---"})

		// Write delimiter for the beginning of the jobs table
		writer.Write([]string{"---BEGINNING OF JOBS TABLE---"})

		// Write jobs header
		writer.Write([]string{"invoiceNumber", "jobName", "quantity", "price", "fullPrice"})

		// Write job rows
		for jobRows.Next() {
			var job database.Job
			err := jobRows.Scan(
				&job.InvoiceNumber,
				&job.JobName,
				&job.Quantity,
				&job.Price,
				&job.FullPrice,
			)
			if err != nil {
				log.Printf("Failed to scan job data: %v", err)
				http.Error(w, "Failed to scan job data", http.StatusInternalServerError)
				return
			}
			writer.Write([]string{
				job.InvoiceNumber,
				job.JobName,
				job.Quantity,
				job.Price,
				job.FullPrice,
			})
		}

		// Write delimiter for the end of the jobs table
		writer.Write([]string{"---END OF JOBS TABLE---"})
	}
}
