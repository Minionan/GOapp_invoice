// main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
)

type Job struct {
	JobName   string `json:"jobName"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	FullPrice string `json:"fullPrice"`
}

type Client struct {
	ClientName   string `json:"clientName"`
	ParentName   string `json:"parentName"`
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Abbreviation string `json:"abbreviation"`
}

type InvoiceData struct {
	ParentName    string  `json:"parentName"`
	Address1      string  `json:"address1"`
	Address2      string  `json:"address2"`
	Phone         string  `json:"phone"`
	Email         string  `json:"email"`
	InvoiceNumber string  `json:"invoiceNumber"`
	InvoiceDate   string  `json:"invoiceDate"`
	Cost          float64 `json:"cost"`
	VAT           float64 `json:"vat"`
	Total         float64 `json:"total"`
	Jobs          []Job   `json:"jobs"`
}

var db *sql.DB

func getMaxJobRows(f *excelize.File) (int, error) {
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

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./db/data.db")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/max-job-rows", func(w http.ResponseWriter, r *http.Request) {
		f, err := excelize.OpenFile("template.xlsx")
		if err != nil {
			http.Error(w, "Failed to open template file", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		maxRows, err := getMaxJobRows(f)
		if err != nil {
			http.Error(w, "Failed to get max rows", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"maxRows": maxRows})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("static/invoice.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT jobName, price FROM jobs")
		if err != nil {
			http.Error(w, "Failed to query jobs data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var jobs []Job
		for rows.Next() {
			var job Job
			if err := rows.Scan(&job.JobName, &job.Price); err != nil {
				http.Error(w, "Failed to scan job data", http.StatusInternalServerError)
				return
			}
			jobs = append(jobs, job)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jobs)
	})

	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT clientName, parentName, address1, address2, phone, email, abbreviation FROM clients")
		if err != nil {
			http.Error(w, "Failed to query clients data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var clients []Client
		for rows.Next() {
			var client Client
			if err := rows.Scan(&client.ClientName, &client.ParentName, &client.Address1, &client.Address2, &client.Phone, &client.Email, &client.Abbreviation); err != nil {
				http.Error(w, "Failed to scan client data", http.StatusInternalServerError)
				return
			}
			clients = append(clients, client)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/generate-xlsx", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode the JSON request body
		var invoiceData InvoiceData
		if err := json.NewDecoder(r.Body).Decode(&invoiceData); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		// Open the template file
		f, err := excelize.OpenFile("template.xlsx")
		if err != nil {
			http.Error(w, "Failed to open template file", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// Basic replacements (existing code)
		replacements := map[string]string{
			"{{PARENT_NAME}}":    invoiceData.ParentName,
			"{{ADDRESS1}}":       invoiceData.Address1,
			"{{ADDRESS2}}":       invoiceData.Address2,
			"{{PHONE}}":          invoiceData.Phone,
			"{{EMAIL}}":          invoiceData.Email,
			"{{INVOICE_NUMBER}}": invoiceData.InvoiceNumber,
			"{{INVOICE_DATE}}":   invoiceData.InvoiceDate,
			"{{COST}}":           fmt.Sprintf("%.2f", invoiceData.Cost),
			"{{VAT}}":            fmt.Sprintf("%.2f", invoiceData.VAT),
			"{{TOTAL}}":          fmt.Sprintf("%.2f", invoiceData.Total),
		}

		// Add job-related replacements
		for i, job := range invoiceData.Jobs {
			jobNum := i + 1
			replacements[fmt.Sprintf("{{JOB_NAME_%d}}", jobNum)] = job.JobName
			replacements[fmt.Sprintf("{{QUANTITY_%d}}", jobNum)] = job.Quantity
			replacements[fmt.Sprintf("{{PRICE_%d}}", jobNum)] = job.Price
			replacements[fmt.Sprintf("{{FULL_PRICE_%d}}", jobNum)] = job.FullPrice
		}

		// Get all sheet names
		sheets := f.GetSheetList()

		// Replace placeholders in all sheets
		for _, sheet := range sheets {
			// Get all cells in the sheet
			rows, err := f.GetRows(sheet)
			if err != nil {
				continue
			}

			for rowIdx, row := range rows {
				for colIdx, cell := range row {
					// Check if the cell contains any of our placeholders
					if replacement, exists := replacements[cell]; exists {
						// Convert column index to Excel column letter
						col, err := excelize.ColumnNumberToName(colIdx + 1)
						if err != nil {
							continue
						}
						// Replace the placeholder
						f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, rowIdx+1), replacement)
					}
				}
			}
		}

		// Clear any remaining placeholders
		for _, sheet := range sheets {
			rows, err := f.GetRows(sheet)
			if err != nil {
				continue
			}

			for rowIdx, row := range rows {
				for colIdx, cell := range row {
					// Check if the cell contains any placeholder pattern (starts with {{ and ends with }})
					if strings.HasPrefix(cell, "{{") && strings.HasSuffix(cell, "}}") {
						col, err := excelize.ColumnNumberToName(colIdx + 1)
						if err != nil {
							continue
						}
						// Clear any remaining placeholder by setting it to empty string
						f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, rowIdx+1), "")
					}
				}
			}
		}

		// Set the appropriate headers for file download
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", invoiceData.InvoiceNumber))

		// Write the modified file directly to the response writer
		if err := f.Write(w); err != nil {
			http.Error(w, "Failed to write XLSX file", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
