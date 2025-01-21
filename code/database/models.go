// code/database/models.go
package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Job struct {
	JobName   string `json:"jobName"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	FullPrice string `json:"fullPrice"`
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

type InvoiceData struct {
	ParentName    string  `json:"parentName"`
	Address1      string  `json:"address1"`
	Address2      string  `json:"address2"`
	Phone         string  `json:"phone"`
	Email         string  `json:"email"`
	InvoiceNumber string  `json:"invoiceNumber"`
	InvoiceDate   string  `json:"invoiceDate"`
	Cost          float64 `json:"cost"`
	VAT           float64 `json:"vat"`
	Total         float64 `json:"total"`
	Jobs          []Job   `json:"jobs"`
}

var db *sql.DB

func InitDB() (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", "./db/data.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}
