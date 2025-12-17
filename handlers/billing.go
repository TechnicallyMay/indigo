package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TechnicallyMay/indigo/db"
)

type BillingHandler struct {
	batchDb db.InvoiceBatchTable
	invDb   db.InvoiceTable
	custDb  db.CustomerTable
}

type billingData struct {
	Batch              db.InvoiceBatch
	IncludedCustomers  []db.Customer
	AvailableCustomers []db.Customer
}

var billingHandlerInstance *BillingHandler

func NewBillingHandler(batchDb db.InvoiceBatchTable, invDb db.InvoiceTable, custDb db.CustomerTable) *BillingHandler {
	if billingHandlerInstance == nil {
		billingHandlerInstance = &BillingHandler{batchDb: batchDb, invDb: invDb, custDb: custDb}
	}

	return billingHandlerInstance
}

func (h *BillingHandler) HandleGetBilling(w http.ResponseWriter, r *http.Request) {
	batches := h.batchDb.List()
	renderTemplate(w, r, newRenderOpts("billingHome", batches, "content", "header"))
}

func (h *BillingHandler) HandleGetNewBilling(w http.ResponseWriter, r *http.Request) {
	newBatch := db.InvoiceBatch{State: db.Draft, DueDate: 0, FinishedSendingAt: 0}
	newId := h.batchDb.Add(newBatch)

	dest := fmt.Sprintf("/billing/%d", newId)
	HtmxHardRedirect(w, dest)
}

func (h *BillingHandler) HandleGetInvoiceBatch(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		http.Error(w, "Invoice batch id "+idStr+" couldn't be parsed to an integer", 400)
		return
	}

	batch, err := h.batchDb.Get(id)

	if err != nil {
		http.Error(w, "Invoice batch id '%v' couldn't be parsed to an integer", 400)
		return
	}

	incCusts := h.custDb.GetByInvoiceBatch(id)
	allCusts := h.getNonIncludedCustomers(incCusts)

	data := &billingData{
		Batch:              batch,
		IncludedCustomers:  incCusts,
		AvailableCustomers: allCusts,
	}

	renderOpts := newRenderOpts("billingById", data)
	renderOpts.prereqTemplates = []string{"customerPicker"}
	renderTemplate(w, r, renderOpts)
}

func (h *BillingHandler) HandleDeleteInvoiceFromBatch(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	bId, err := strconv.ParseInt(idStr, 10, 64)
	handleHttpError(w, err, 400)

	params := r.URL.Query()
	cId, err := strconv.ParseInt(params.Get("customerId"), 10, 64)
	handleHttpError(w, err, 400)

	incCusts := h.custDb.GetByInvoiceBatch(bId)

	for _, cust := range incCusts {
		if cust.Id == cId {
			handleHttpError(w, errors.New("The customer is already included in this batch"), 400)
			return
		}
	}

	batch, err := h.batchDb.Get(bId)
	handleHttpError(w, err, 500)
	addedCust, err := h.custDb.Get(cId)
	handleHttpError(w, err, 500)

	inv := db.Invoice{
		BatchId:         bId,
		CustomerId:      cId,
		CustomerVersion: addedCust.Version,
		IsPaid:          false,
	}

	h.invDb.Add(inv)
	newIncludedCusts := append(incCusts, addedCust)

	data := &billingData{
		Batch:              batch,
		IncludedCustomers:  newIncludedCusts,
		AvailableCustomers: h.getNonIncludedCustomers(newIncludedCusts),
	}
	fmt.Println(data.AvailableCustomers)

	renderTemplate(w, r, newRenderOpts("customerPicker", data, "customerPicker"))
}
func (h *BillingHandler) HandleAddInvoiceToBatch(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	custStr := r.PathValue("customerId")

	bId, err := strconv.ParseInt(idStr, 10, 64)
	handleHttpError(w, err, 400)

	cId, err := strconv.ParseInt(custStr, 10, 64)
	handleHttpError(w, err, 400)

	incCusts := h.custDb.GetByInvoiceBatch(bId)

	for _, cust := range incCusts {
		if cust.Id == cId {
			handleHttpError(w, errors.New("The customer is already included in this batch"), 400)
			return
		}
	}

	batch, err := h.batchDb.Get(bId)
	handleHttpError(w, err, 500)
	addedCust, err := h.custDb.Get(cId)
	handleHttpError(w, err, 500)

	inv := db.Invoice{
		BatchId:         bId,
		CustomerId:      cId,
		CustomerVersion: addedCust.Version,
		IsPaid:          false,
	}

	h.invDb.Add(inv)
	newIncludedCusts := append(incCusts, addedCust)

	data := &billingData{
		Batch:              batch,
		IncludedCustomers:  newIncludedCusts,
		AvailableCustomers: h.getNonIncludedCustomers(newIncludedCusts),
	}
	fmt.Println(data.AvailableCustomers)

	renderTemplate(w, r, newRenderOpts("customerPicker", data, "customerPicker"))
}

func (h *BillingHandler) getNonIncludedCustomers(incCusts []db.Customer) []db.Customer {
	allCusts := h.custDb.List()
	nonIncluded := make([]db.Customer, 0)
	for _, c := range allCusts {
		include := true
		for _, inc := range incCusts {
			if inc.Id == c.Id {
				include = false
				break
			}
		}

		if include {
			nonIncluded = append(nonIncluded, c)
		}
	}
	return nonIncluded
}
