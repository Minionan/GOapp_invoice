// code/pages/invoice/job-edit.go
package invoice

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
)

type AddJobRequest struct {
	JobName string `json:"jobName"`
	Price   string `json:"price"`
	Status  bool   `json:"status"`
}

// Fetching job details is handled by JobGetHandler function in code/pages/invoice/invoice.go

// Add new job
func JobAddHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddJobRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO jobs (jobName, price, status) VALUES (?, ?, ?)", req.JobName, req.Price, req.Status)
		if err != nil {
			http.Error(w, "Failed to insert job", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// Update job details
func JobUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		var req AddJobRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("UPDATE jobs SET jobName = ?, price = ?, status = ? WHERE id = ?", req.JobName, req.Price, req.Status, id)
		if err != nil {
			http.Error(w, "Failed to update job", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// Job status update (used only for invoiceJob.js to update job status on checkmark changes)
func JobStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		status := r.URL.Query().Get("status")

		// Convert status to boolean
		statusBool := status == "true"

		_, err := db.Exec("UPDATE jobs SET status = ? WHERE id = ?", statusBool, id)
		if err != nil {
			http.Error(w, "Failed to update job status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// Delete job
func JobDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		_, err := db.Exec("DELETE FROM jobs WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Failed to delete job", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// ExportJobsHandler exports the list of jobs as a CSV file
func JobExportHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT jobName, price, status FROM jobs")
		if err != nil {
			http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=jobs.csv")
		w.WriteHeader(http.StatusOK)

		csvWriter := csv.NewWriter(w)
		defer csvWriter.Flush()

		// Write CSV header
		csvWriter.Write([]string{"Job Name", "Price", "Status"})

		// Write job data
		for rows.Next() {
			var jobName, price string
			var status bool // Added status variable declaration
			if err := rows.Scan(&jobName, &price, &status); err != nil {
				http.Error(w, "Failed to read job data", http.StatusInternalServerError)
				return
			}
			csvWriter.Write([]string{jobName, price, strconv.FormatBool(status)})
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating over job data", http.StatusInternalServerError)
			return
		}
	}
}

// Import jobs from CSV
func JobImportHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			http.Error(w, "Failed to parse CSV", http.StatusInternalServerError)
			return
		}

		importedJobs := []string{} // Ensure this is always an empty array, not nil
		skippedJobs := []string{}  // Ensure this is always an empty array, not nil

		for i, record := range records {
			if i == 0 {
				continue // Skip header
			}

			jobName := record[0]
			price := record[1]
			statusStr := record[2]

			// Convert status string to bool
			statusBool, err := strconv.ParseBool(statusStr)
			if err != nil {
				http.Error(w, "Invalid status value in CSV", http.StatusBadRequest)
				return
			}

			// Check if job already exists
			var exists bool
			err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM jobs WHERE jobName = ?)", jobName).Scan(&exists)
			if err != nil {
				http.Error(w, "Failed to check job existence", http.StatusInternalServerError)
				return
			}

			if exists {
				skippedJobs = append(skippedJobs, jobName)
				continue
			}

			// Insert new job
			_, err = db.Exec("INSERT INTO jobs (jobName, price, status) VALUES (?, ?, ?)", jobName, price, statusBool)
			if err != nil {
				http.Error(w, "Failed to insert job", http.StatusInternalServerError)
				return
			}

			importedJobs = append(importedJobs, jobName)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"imported": importedJobs,
			"skipped":  skippedJobs,
		})
	}
}
