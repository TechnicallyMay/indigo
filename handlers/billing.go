package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TechnicallyMay/indigo/db"
)

type BillingHandler struct {
	db db.InvoiceBatchTable
}

type billingData struct {
	Batch             db.InvoiceBatch
	IncludedCustomers []db.Customer
}

var billingHandlerInstance *BillingHandler

func NewBillingHandler(db db.InvoiceBatchTable) *BillingHandler {
	if billingHandlerInstance == nil {
		billingHandlerInstance = &BillingHandler{db: db}
	}

	return billingHandlerInstance
}

func (h *BillingHandler) HandleGetBilling(w http.ResponseWriter, r *http.Request) {
	batches := h.db.List()
	renderTemplate(w, r, "billingHome", nil, batches)
}

func (h *BillingHandler) HandleGetNewBilling(w http.ResponseWriter, r *http.Request) {
	newBatch := db.InvoiceBatch{State: db.Draft, DueDate: 0, FinishedSendingAt: 0}
	newId := h.db.Add(newBatch)

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

	fmt.Println("Billing with id", id)

	batch, err := h.db.Get(id)

	if err != nil {
		http.Error(w, "Invoice batch id '%v' couldn't be parsed to an integer", 400)
		return
	}

	//TODO: Correctly populate this data
	data := &billingData{
		Batch: batch,
		IncludedCustomers: []db.Customer{
			db.Customer{
				Id:        1,
				Version:   1,
				CreatedAt: 0,

				FirstName: "test",
				LastName:  "whatever",
				Email:     "masil",
			},
			db.Customer{
				Id:        1,
				Version:   1,
				CreatedAt: 0,

				FirstName: "test",
				LastName:  "whatever2",
				Email:     "masil",
			},
		},
	}

	renderTemplate(w, r, "billingById", []string{"customerPicker"}, data)
}
