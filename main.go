package main

import (
	"log"
	"net/http"

	"github.com/TechnicallyMay/indigo/db"
	"github.com/TechnicallyMay/indigo/handlers"
)

func main() {
	log.Println("Starting!")

	// // PDF
	// log.Println("making example pdf")
	// pdfPath, err := pdf.MakeInvoicePdf()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// // MAIL
	// log.Println("Sending an email")
	//
	// var smtpPass string
	// fmt.Println("Enter the smtp pass")
	// fmt.Scan(&smtpPass)
	//
	// auth := smtp.PlainAuth("", "postmaster@sandbox6799805213174515aa97c14ef663c50e.mailgun.org", smtpPass, "smtp.mailgun.org")
	// client := mail.SmtpClient{Host: "smtp.mailgun.org", Port: 587, Auth: auth}
	//
	// msg := mail.Mail{From: "mail@chickpea-home.duckdns.org", To: []string{"masonwells01@gmail.com"}, Subject: "Test Mail", Body: "Test Message Body", AttachmentFilePath: pdfPath}
	// err = client.SendMail(msg)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// MAIN

	pool := db.OpenDb()
	defer pool.Close()

	custTable := db.InitCustomerTable(pool)
	invoiceBatchTable := db.InitInvoiceBatchTable(pool)
	invoiceTable := db.InitInvoiceTable(pool)
	invoiceItemTable := db.InitInvoiceItemTable(pool)
	productsTable := db.InitProductTable(pool)

	customerHandler := handlers.NewCustomerHandler(*custTable)
	billingHandler := handlers.NewBillingHandler(*invoiceBatchTable, *invoiceTable, *custTable)
	invoiceHandler := handlers.NewInvoiceHandler(*invoiceTable, *custTable, *invoiceItemTable, *productsTable)
	invoiceItemHandler := handlers.NewInvoiceItemHandler(*invoiceItemTable, *productsTable)
	productsHandler := handlers.NewProductHandler(*productsTable)

	mux := http.NewServeMux()
	mux.Handle("GET /js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
	mux.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))

	mux.HandleFunc("GET /{$}", handlers.GetRootHandler)

	mux.HandleFunc("GET /billing/", billingHandler.HandleGetBilling)
	mux.HandleFunc("GET /billing/new", billingHandler.HandleGetNewBilling)
	mux.HandleFunc("GET /billing/{id}", billingHandler.HandleGetInvoiceBatch)
	mux.HandleFunc("PUT /billing/{id}/invoice", billingHandler.HandleAddInvoiceToBatch)

	mux.HandleFunc("GET /customers/", customerHandler.HandleGetCustomers)
	mux.HandleFunc("GET /customers/new", customerHandler.HandleGetEditCustomerForm)
	mux.HandleFunc("GET /customers/edit/{id}", customerHandler.HandleGetEditCustomerForm)
	mux.HandleFunc("POST /customers", customerHandler.HandlePostCustomer)
	mux.HandleFunc("PUT /customers/{id}", customerHandler.HandlePutCustomer)

	mux.HandleFunc("GET /invoice", invoiceHandler.HandleGetInvoice)

	mux.HandleFunc("POST /invoiceItem", invoiceItemHandler.HandleNewInvoiceItem)
	mux.HandleFunc("PUT /invoiceItem/{invoiceId}/{productId}", invoiceItemHandler.HandleUpdateInvoiceItem)
	mux.HandleFunc("GET /invoiceItem/{invoiceId}/{productId}", invoiceItemHandler.HandleGetInvoiceItem)
	mux.HandleFunc("GET /invoiceItem/edit/{invoiceId}/{productId}", invoiceItemHandler.HandleGetInvoiceItemEditForm)

	mux.HandleFunc("GET /products/", productsHandler.HandleGetProducts)
	mux.HandleFunc("GET /products/new", productsHandler.HandleGetEditProductForm)
	mux.HandleFunc("GET /products/edit/{id}", productsHandler.HandleGetEditProductForm)
	mux.HandleFunc("POST /products", productsHandler.HandlePostProduct)
	mux.HandleFunc("PUT /products/{id}", productsHandler.HandlePutProduct)

	mux.HandleFunc("GET /records/", handlers.HandleGetRecords)

	mux.HandleFunc("GET /settings/", handlers.HandleGetSettings)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
