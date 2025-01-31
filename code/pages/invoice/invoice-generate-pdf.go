// code/pages/invoice/invoice-generate-pdf.go
package invoice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"GOapp_invoice/code/database"

	"github.com/xuri/excelize/v2"
)

func InvoiceGeneratePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var invoiceData database.InvoiceData
	if err := json.NewDecoder(r.Body).Decode(&invoiceData); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// First, generate the XLSX file
	f, err := excelize.OpenFile("template.xlsx")
	if err != nil {
		http.Error(w, "Failed to open template file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Apply replacements (same as in /generate-xlsx)
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
		rows, err := f.GetRows(sheet)
		if err != nil {
			continue
		}

		for rowIdx, row := range rows {
			for colIdx, cell := range row {
				if replacement, exists := replacements[cell]; exists {
					col, err := excelize.ColumnNumberToName(colIdx + 1)
					if err != nil {
						continue
					}
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

	// Save the modified XLSX file
	tempXlsxPath := filepath.Join(os.TempDir(), invoiceData.InvoiceNumber+".xlsx")
	if err := f.SaveAs(tempXlsxPath); err != nil {
		http.Error(w, "Failed to save temporary XLSX file", http.StatusInternalServerError)
		return
	}

	// Convert XLSX to PDF using LibreOffice
	tempPdfPath := filepath.Join(os.TempDir(), invoiceData.InvoiceNumber+".pdf")
	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "pdf", "--outdir", os.TempDir(), tempXlsxPath)
	if err := cmd.Run(); err != nil {
		http.Error(w, "Failed to convert XLSX to PDF", http.StatusInternalServerError)
		return
	}

	// Serve the PDF file
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", invoiceData.InvoiceNumber))
	http.ServeFile(w, r, tempPdfPath)

	// Clean up temporary files
	os.Remove(tempXlsxPath)
	os.Remove(tempPdfPath)
}
