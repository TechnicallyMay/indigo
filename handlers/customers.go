package handlers

import (
	"log"
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
    customers := h.db.List()
    renderTemplate(w, r, "customers", []string{"base"}, customers)
}

func (h *CustomerHandler) HandleGetNewCustomer(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, r, "newCustomer", []string{"base"}, nil)
}

func (h *CustomerHandler) HandlePostCustomer(w http.ResponseWriter, r *http.Request) {
    log.Println("Adding a new customer")
    r.ParseForm()
    log.Println(r.PostForm)

    newCustomer:= db.Customer{FirstName: r.PostForm.Get("firstName"), LastName: r.PostForm.Get("lastName"), Email: r.PostForm.Get("email") }

    log.Printf("Got customer %v %v %v\n", newCustomer.FirstName, newCustomer.LastName, newCustomer.Email)

    h.db.Add(newCustomer)

    HtmxRedirect(w, "/customers", "#main-content")
}

