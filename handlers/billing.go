package handlers

import (
	"github.com/TechnicallyMay/indigo/db"
	"net/http"
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
	renderTemplate(w, r, "billingNew", []string{"base"}, nil)
}
