package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TechnicallyMay/indigo/db"
)

type InvoiceHandler struct {
	invDb     db.InvoiceTable
	custDb    db.CustomerTable
	invItemDb db.InvoiceItemTable
	prodDb    db.ProductTable
}

// Data required to populate the invoice template
type invoiceData struct {
	Invoice  db.Invoice
	Customer db.Customer
	Items    []db.InvoiceItem
	Products map[int64]db.Product
}

func (d *invoiceData) GetInvoiceItemData(item db.InvoiceItem) invoiceItemData {
	product := d.Products[item.ProductId]
	return invoiceItemData{
		InvoiceItem: item,
		Product:     product,
	}
}

func (d *invoiceData) GetUnusedProducts() []db.Product {
	res := make([]db.Product, 0)

	for _, product := range d.Products {
		canUse := true
		for _, item := range d.Items {
			if item.ProductId == product.Id {
				canUse = false
				break
			}
		}

		if canUse {
			res = append(res, product)
		}
	}

	return res
}

var invoiceHandlerInstance *InvoiceHandler

func NewInvoiceHandler(invDb db.InvoiceTable, custDb db.CustomerTable, invItemDb db.InvoiceItemTable, prodDb db.ProductTable) *InvoiceHandler {
	if invoiceHandlerInstance == nil {
		invoiceHandlerInstance = &InvoiceHandler{invDb: invDb, custDb: custDb, invItemDb: invItemDb, prodDb: prodDb}
	}

	return invoiceHandlerInstance
}

func (h *InvoiceHandler) HandleDeleteInvoice(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	handleHttpError(w, err, 400)

	err = h.invDb.Delete(id)
	handleHttpError(w, err, 500)

	HtmxRefresh(w)
}

func (h *InvoiceHandler) HandleQueryInvoice(w http.ResponseWriter, r *http.Request) {
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

	invoice, err := h.invDb.QueryRow(query)
	handleHttpError(w, err, 500)
	items, err := h.invItemDb.List(invoice.Id)
	handleHttpError(w, err, 500)
	products, err := h.prodDb.List()
	handleHttpError(w, err, 500)
	customer, err := h.custDb.Get(invoice.CustomerId)
	handleHttpError(w, err, 500)

	productMap := make(map[int64]db.Product, 0)
	for _, prod := range products {
		productMap[prod.Id] = prod
	}

	data := &invoiceData{
		Invoice:  *invoice,
		Items:    items,
		Products: productMap,
		Customer: customer,
	}

	opts := newRenderOpts("invoice", data)
	opts.prereqTemplates = []string{"viewInvoiceItem"}
	renderTemplate(w, r, opts)
}
