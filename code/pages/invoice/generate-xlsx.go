// code/pages/invoice/generate-xlsx.go
package invoice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"GOapp_invoice/code/database"

	"github.com/xuri/excelize/v2"
)

func GenerateXLSX(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON request body
	var invoiceData database.InvoiceData
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
}
