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
	invId := r.PathValue("invoiceId")
	prodId := r.PathValue("productId")

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
	renderTemplate(w, r, newRenderOpts("viewInvoiceItem", data, "invoiceItem"))
}

func (h *InvoiceItemHandler) HandleDeleteInvoiceItem(w http.ResponseWriter, r *http.Request) {
	invId := r.PathValue("invoiceId")
	prodId := r.PathValue("productId")

	parsedInvId, err := strconv.ParseInt(invId, 10, 64)
	handleHttpError(w, err, 500)
	parsedProdId, err := strconv.ParseInt(prodId, 10, 64)
	handleHttpError(w, err, 500)

	err = h.invItemDb.Delete(parsedInvId, parsedProdId)
	handleHttpError(w, err, 500)

	renderTemplate(w, r, newRenderOpts("empty", nil))
}

func (h *InvoiceItemHandler) HandleUpdateInvoiceItem(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	handleHttpError(w, err, 500)

	invId := r.FormValue("invoiceId")
	prodId := r.FormValue("productId")
	qty := r.FormValue("quantity")

	parsedInvId, err := strconv.ParseInt(invId, 10, 64)
	handleHttpError(w, err, 500)
	parsedProdId, err := strconv.ParseInt(prodId, 10, 64)
	handleHttpError(w, err, 500)
	parsedQty, err := strconv.ParseInt(qty, 10, 64)
	handleHttpError(w, err, 500)

	product, err := h.prodDb.Get(parsedProdId)
	handleHttpError(w, err, 500)

	item := db.InvoiceItem{
		InvoiceId: parsedInvId,
		ProductId: parsedProdId,
		Quantity:  parsedQty,
	}

	err = h.invItemDb.Update(item)
	handleHttpError(w, err, 500)

	data := invoiceItemData{
		InvoiceItem: item,
		Product:     product,
	}
	renderTemplate(w, r, newRenderOpts("viewInvoiceItem", data, "invoiceItem"))
}

func (h *InvoiceItemHandler) HandleGetInvoiceItem(w http.ResponseWriter, r *http.Request) {
	invIdStr := r.PathValue("invoiceId")
	prodIdStr := r.PathValue("productId")

	invId, err := strconv.ParseInt(invIdStr, 10, 64)
	handleHttpError(w, err, 500)
	prodId, err := strconv.ParseInt(prodIdStr, 10, 64)
	handleHttpError(w, err, 500)

	item, err := h.invItemDb.Get(invId, prodId)
	handleHttpError(w, err, 500)

	product, err := h.prodDb.Get(prodId)
	handleHttpError(w, err, 500)

	data := invoiceItemData{
		InvoiceItem: *item,
		Product:     product,
	}

	renderTemplate(w, r, newRenderOpts("viewInvoiceItem", data, "invoiceItem"))
}

func (h *InvoiceItemHandler) HandleGetInvoiceItemEditForm(w http.ResponseWriter, r *http.Request) {
	invIdStr := r.PathValue("invoiceId")
	prodIdStr := r.PathValue("productId")

	invId, err := strconv.ParseInt(invIdStr, 10, 64)
	handleHttpError(w, err, 500)
	prodId, err := strconv.ParseInt(prodIdStr, 10, 64)
	handleHttpError(w, err, 500)

	item, err := h.invItemDb.Get(invId, prodId)
	handleHttpError(w, err, 500)

	product, err := h.prodDb.Get(prodId)
	handleHttpError(w, err, 500)

	data := invoiceItemData{
		InvoiceItem: *item,
		Product:     product,
	}

	renderTemplate(w, r, newRenderOpts("editInvoiceItem", data))
}
