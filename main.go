package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("static/invoice.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		jobsData, err := os.ReadFile("static/jobs.json")
		if err != nil {
			http.Error(w, "Failed to read jobs data", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jobsData)
	})

	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		clientsData, err := os.ReadFile("static/clients.json")
		if err != nil {
			http.Error(w, "Failed to read clients data", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(clientsData)
	})

	http.HandleFunc("/generate-invoice", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		clientName := r.FormValue("clientName")
		invoiceDate := r.FormValue("invoiceDate")
		invoiceNumber := r.FormValue("invoiceNumber")

		jobs := []struct {
			JobDescription string `json:"jobDescription"`
			Quantity       string `json:"quantity"`
			Price          string `json:"price"`
			FullPrice      string `json:"fullPrice"`
		}{}

		for i := 0; ; i++ {
			jobDescription := r.FormValue(fmt.Sprintf("job[%d][description]", i))
			if jobDescription == "" {
				break
			}
			jobs = append(jobs, struct {
				JobDescription string `json:"jobDescription"`
				Quantity       string `json:"quantity"`
				Price          string `json:"price"`
				FullPrice      string `json:"fullPrice"`
			}{
				JobDescription: jobDescription,
				Quantity:       r.FormValue(fmt.Sprintf("job[%d][quantity]", i)),
				Price:          r.FormValue(fmt.Sprintf("job[%d][price]", i)),
				FullPrice:      r.FormValue(fmt.Sprintf("job[%d][fullPrice]", i)),
			})
		}

		cost := r.FormValue("cost")
		vat := r.FormValue("vat")
		total := r.FormValue("total")

		invoiceContent := fmt.Sprintf(
			"Client Name: %s\nInvoice Date: %s\nInvoice Number: %s\n",
			clientName, invoiceDate, invoiceNumber,
		)
		for _, job := range jobs {
			invoiceContent += fmt.Sprintf(
				"Job Description: %s\nQuantity: %s\nPrice: %s\nFull Price: %s\n",
				job.JobDescription, job.Quantity, job.Price, job.FullPrice,
			)
		}
		invoiceContent += fmt.Sprintf("Cost: %s\nVAT (5%%): %s\nTotal: %s\n", cost, vat, total)

		w.Header().Set("Content-Disposition", "attachment; filename="+invoiceNumber+".txt")
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(invoiceContent))
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
