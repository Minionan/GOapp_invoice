// code/pages/invoice/job-edit.go
package invoice

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"net/http"
)

type AddJobRequest struct {
	JobName string `json:"jobName"`
	Price   string `json:"price"`
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

		_, err := db.Exec("INSERT INTO jobs (jobName, price) VALUES (?, ?)", req.JobName, req.Price)
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

		_, err := db.Exec("UPDATE jobs SET jobName = ?, price = ? WHERE id = ?", req.JobName, req.Price, id)
		if err != nil {
			http.Error(w, "Failed to update job", http.StatusInternalServerError)
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
		rows, err := db.Query("SELECT jobName, price FROM jobs")
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
		csvWriter.Write([]string{"Job Name", "Price"})

		// Write job data
		for rows.Next() {
			var jobName, price string
			if err := rows.Scan(&jobName, &price); err != nil {
				http.Error(w, "Failed to read job data", http.StatusInternalServerError)
				return
			}
			csvWriter.Write([]string{jobName, price})
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating over job data", http.StatusInternalServerError)
			return
		}
	}
}
