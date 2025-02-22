// code/pages/invoice/invoice-import.go
package invoice

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"GOapp_invoice/code/database" // Import the database package
)

func CheckInvoiceExists(db *sql.DB, invoiceNumber string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM invoices WHERE invoiceNumber = ?)", invoiceNumber).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func InvoicesImportCSVHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			log.Printf("Failed to read file: %v", err)
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Failed to read file content: %v", err)
			http.Error(w, "Failed to read file content", http.StatusInternalServerError)
			return
		}

		log.Printf("Raw file content:\n%s", string(content))

		lines := strings.Split(string(content), "\n")

		var importedInvoices []string
		var duplicateInvoices []string
		var malformedRows []int
		var currentTable string
		var isHeaderRow bool

		// Regular expression to split on 2 or more spaces or tabs
		splitRegex := regexp.MustCompile(`\s{2,}|\t+`)

		// Function to handle empty address2 field
		processInvoiceLine := func(line string) []string {
			fields := splitRegex.Split(line, -1)
			result := make([]string, 11)

			if len(fields) == 10 {
				// Copy everything up to address1 (0-4)
				copy(result[0:5], fields[0:5])
				// Insert empty string for address2
				result[5] = ""
				// Copy remaining fields (5-9 to 6-10)
				copy(result[6:], fields[5:])
			} else {
				copy(result, fields)
			}

			return result
		}

		// First pass: Process invoices table
		for lineNumber, line := range lines {
			log.Printf("Processing invoice line %d: %s", lineNumber+1, line)
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "---") {
				switch line {
				case "---BEGINNING OF INVOICES TABLE---":
					currentTable = "invoices"
					isHeaderRow = true
				case "---END OF INVOICES TABLE---":
					currentTable = ""
				case "---BEGINNING OF JOBS TABLE---":
					// Stop processing after invoices table
					break
				}
				continue
			}

			if isHeaderRow {
				isHeaderRow = false
				continue
			}

			if currentTable == "invoices" {
				fields := processInvoiceLine(line)

				exists, err := CheckInvoiceExists(db, fields[0])
				if err != nil {
					log.Printf("Failed to check if invoice exists: %v", err)
					http.Error(w, "Failed to check if invoice exists", http.StatusInternalServerError)
					return
				}

				if exists {
					log.Printf("Skipping duplicate invoice: %s", fields[0])
					duplicateInvoices = append(duplicateInvoices, fields[0])
					continue
				}

				_, err = db.Exec(`
                    INSERT INTO invoices (invoiceNumber, invoiceDate, clientName, parentName, address1, address2, phone, email, cost, VAT, total)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                `, fields[0], fields[1], fields[2], fields[3], fields[4], fields[5], fields[6], fields[7], fields[8], fields[9], fields[10])

				if err != nil {
					log.Printf("Failed to insert invoice data: %v", err)
					http.Error(w, "Failed to insert invoice data", http.StatusInternalServerError)
					return
				}

				importedInvoices = append(importedInvoices, fields[0])
				log.Printf("Successfully imported invoice: %s", fields[0])
			}
		}

		// Second pass: Process jobs table
		currentTable = ""
		isHeaderRow = false

		for lineNumber, line := range lines {
			log.Printf("Processing job line %d: %s", lineNumber+1, line)
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "---") {
				switch line {
				case "---BEGINNING OF JOBS TABLE---":
					currentTable = "jobs"
					isHeaderRow = true
				case "---END OF JOBS TABLE---":
					currentTable = ""
				}
				continue
			}

			if isHeaderRow {
				isHeaderRow = false
				continue
			}

			if currentTable == "jobs" {
				fields := splitRegex.Split(line, -1)
				if len(fields) < 5 {
					log.Printf("Malformed job row at line %d: expected at least 5 fields", lineNumber+1)
					malformedRows = append(malformedRows, lineNumber+1)
					continue
				}

				invoiceNumber := fields[0]

				// Check if this invoice is in our white list
				isWhitelisted := false
				for _, importedInvoice := range importedInvoices {
					if importedInvoice == invoiceNumber {
						isWhitelisted = true
						break
					}
				}

				if !isWhitelisted {
					log.Printf("Skipping job for non-whitelisted invoice: %s", invoiceNumber)
					continue
				}

				// For jobs with names containing spaces, join all fields between invoice number and quantity
				jobName := strings.Join(fields[1:len(fields)-3], " ")
				quantity := fields[len(fields)-3]
				price := fields[len(fields)-2]
				fullPrice := fields[len(fields)-1]

				_, err := db.Exec(`
                    INSERT INTO invoices_job_row (invoiceNumber, jobName, quantity, price, fullPrice)
                    VALUES (?, ?, ?, ?, ?)
                `, invoiceNumber, jobName, quantity, price, fullPrice)

				if err != nil {
					log.Printf("Failed to insert job data: %v", err)
					http.Error(w, "Failed to insert job data", http.StatusInternalServerError)
					return
				}

				log.Printf("Successfully imported job for invoice %s: %s", invoiceNumber, jobName)
			}
		}

		log.Printf("Import process completed. Imported: %d, Duplicates: %d, Malformed: %d",
			len(importedInvoices), len(duplicateInvoices), len(malformedRows))

		response := map[string]interface{}{
			"message":           "Import process completed",
			"importedInvoices":  importedInvoices,
			"duplicateInvoices": duplicateInvoices,
			"malformedRows":     malformedRows,
		}

		if len(duplicateInvoices) > 0 || len(malformedRows) > 0 {
			response["message"] = "Some items were not imported due to duplicates or malformed rows."
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func InvoicesImportTXTHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			log.Printf("Failed to read file: %v", err)
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Failed to read file content: %v", err)
			http.Error(w, "Failed to read file content", http.StatusInternalServerError)
			return
		}

		log.Printf("Raw file content:\n%s", string(content))

		lines := strings.Split(string(content), "\n")

		var invoiceNumber string
		var invoiceDate string
		var clientName string
		var parentName string
		var address1 string
		var address2 string
		var phone string
		var email string
		var cost float64
		var vat float64
		var total float64
		var jobs []database.Job

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "Invoice Number:") {
				invoiceNumber = strings.TrimSpace(strings.TrimPrefix(line, "Invoice Number:"))
				exists, err := CheckInvoiceExists(db, invoiceNumber)
				if err != nil {
					log.Printf("Failed to check if invoice exists: %v", err)
					http.Error(w, "Failed to check if invoice exists", http.StatusInternalServerError)
					return
				}
				if exists {
					log.Printf("Duplicate invoice: %s", invoiceNumber)
					http.Error(w, "Duplicate invoice", http.StatusConflict)
					return
				}
			} else if strings.HasPrefix(line, "Invoice Date:") {
				invoiceDate = strings.TrimSpace(strings.TrimPrefix(line, "Invoice Date:"))
			} else if strings.HasPrefix(line, "Client Name:") {
				clientName = strings.TrimSpace(strings.TrimPrefix(line, "Client Name:"))
			} else if strings.HasPrefix(line, "Parent Name:") {
				parentName = strings.TrimSpace(strings.TrimPrefix(line, "Parent Name:"))
			} else if strings.HasPrefix(line, "Address 1:") {
				address1 = strings.TrimSpace(strings.TrimPrefix(line, "Address 1:"))
			} else if strings.HasPrefix(line, "Address 2:") {
				address2 = strings.TrimSpace(strings.TrimPrefix(line, "Address 2:"))
			} else if strings.HasPrefix(line, "Phone:") {
				phone = strings.TrimSpace(strings.TrimPrefix(line, "Phone:"))
			} else if strings.HasPrefix(line, "Email:") {
				email = strings.TrimSpace(strings.TrimPrefix(line, "Email:"))
			} else if strings.HasPrefix(line, "Job Cost:") {
				cost, _ = strconv.ParseFloat(strings.TrimSpace(strings.TrimPrefix(line, "Job Cost:")), 64)
			} else if strings.HasPrefix(line, "VAT (5%):") {
				vat, _ = strconv.ParseFloat(strings.TrimSpace(strings.TrimPrefix(line, "VAT (5%):")), 64)
			} else if strings.HasPrefix(line, "Total Amount:") {
				total, _ = strconv.ParseFloat(strings.TrimSpace(strings.TrimPrefix(line, "Total Amount:")), 64)
			} else if strings.HasPrefix(line, "Jobs:") {
				// Process job lines
				for _, jobLine := range lines {
					if strings.Contains(jobLine, ":") && strings.Contains(jobLine, "Quantity:") {
						parts := strings.Split(jobLine, ",")
						jobName := strings.TrimSpace(strings.Split(parts[0], ":")[1])
						quantity := strings.TrimSpace(strings.Split(parts[1], ":")[1])
						price := strings.TrimSpace(strings.Split(parts[2], ":")[1])
						fullPrice := strings.TrimSpace(strings.Split(parts[3], ":")[1])

						job := database.Job{
							InvoiceNumber: invoiceNumber,
							JobName:       jobName,
							Quantity:      quantity,
							Price:         price,
							FullPrice:     fullPrice,
						}
						jobs = append(jobs, job)
					}
				}
			}
		}

		// Insert invoice data
		_, err = db.Exec(`
            INSERT INTO invoices (invoiceNumber, invoiceDate, clientName, parentName, address1, address2, phone, email, cost, VAT, total)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `, invoiceNumber, invoiceDate, clientName, parentName, address1, address2, phone, email, cost, vat, total)

		if err != nil {
			log.Printf("Failed to insert invoice data: %v", err)
			http.Error(w, "Failed to insert invoice data", http.StatusInternalServerError)
			return
		}

		// Insert job data
		for _, job := range jobs {
			_, err := db.Exec(`
                INSERT INTO invoices_job_row (invoiceNumber, jobName, quantity, price, fullPrice)
                VALUES (?, ?, ?, ?, ?)
            `, job.InvoiceNumber, job.JobName, job.Quantity, job.Price, job.FullPrice)

			if err != nil {
				log.Printf("Failed to insert job data: %v", err)
				http.Error(w, "Failed to insert job data", http.StatusInternalServerError)
				return
			}
		}

		log.Printf("Successfully imported invoice: %s", invoiceNumber)

		response := map[string]interface{}{
			"message":         "Invoice imported successfully",
			"importedInvoice": invoiceNumber,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
