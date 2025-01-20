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
        abbreviation TEXT NOT NULL
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
        price TEXT NOT NULL
    );
    `
	_, err = db.Exec(createJobsTable)
	if err != nil {
		log.Fatal(err)
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
	}{
		{
			ClientName:   "Johnathan Michael Doe",
			ParentName:   "Michael John Doe",
			Address1:     "1234 Elm Street",
			Address2:     "Apt 5B",
			Phone:        "(555) 123-4567",
			Email:        "johndoe@example.com",
			Abbreviation: "JMD",
		},
		{
			ClientName:   "Emily Elizabeth Smith",
			ParentName:   "James William Smith",
			Address1:     "987 Oak Avenue",
			Address2:     "",
			Phone:        "(555) 987-6543",
			Email:        "emilysmith@website.org",
			Abbreviation: "EES",
		},
		{
			ClientName:   "William Thomas Johnson",
			ParentName:   "David Robert Johnson",
			Address1:     "456 Maple Road",
			Address2:     "Suite 302",
			Phone:        "(555) 555-5555",
			Email:        "williamjohnson@service.net",
			Abbreviation: "WTJ",
		},
	}

	// Insert hard-coded clients data into the database
	for _, client := range clients {
		_, err := db.Exec(`
            INSERT INTO clients (clientName, parentName, address1, address2, phone, email, abbreviation)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        `, client.ClientName, client.ParentName, client.Address1, client.Address2, client.Phone, client.Email, client.Abbreviation)
		if err != nil {
			log.Fatal("Failed to insert client:", err)
		}
	}

	// Hard-coded jobs data
	jobs := []struct {
		JobName string
		Price   string
	}{
		{JobName: "Software Development Services", Price: "500"},
		{JobName: "Bug Fixing and Maintenance", Price: "400"},
		{JobName: "API Integration", Price: "250"},
		{JobName: "Database Design and Management", Price: "350"},
		{JobName: "Project Management", Price: "500"},
		{JobName: "Consultation Services", Price: "800"},
	}

	// Insert hard-coded jobs data into the database
	for _, job := range jobs {
		_, err := db.Exec(`
            INSERT INTO jobs (jobName, price)
            VALUES (?, ?)
        `, job.JobName, job.Price)
		if err != nil {
			log.Fatal("Failed to insert job:", err)
		}
	}

	log.Println("Database initialized and populated with demo data.")
}
