// code/pages/invoice/job-edit.go
package invoice

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"GOapp_invoice/code/database" // Import the database package
)

type AddJobRequest struct {
	JobName string `json:"jobName"`
	Price   string `json:"price"`
}

// Fetch job details by ID
func JobDetailsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		var job database.Job // Use database.Job instead of Job
		err := db.QueryRow("SELECT id, jobName, price FROM jobs WHERE id = ?", id).Scan(&job.ID, &job.JobName, &job.Price)
		if err != nil {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(job)
	}
}

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
