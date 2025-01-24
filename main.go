// main.go
package main

import (
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"

	"GOapp_invoice/code/database"
	"GOapp_invoice/code/pages/invoice"
)

func main() {
	// Initialize the database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Serve static files
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("pages"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve the invoice page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("pages/invoice.html"))
		tmpl.Execute(w, nil)
	})

	// Serve the invoiceForm page
	http.HandleFunc("/invoice-form", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("pages/invoiceForm.html"))
		tmpl.Execute(w, nil)
	})

	// Serve the invoiceJob page
	http.HandleFunc("/invoice-job", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("pages/invoiceJob.html"))
		tmpl.Execute(w, nil)
	})

	// Serve the invoiceJobEdit page
	http.HandleFunc("/pages/invoice-job-edit", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("pages/invoiceJobEdit.html"))
		tmpl.Execute(w, nil)
	})

	// Serve the invoiceJobAdd page
	http.HandleFunc("/pages/invoice-job-add", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("pages/invoiceJobAdd.html"))
		tmpl.Execute(w, nil)
	})

	// Register handlers from the invoice package
	http.HandleFunc("/max-job-rows", invoice.GetMaxJobRowsHandler)
	http.HandleFunc("/jobs", invoice.GetJobsHandler(db))
	http.HandleFunc("/job-details", invoice.JobDetailsHandler(db))
	http.HandleFunc("/job-update", invoice.JobUpdateHandler(db))
	http.HandleFunc("/job-add", invoice.JobAddHandler(db))
	http.HandleFunc("/job-delete", invoice.JobDeleteHandler(db))
	http.HandleFunc("/clients", invoice.GetClientsHandler(db))
	http.HandleFunc("/generate-xlsx", invoice.GenerateXLSX)
	http.HandleFunc("/generate-pdf", invoice.GeneratePDF)

	// Start the server
	log.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
