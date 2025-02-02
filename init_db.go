// init_db.go
package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open or create the SQLite database file
	db, err := sql.Open("sqlite3", "./db/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the clients table
	createClientsTable := `
    CREATE TABLE IF NOT EXISTS clients (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        clientName TEXT NOT NULL,
        parentName TEXT NOT NULL,
        address1 TEXT NOT NULL,
        address2 TEXT NOT NULL,
        phone TEXT NOT NULL,
        email TEXT NOT NULL,
        abbreviation TEXT NOT NULL,
		status BOOLEAN NOT NULL
    );
    `
	_, err = db.Exec(createClientsTable)
	if err != nil {
		log.Fatal(err)
	}

	// Create the jobs table
	createJobsTable := `
    CREATE TABLE IF NOT EXISTS jobs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        jobName TEXT NOT NULL,
        price TEXT NOT NULL,
		status BOOLEAN NOT NULL
    );
    `
	_, err = db.Exec(createJobsTable)
	if err != nil {
		log.Fatal(err)
	}

	// Create the vat table
	createVatTable := `
	CREATE TABLE IF NOT EXISTS vat (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		rate REAL NOT NULL
	);
	`
	_, err = db.Exec(createVatTable)
	if err != nil {
		log.Fatal(err)
	}

	// Create the invoices table
	createInvoicesTable := `
	CREATE TABLE IF NOT EXISTS invoices (
		invoiceNumber TEXT PRIMARY KEY,
		invoiceDate TEXT NOT NULL,
		clientName TEXT NOT NULL,
		parentName TEXT NOT NULL,
		address1 TEXT NOT NULL,
		address2 TEXT NOT NULL,
		phone TEXT NOT NULL,
		email TEXT NOT NULL,
		cost REAL NOT NULL,
		VAT REAL NOT NULL,
		total REAL NOT NULL
	);
	`
	_, err = db.Exec(createInvoicesTable)
	if err != nil {
		log.Fatal(err)
	}

	// Create the invoices-job-row table
	createInvoicesJobRowTable := `
	CREATE TABLE IF NOT EXISTS invoices_job_row (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		invoiceNumber TEXT NOT NULL,
		jobName TEXT NOT NULL,
		quantity TEXT NOT NULL,
		price TEXT NOT NULL,
		fullPrice TEXT NOT NULL,
		FOREIGN KEY (invoiceNumber) REFERENCES invoices(invoiceNumber)
	);
	`
	_, err = db.Exec(createInvoicesJobRowTable)
	if err != nil {
		log.Fatal(err)
	}

	// Set default vat value to 20%
	_, err = db.Exec("INSERT INTO vat (rate) VALUES (?)", 20.0) // Default VAT rate of 20%
	if err != nil {
		log.Fatal("Failed to insert default VAT rate:", err)
	}

	// Hard-coded clients data
	clients := []struct {
		ClientName   string
		ParentName   string
		Address1     string
		Address2     string
		Phone        string
		Email        string
		Abbreviation string
		Status       bool
	}{
		{
			ClientName:   "Johnathan Michael Doe",
			ParentName:   "Michael John Doe",
			Address1:     "1234 Elm Street",
			Address2:     "Apt 5B",
			Phone:        "(555) 123-4567",
			Email:        "johndoe@example.com",
			Abbreviation: "JMD",
			Status:       true,
		},
		{
			ClientName:   "Emily Elizabeth Smith",
			ParentName:   "James William Smith",
			Address1:     "987 Oak Avenue",
			Address2:     "",
			Phone:        "(555) 987-6543",
			Email:        "emilysmith@website.org",
			Abbreviation: "EES",
			Status:       true,
		},
		{
			ClientName:   "William Thomas Johnson",
			ParentName:   "David Robert Johnson",
			Address1:     "456 Maple Road",
			Address2:     "Suite 302",
			Phone:        "(555) 555-5555",
			Email:        "williamjohnson@service.net",
			Abbreviation: "WTJ",
			Status:       true,
		},
	}

	// Insert hard-coded clients data into the database
	for _, client := range clients {
		_, err := db.Exec(`
            INSERT INTO clients (clientName, parentName, address1, address2, phone, email, abbreviation, status)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        `, client.ClientName, client.ParentName, client.Address1, client.Address2, client.Phone, client.Email, client.Abbreviation, client.Status)
		if err != nil {
			log.Fatal("Failed to insert client:", err)
		}
	}

	// Hard-coded jobs data
	jobs := []struct {
		JobName string
		Price   string
		Status  bool
	}{
		{JobName: "Software Development Services", Price: "500", Status: true},
		{JobName: "Bug Fixing and Maintenance", Price: "400", Status: true},
		{JobName: "API Integration", Price: "250", Status: true},
		{JobName: "Database Design and Management", Price: "350", Status: true},
		{JobName: "Project Management", Price: "500", Status: true},
		{JobName: "Consultation Services", Price: "800", Status: true},
	}

	// Insert hard-coded jobs data into the database
	for _, job := range jobs {
		_, err := db.Exec(`
            INSERT INTO jobs (jobName, price, status)
            VALUES (?, ?, ?)
        `, job.JobName, job.Price, job.Status)
		if err != nil {
			log.Fatal("Failed to insert job:", err)
		}
	}

	log.Println("Database initialized and populated with demo data.")
}
