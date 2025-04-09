package handlers

import (
	"log"
	"strconv"
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

func (h *CustomerHandler) HandleGetAddOrUpdateCustomerForm(w http.ResponseWriter, r *http.Request) {
	id:=r.PathValue("id")

	if id != "" {
		id, err:= strconv.ParseInt(r.PathValue("id"), 10, 64);
		if err != nil {
			http.Error(w, "Customer id '%v' could not be parsed to an Int64", 400);
		}

		cust:= h.db.Get(id)
		renderTemplate(w, r, "addOrUpdateCustomer", []string{"base"}, cust)
	} else {
		renderTemplate(w, r, "addOrUpdateCustomer", []string{"base"}, nil)
	}
}

func (h *CustomerHandler) HandlePostCustomer(w http.ResponseWriter, r *http.Request) {
    log.Println("Adding a new customer")
    r.ParseForm()

    newCustomer:= db.Customer{FirstName: r.PostForm.Get("firstName"), LastName: r.PostForm.Get("lastName"), Email: r.PostForm.Get("email") }

    log.Printf("Got customer %v %v %v\n", newCustomer.FirstName, newCustomer.LastName, newCustomer.Email)

    h.db.Add(newCustomer)

	HtmxRedirect(w, "/customers", "#main-content")
}

func (h *CustomerHandler) HandlePutCustomer(w http.ResponseWriter, r *http.Request) {
    log.Println("Updating an existing customer")
    r.ParseForm()

	id, err:= strconv.ParseInt(r.PathValue("id"), 10, 64);

	if err != nil {
		http.Error(w, "Customer id '%v' could not be parsed to an Int64", 400);
	}

	updatedCustomer:= db.Customer{Id: id, FirstName: r.PostForm.Get("firstName"), LastName: r.PostForm.Get("lastName"), Email: r.PostForm.Get("email") }

	h.db.Update(updatedCustomer)

	HtmxRedirect(w, "/customers", "#main-content")
}
