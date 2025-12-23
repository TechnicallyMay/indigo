package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TechnicallyMay/indigo/db"
	"github.com/TechnicallyMay/indigo/pdf"
)

type InvoiceHandler struct {
	batchDb    db.InvoiceBatchTable
	invDb      db.InvoiceTable
	custDb     db.CustomerTable
	invItemDb  db.InvoiceItemTable
	prodDb     db.ProductTable
	settingsDb db.SettingsTable
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

func NewInvoiceHandler(invDb db.InvoiceTable, custDb db.CustomerTable, invItemDb db.InvoiceItemTable, prodDb db.ProductTable, batchDb db.InvoiceBatchTable, settingsDb db.SettingsTable) *InvoiceHandler {
	if invoiceHandlerInstance == nil {
		invoiceHandlerInstance = &InvoiceHandler{invDb: invDb, custDb: custDb, invItemDb: invItemDb, prodDb: prodDb, batchDb: batchDb, settingsDb: settingsDb}
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

func (h *InvoiceHandler) HandleGetInvoice(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	handleHttpError(w, err, 400)

	data, err := h.getInvoiceData(id)
	handleHttpError(w, err, 500)

	opts := newRenderOpts("invoice", data)
	opts.prereqTemplates = []string{"viewInvoiceItem"}
	renderTemplate(w, r, opts)
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

	data, err := h.getInvoiceData(invoice.Id)
	handleHttpError(w, err, 500)

	opts := newRenderOpts("invoice", data)
	opts.prereqTemplates = []string{"viewInvoiceItem"}
	renderTemplate(w, r, opts)
}

func (h *InvoiceHandler) HandlePreviewInvoicePdf(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	handleHttpError(w, err, 400)

	data, err := h.getInvoiceData(id)
	handleHttpError(w, err, 500)

	batch, err := h.batchDb.Get(data.Invoice.BatchId)
	handleHttpError(w, err, 500)
	settings, err := h.settingsDb.GetAll()
	handleHttpError(w, err, 500)
	path, err := pdf.MakeInvoicePdf(settings, batch, data.Invoice, data.Customer, data.Items, data.Products)
	handleHttpError(w, err, 500)
	AddAttachment(w, path, "invoice.pdf")
}

func (h *InvoiceHandler) getInvoiceData(id int64) (*invoiceData, error) {
	inv, err := h.invDb.Get(id)
	if err != nil {
		return nil, err
	}

	items, err := h.invItemDb.List(id)
	if err != nil {
		return nil, err
	}

	products, err := h.prodDb.List()
	if err != nil {
		return nil, err
	}

	productMap := make(map[int64]db.Product, 0)
	for _, prod := range products {
		productMap[prod.Id] = prod
	}

	customer, err := h.custDb.Get(inv.CustomerId)
	if err != nil {
		return nil, err
	}

	return &invoiceData{
		Invoice:  inv,
		Items:    items,
		Products: productMap,
		Customer: customer,
	}, nil
}
