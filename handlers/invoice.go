package handlers

import (
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
	// id := r.PathValue("id")
	// fmt.Println("Invoice with id", id)

	// idInt, err := strconv.ParseInt(id, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Error while retrieving invoice: Id couldn't be parsed as an int.", 400)
	// 	return
	// }
	params := r.URL.Query()

	query := db.CreateInvoiceQuery()

	batchFilter, berr := strconv.ParseInt(params.Get("batchId"), 10, 64)
	customerIdFilter, cerr := strconv.ParseInt(params.Get("customerId"), 10, 64)
	customerVersionFilter, verr := strconv.ParseInt(params.Get("customerVersion"), 10, 64)
	isPaidFilter, perr := strconv.ParseBool(params.Get("isPaid"))

	if berr == nil {
		query.BatchId = batchFilter
	}
	if cerr != nil {
		query.CustomerId = customerIdFilter
	}
	if verr != nil {
		query.CustomerVersion = customerVersionFilter
	}
	if perr != nil {
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

	renderTemplate(w, r, "invoice", nil, invoice)
}
