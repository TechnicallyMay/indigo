package handlers

import (
	"net/http"
	"strconv"

	"github.com/TechnicallyMay/indigo/db"
)

type InvoiceItemHandler struct {
	invItemDb db.InvoiceItemTable
	prodDb    db.ProductTable
}

// Data required to populate the invoiceItem template
type invoiceItemData struct {
	InvoiceItem db.InvoiceItem
	Product     db.Product
}

var invoiceItemHandlerInstance *InvoiceItemHandler

func NewInvoiceItemHandler(invItemDb db.InvoiceItemTable, prodDb db.ProductTable) *InvoiceItemHandler {
	if invoiceItemHandlerInstance == nil {
		invoiceItemHandlerInstance = &InvoiceItemHandler{invItemDb: invItemDb, prodDb: prodDb}
	}

	return invoiceItemHandlerInstance
}

func (h *InvoiceItemHandler) HandleNewInvoiceItem(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	invId := query.Get("invoiceId")
	prodId := query.Get("productId")

	parsedInvId, err := strconv.ParseInt(invId, 10, 64)
	handleHttpError(w, err, 500)
	parsedProdId, err := strconv.ParseInt(prodId, 10, 64)
	handleHttpError(w, err, 500)

	product, err := h.prodDb.Get(parsedProdId)
	handleHttpError(w, err, 500)

	item := db.InvoiceItem{
		InvoiceId: parsedInvId,
		ProductId: parsedProdId,
		Quantity:  1,
	}

	err = h.invItemDb.Add(item)
	handleHttpError(w, err, 500)

	data := invoiceItemData{
		InvoiceItem: item,
		Product:     product,
	}
	opts := newRenderOpts("invoiceItem", data)
	opts.entrypoint = "invoiceItem"
	renderTemplate(w, r, opts)
}
