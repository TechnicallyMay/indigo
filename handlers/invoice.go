package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TechnicallyMay/indigo/db"
)

type InvoiceHandler struct {
	db db.InvoiceTable
}

var invoiceHandlerInstance *InvoiceHandler

func NewInvoiceHandler(db db.InvoiceTable) *InvoiceHandler {
	if invoiceHandlerInstance == nil {
		invoiceHandlerInstance = &InvoiceHandler{db: db}
	}

	return invoiceHandlerInstance
}

func (h *InvoiceHandler) HandleGetInvoice(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	query := db.CreateInvoiceQuery()

	batchFilter, berr := strconv.ParseInt(params.Get("batchId"), 10, 64)
	customerIdFilter, cerr := strconv.ParseInt(params.Get("customerId"), 10, 64)
	customerVersionFilter, verr := strconv.ParseInt(params.Get("customerVersion"), 10, 64)
	isPaidFilter, perr := strconv.ParseBool(params.Get("isPaid"))

	if berr == nil {
		query.BatchId = batchFilter
		log.Println("batch id", query.BatchId)
	}
	if cerr == nil {
		query.CustomerId = customerIdFilter
		log.Println("customer id", query.CustomerId)
	}
	if verr == nil {
		query.CustomerVersion = customerVersionFilter
	}
	if perr == nil {
		if isPaidFilter {
			query.IsPaid = 1
		} else {
			query.IsPaid = 0
		}
	}

	invoice, err := h.db.QueryRow(query)

	if err != nil {
		http.Error(w, "Error while querying invoice", 400)
		return
	}
	log.Println("Found an invoice with id", invoice.Id)

	renderTemplate(w, r, newRenderOpts("invoice", invoice))
}
