package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/TechnicallyMay/indigo/appsettings"
	"github.com/TechnicallyMay/indigo/db"
	"github.com/TechnicallyMay/indigo/handlers"
	"github.com/TechnicallyMay/indigo/mail"
	"github.com/TechnicallyMay/indigo/sender"
)

const DefaultSettingsPath = "indigo.toml"

func main() {
	log.Println("Starting!")

	path := os.Getenv("INDIGO_CONFIG")
	if path == "" {
		path = DefaultSettingsPath
	}

	settings, err := appsettings.GetSettings(path)
	if err != nil {
		log.Fatal("Unable to read settings file. Ensure a settings file is available at " + path)
	}

	passBytes, err := os.ReadFile(settings.Smtp.PassFile)
	smtpPass := strings.TrimSpace(string(passBytes))

	auth := smtp.PlainAuth("", settings.Smtp.From, smtpPass, settings.Smtp.Host)
	client := mail.SmtpClient{Host: settings.Smtp.Host, Port: settings.Smtp.Port, Auth: auth}
	sender := sender.InvoiceSender{MailClient: &client}

	pool := db.OpenDb(settings)
	defer pool.Close()

	custTable := db.InitCustomerTable(pool)
	invoiceBatchTable := db.InitInvoiceBatchTable(pool)
	invoiceTable := db.InitInvoiceTable(pool)
	invoiceItemTable := db.InitInvoiceItemTable(pool)
	productsTable := db.InitProductTable(pool)
	notTable := db.InitInvoiceNotificationTable(pool)
	settingsTable := db.InitSettingsTable(pool)

	customerHandler := handlers.NewCustomerHandler(*custTable)
	billingHandler := handlers.NewBillingHandler(*invoiceBatchTable, *invoiceTable, *custTable, *productsTable, *invoiceItemTable, sender, *notTable, *settingsTable, settings)
	invoiceHandler := handlers.NewInvoiceHandler(*invoiceTable, *custTable, *invoiceItemTable, *productsTable, *invoiceBatchTable, *settingsTable)
	invoiceItemHandler := handlers.NewInvoiceItemHandler(*invoiceItemTable, *productsTable, *invoiceHandler)
	productsHandler := handlers.NewProductHandler(*productsTable)
	settingsHandler := handlers.NewSettingsHandler(*settingsTable)

	mux := http.NewServeMux()
	mux.Handle("GET /js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
	mux.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	mux.Handle("GET /image/", http.StripPrefix("/image/", http.FileServer(http.Dir("./static/image"))))

	mux.HandleFunc("GET /{$}", handlers.GetRootHandler)

	mux.HandleFunc("GET /billing/", billingHandler.HandleGetBilling)
	mux.HandleFunc("GET /billing/new", billingHandler.HandleGetNewBilling)
	mux.HandleFunc("GET /billing/{id}", billingHandler.HandleGetInvoiceBatch)
	mux.HandleFunc("GET /billing/edit/{id}", billingHandler.HandleEditBatchDetails)
	mux.HandleFunc("GET /billing/details/{id}", billingHandler.HandleViewBatchDetails)
	mux.HandleFunc("PUT /billing/{id}", billingHandler.HandleUpdateInvoiceBatch)
	mux.HandleFunc("POST /billing/duplicate/{id}", billingHandler.HandleDuplicateInvoiceBatch)
	mux.HandleFunc("POST /billing/send/{id}", billingHandler.HandleSendBatch)
	mux.HandleFunc("POST /billing/{id}/invoice/{customerId}", billingHandler.HandleAddInvoiceToBatch)

	mux.HandleFunc("GET /customers/", customerHandler.HandleGetCustomers)
	mux.HandleFunc("GET /customers/new", customerHandler.HandleGetEditCustomerForm)
	mux.HandleFunc("GET /customers/edit/{id}", customerHandler.HandleGetEditCustomerForm)
	mux.HandleFunc("POST /customers", customerHandler.HandlePostCustomer)
	mux.HandleFunc("PUT /customers/{id}", customerHandler.HandlePutCustomer)

	mux.HandleFunc("GET /invoice", invoiceHandler.HandleQueryInvoice)
	mux.HandleFunc("GET /invoice/{id}", invoiceHandler.HandleGetInvoice)
	mux.HandleFunc("DELETE /invoice/{id}", invoiceHandler.HandleDeleteInvoice)
	mux.HandleFunc("GET /invoice/preview/{id}", invoiceHandler.HandlePreviewInvoicePdf)
	mux.HandleFunc("GET /invoice/preview/sample", invoiceHandler.HandleSampleInvoicePdf)

	mux.HandleFunc("POST /invoiceItem/{invoiceId}/{productId}", invoiceItemHandler.HandleNewInvoiceItem)
	mux.HandleFunc("PUT /invoiceItem/{invoiceId}/{productId}", invoiceItemHandler.HandleUpdateInvoiceItem)
	mux.HandleFunc("GET /invoiceItem/{invoiceId}/{productId}", invoiceItemHandler.HandleGetInvoiceItem)
	mux.HandleFunc("DELETE /invoiceItem/{invoiceId}/{productId}", invoiceItemHandler.HandleDeleteInvoiceItem)
	mux.HandleFunc("GET /invoiceItem/edit/{invoiceId}/{productId}", invoiceItemHandler.HandleGetInvoiceItemEditForm)

	mux.HandleFunc("GET /products/", productsHandler.HandleGetProducts)
	mux.HandleFunc("GET /products/new", productsHandler.HandleGetEditProductForm)
	mux.HandleFunc("GET /products/edit/{id}", productsHandler.HandleGetEditProductForm)
	mux.HandleFunc("POST /products", productsHandler.HandlePostProduct)
	mux.HandleFunc("PUT /products/{id}", productsHandler.HandlePutProduct)

	mux.HandleFunc("GET /records/", handlers.HandleGetRecords)

	mux.HandleFunc("GET /settings/", settingsHandler.HandleGetSettings)
	mux.HandleFunc("PUT /settings/", settingsHandler.HandlePutSettings)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	http.ServeFile(w, r, "/static/favicon.ico")
}
