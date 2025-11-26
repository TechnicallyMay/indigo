package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"regexp"

	"codeberg.org/go-pdf/fpdf"
	"github.com/TechnicallyMay/indigo/db"
	"github.com/TechnicallyMay/indigo/handlers"
)

var validPath = regexp.MustCompile("^(/[a-zA-Z0-9]+/{0,1}?)*$")

func readFile(fileName string) []byte {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func BuildMail(filename string) []byte {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("From: %s\r\n", "postmaster@sandbox6799805213174515aa97c14ef663c50e.mailgun.org"))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", "masonwells01@gmail.com"))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", "some Email stuff"))

	boundary := "doesthismatter"
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/plain;\r\n")
	buf.WriteString(fmt.Sprintf("\r\n%s", "someEmailBody"))

	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/plain;\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
	buf.WriteString("Content-Disposition: attachment; filename=" + filename + "\r\n")
	buf.WriteString("Content-ID: <" + filename + ">\r\n\r\n")

	data := readFile(filename)
	b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(b, data)
	buf.Write(b)
	buf.WriteString(fmt.Sprintf("\r\n--%s", boundary))
	buf.WriteString("--")

	return buf.Bytes()
}

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

	msg := BuildMail("hello.pdf")
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
