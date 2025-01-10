// main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Job struct {
	JobName string `json:"jobName"`
	Price   string `json:"price"`
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

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./db/data.db")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
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

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
