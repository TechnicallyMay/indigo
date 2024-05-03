package handlers

import (
	"net/http"
	"github.com/TechnicallyMay/indigo/db"
)

type CustomerHandler struct {
    db db.CustomerTable
}

var handlerInstance *CustomerHandler
func NewCustomerHandler(db db.CustomerTable) *CustomerHandler {
    if handlerInstance == nil {
        handlerInstance = &CustomerHandler{db: db}
    }

    return handlerInstance
}

func (h *CustomerHandler) HandleGetCustomers(w http.ResponseWriter, r *http.Request) {
    h.db.List()
    h.db.Add(db.Customer{FirstName: "heyo", LastName: "heyooo", Email: "test"})
    renderTemplate(w, r, "customers", []string{"base"})
}

func PostCustomersHandler(w http.ResponseWriter, r *http.Request) {
}

