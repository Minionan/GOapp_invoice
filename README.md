# GOapp_invoice

This is a simple invoice-generating app written with JS and GO.
It is a good start when creating your web app where invoicing functionality is required.
This version of the app can generate the invoice in txt, xlsx and pdf format.

## Setup

Clone the repository with `git clone https://github.com/Minionan/GOapp_invoice.git`

### Initialising user database

1. Run `init_db.go` script by typing in terminal `go run init_db.go`
2. Verify if a new SQLite database file was created in db folder
3. The default data.db file has 3 default client and 6 job records

### Installing libreoffice

Install libre office on your system
For linux system a headless installation with calc will suffice

### Replacing template.xlsx file

Please replace the template.xlsx file with your company file.
Be advised that libreoffice is used for pdf generation, thus formatting might differ from one used by excel. Some adjustments might be required to get the pdf output file to look as intended.

## Run app

1. Run `go mod tidy`
2. Run `go run main.go`
