package handlers

import (
	"fmt"
	"net/http"

	"github.com/TechnicallyMay/indigo/db"
)

type BillingHandler struct {
	db db.InvoiceBatchTable
}

var billingHandlerInstance *BillingHandler

func NewBillingHandler(db db.InvoiceBatchTable) *BillingHandler {
	if billingHandlerInstance == nil {
		billingHandlerInstance = &BillingHandler{db: db}
	}

	return billingHandlerInstance
}

func (h *BillingHandler) HandleGetBilling(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "billing", []string{"base"}, nil)
}

func (h *BillingHandler) HandleGetNewBilling(w http.ResponseWriter, r *http.Request) {
	newBatch := db.InvoiceBatch{State: db.Draft, DueDate: 0, FinishedSendingAt: 0}
	newId := h.db.Add(newBatch)

	dest := fmt.Sprintf("/billing/%d", newId)
	// HtmxSoftRedirect(w, dest, "#main-content")
	HtmxHardRedirect(w, dest)
}

func (h *BillingHandler) HandleGetInvoiceBatch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Println("Billing with id", id)
	renderTemplate(w, r, "billingNew", []string{"base"}, nil)
}
