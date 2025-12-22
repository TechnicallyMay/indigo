package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/TechnicallyMay/indigo/db"
	"github.com/TechnicallyMay/indigo/sender"
)

type BillingHandler struct {
	batchDb db.InvoiceBatchTable
	invDb   db.InvoiceTable
	custDb  db.CustomerTable
	prodDb  db.ProductTable
	itemDb  db.InvoiceItemTable
	notDb   db.InvoiceNotificationTable
	sender  sender.InvoiceSender
}

type billingData struct {
	Batch              db.InvoiceBatch
	IncludedCustomers  []db.Customer
	AvailableCustomers []db.Customer
}

var billingHandlerInstance *BillingHandler

func NewBillingHandler(batchDb db.InvoiceBatchTable,
	invDb db.InvoiceTable,
	custDb db.CustomerTable,
	prodDb db.ProductTable,
	itemDb db.InvoiceItemTable,
	sender sender.InvoiceSender,
	notDb db.InvoiceNotificationTable) *BillingHandler {
	if billingHandlerInstance == nil {
		billingHandlerInstance = &BillingHandler{batchDb: batchDb, invDb: invDb, custDb: custDb, prodDb: prodDb, itemDb: itemDb, sender: sender, notDb: notDb}
	}

	return billingHandlerInstance
}

func (h *BillingHandler) HandleGetBilling(w http.ResponseWriter, r *http.Request) {
	batches := h.batchDb.List()
	renderTemplate(w, r, newRenderOpts("invbatch/list", batches, "content", "header"))
}

func (h *BillingHandler) HandleGetNewBilling(w http.ResponseWriter, r *http.Request) {
	due := time.Now()
	if due.Day() >= 15 {
		due = due.AddDate(0, 1, 0)
	}
	due = time.Date(due.Year(), due.Month(), 15, 0, 0, 0, 0, due.Location())
	subj := "Invoice for " + due.Month().String()
	desc := "Hello,\n\nThis is the invoice for " + due.Month().String()

	newBatch := db.InvoiceBatch{State: db.Draft, DueDate: due.Unix(), FinishedSendingAt: 0, NotificationSubject: subj, NotificationDescription: desc}
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

	data, err := h.getBatchData(id)
	if err != nil {
		http.Error(w, "Invoice batch id "+idStr+" couldn't be parsed to an integer", 400)
		return
	}

	renderOpts := newRenderOpts("invbatch/batch", data, "content", "header")
	renderOpts.prereqTemplates = []string{"invbatch/customerPicker", "invbatch/detailsview"}
	renderTemplate(w, r, renderOpts)
}

func (h *BillingHandler) HandleEditBatchDetails(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		http.Error(w, "Invoice batch id "+idStr+" couldn't be parsed to an integer", 400)
		return
	}

	batch, err := h.batchDb.Get(id)
	handleHttpError(w, err, 404)

	if batch.State == db.Draft {
		renderTemplate(w, r, newRenderOpts("invbatch/detailsedit", &batch, "invBatchDetailsEdit"))
	} else {
		handleHttpError(w, errors.New("Can't edit details of a non-draft batch"), 400)
	}
}

func (h *BillingHandler) HandleViewBatchDetails(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		http.Error(w, "Invoice batch id "+idStr+" couldn't be parsed to an integer", 400)
		return
	}

	batch, err := h.batchDb.Get(id)
	handleHttpError(w, err, 404)

	renderTemplate(w, r, newRenderOpts("invbatch/detailsview", &batch, "invBatchDetailsView"))
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

	renderTemplate(w, r, newRenderOpts("invbatch/customerPicker", data, "customerPicker"))
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

	renderTemplate(w, r, newRenderOpts("invbatch/customerPicker", data, "customerPicker"))
}

func (h *BillingHandler) HandleUpdateInvoiceBatch(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	handleHttpError(w, err, 400)

	batch, err := h.batchDb.Get(id)
	handleHttpError(w, err, 500)

	if batch.State != db.Draft {
		handleHttpError(w, errors.New("Tried to update an invoice not in draft state"), 400)
	}

	err = r.ParseForm()
	handleHttpError(w, err, 400)

	subj := r.FormValue("subject")
	body := r.FormValue("body")
	dueDateStr := r.FormValue("dueDate")
	dueDate, err := time.Parse(time.DateOnly, dueDateStr)
	handleHttpError(w, err, 400)

	batch.NotificationSubject = subj
	batch.NotificationDescription = body
	batch.DueDate = dueDate.UnixMilli() / 1000

	h.batchDb.Update(batch)

	renderTemplate(w, r, newRenderOpts("invbatch/detailsview", &batch, "invBatchDetailsView"))
}

func (h *BillingHandler) HandleSendBatch(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	handleHttpError(w, err, 400)

	batch, err := h.batchDb.Get(id)
	handleHttpError(w, err, 500)

	if batch.State != db.Draft && batch.State != db.Failed {
		handleHttpError(w, errors.New("Invoice batches can only be sent from draft or failed state currently."), 400)
	}

	batch.State = db.Sending
	h.batchDb.Update(batch)

	go h.sendInvoices(batch)

	HtmxRefresh(w)
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

func (h *BillingHandler) getBatchData(id int64) (*billingData, error) {
	batch, err := h.batchDb.Get(id)

	if err != nil {
		return nil, err
	}

	incCusts := h.custDb.GetByInvoiceBatch(id)
	allCusts := h.getNonIncludedCustomers(incCusts)

	return &billingData{
		Batch:              batch,
		IncludedCustomers:  incCusts,
		AvailableCustomers: allCusts,
	}, nil
}

func (h *BillingHandler) sendInvoices(batch db.InvoiceBatch) error {
	query := db.CreateInvoiceQuery()
	query.BatchId = batch.Id

	invs, err := h.invDb.Query(query)
	if err != nil {
		return err
	}

	customers := h.custDb.List()
	custMap := make(map[int64]db.Customer, 0)

	for _, cust := range customers {
		custMap[cust.Id] = cust
	}

	products, err := h.prodDb.List()
	if err != nil {
		return err
	}

	prodMap := make(map[int64]db.Product, 0)
	for _, prod := range products {
		prodMap[prod.Id] = prod
	}

	var errs error
	failCount := 0
	for _, inv := range invs {
		sent, _ := h.sendInvoice(batch, inv, custMap[inv.CustomerId], prodMap)
		if !sent {
			failCount++
		}
	}

	if failCount == len(invs) {
		batch.State = db.Failed
	} else if failCount > 0 {
		batch.State = db.PartialFailure
	} else {
		batch.State = db.Sent
		batch.FinishedSendingAt = time.Now().Unix()
	}

	h.batchDb.Update(batch)

	return errs
}

func (h *BillingHandler) sendInvoice(batch db.InvoiceBatch, inv db.Invoice, cust db.Customer, prodMap map[int64]db.Product) (sent bool, error error) {
	items, error := h.itemDb.List(inv.Id)

	if error != nil {
		return
	}

	error = h.sender.SendInvoice(batch, cust, items, prodMap)

	notification := db.InvoiceNotification{
		InvoiceId:  inv.Id,
		DueDate:    batch.DueDate,
		Successful: true,
		SentAt:     time.Now().Unix(),
	}

	if error != nil {
		notification.Error = error.Error()
		notification.Successful = false
	} else {
		sent = true
	}

	_, err := h.notDb.Add(notification)
	error = errors.Join(err, error)

	return
}
