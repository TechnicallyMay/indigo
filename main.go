package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"regexp"

	"codeberg.org/go-pdf/fpdf"
	"github.com/TechnicallyMay/indigo/db"
	"github.com/TechnicallyMay/indigo/handlers"
)

var validPath = regexp.MustCompile("^(/[a-zA-Z0-9]+/{0,1}?)*$")

func main() {
	log.Println("Starting!")

	// PDF
	log.Println("making example pdf")
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "you owe me so much money")
	pdf.OutputFileAndClose("hello.pdf")

	// MAIL
	log.Println("Sending an email")
	host := "smtp.mailgun.org"
	addr := host + ":587"
	var smtpPass string
	fmt.Println("Enter the smtp pass")
	fmt.Scan(&smtpPass)
	auth := smtp.PlainAuth("", "postmaster@sandbox6799805213174515aa97c14ef663c50e.mailgun.org", smtpPass, host)

	msg := []byte("From: postmaster@sandbox6799805213174515aa97c14ef663c50e.mailgun.org\r\n" +
		"To: masonwells01@gmail.com\r\n" +
		"Subject: Test mail\r\n\r\n" +
		"Email body\r\n")
	mail := Multipart
	err := smtp.SendMail(addr, auth, "mail@chickpea-home.duckdns.org", []string{"masonwells01@gmail.com"}, msg)
	if err != nil {
		log.Fatal(err)
	}

	pool := db.OpenDb()
	defer pool.Close()

	customerHandler := handlers.NewCustomerHandler(*db.InitCustomerTable(pool))

	mux := http.NewServeMux()
	mux.Handle("GET /js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
	mux.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))

	mux.HandleFunc("GET /{$}", handlers.GetRootHandler)
	mux.HandleFunc("GET /home/", handlers.GetHomeHandler)

	mux.HandleFunc("GET /billing/", handlers.GetBillingHandler)

	mux.HandleFunc("GET /customers/", customerHandler.HandleGetCustomers)
	mux.HandleFunc("GET /customers/new", customerHandler.HandleGetAddOrUpdateCustomerForm)
	mux.HandleFunc("GET /customers/new/{id}", customerHandler.HandleGetAddOrUpdateCustomerForm)
	mux.HandleFunc("POST /customers", customerHandler.HandlePostCustomer)
	mux.HandleFunc("PUT /customers/{id}", customerHandler.HandlePutCustomer)

	mux.HandleFunc("GET /products/", handlers.GetProductsHandler)

	mux.HandleFunc("GET /records/", handlers.GetRecordsHandler)

	mux.HandleFunc("GET /reports/", handlers.GetReportsHandler)

	mux.HandleFunc("GET /settings/", handlers.GetSettingsHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
