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

	// Serve the client page
	http.HandleFunc("/pages/client", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("pages/clients.html"))
		tmpl.Execute(w, nil)
	})

	// Serve the clientAdd page
	http.HandleFunc("/pages/clients-add", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("pages/clientsAdd.html"))
		tmpl.Execute(w, nil)
	})

	// Serve the clientEdit page
	http.HandleFunc("/pages/clients-edit", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("pages/clientsEdit.html"))
		tmpl.Execute(w, nil)
	})

	// Register handlers from the invoice package
	http.HandleFunc("/max-job-rows", invoice.GetMaxJobRowsHandler)
	http.HandleFunc("/jobs", invoice.JobGetHandler(db))
	http.HandleFunc("/job-update", invoice.JobUpdateHandler(db))
	http.HandleFunc("/job-add", invoice.JobAddHandler(db))
	http.HandleFunc("/job-delete", invoice.JobDeleteHandler(db))
	http.HandleFunc("/job-export", invoice.JobExportHandler(db))
	http.HandleFunc("/job-import", invoice.JobImportHandler(db))
	http.HandleFunc("/invoice-save", invoice.SaveInvoiceHandler(db))
	http.HandleFunc("/invoice-list", invoice.GetInvoicesHandler(db))
	http.HandleFunc("/invoice-number-check", invoice.CheckInvoiceNumberExistsHandler(db))
	http.HandleFunc("/invoice-update", invoice.UpdateInvoiceHandler(db))
	http.HandleFunc("/invoice-delete", invoice.DeleteInvoiceHandler(db))
	http.HandleFunc("/invoice-export-csv", invoice.InvoicesExportCSVHandler(db))
	http.HandleFunc("/invoice-import-csv", invoice.InvoicesImportCSVHandler(db))
	http.HandleFunc("/invoice-import-txt", invoice.InvoicesImportTXTHandler(db))
	http.HandleFunc("/generate-xlsx", invoice.InvoiceGenerateXLSX)
	http.HandleFunc("/generate-pdf", invoice.InvoiceGeneratePDF)
	http.HandleFunc("/clients", invoice.GetClientsHandler(db))
	http.HandleFunc("/client-add", invoice.ClientAddHandler(db))
	http.HandleFunc("/client-update", invoice.ClientUpdateHandler(db))
	http.HandleFunc("/client-delete", invoice.ClientDeleteHandler(db))
	http.HandleFunc("/client-export", invoice.ClientExportCSVHandler(db))
	http.HandleFunc("/client-import", invoice.ClientImportHandler(db))
	http.HandleFunc("/vat-get", invoice.GetVatHandler)
	http.HandleFunc("/vat-update", invoice.UpdateVatHandler)

	// Start the server
	log.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
