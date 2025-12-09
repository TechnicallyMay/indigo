package handlers

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/TechnicallyMay/indigo/db"
)

type ProductHandler struct {
	db db.ProductTable
}

var productHandlerInstance *ProductHandler

func NewProductHandler(db db.ProductTable) *ProductHandler {
	if productHandlerInstance == nil {
		productHandlerInstance = &ProductHandler{db: db}
	}

	return productHandlerInstance
}

func (h *ProductHandler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.db.List()
	handleHttpError(w, err, 500)
	renderTemplate(w, r, newRenderOpts("products", products))
}

func (h *ProductHandler) HandleGetAddOrUpdateProductForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id != "" {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "Product id '%v' could not be parsed to an Int64", 400)
		}

		prod, err := h.db.Get(id)
		handleHttpError(w, err, 500)
		renderTemplate(w, r, newRenderOpts("addOrUpdateProduct", prod))
	} else {
		renderTemplate(w, r, newRenderOpts("addOrUpdateProduct", nil))
	}
}

func (h *ProductHandler) HandlePostProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding a new product")
	r.ParseForm()

	unitPrice, err := strconv.ParseFloat(r.PostForm.Get("unitPrice"), 64)
	handleHttpError(w, err, 400)
	newProduct := db.Product{Name: r.PostForm.Get("name"), Description: r.PostForm.Get("description"), UnitPrice: unitPrice}

	log.Printf("Got product %v %v %v\n", newProduct.Name, newProduct.Description, newProduct.UnitPrice)

	h.db.Add(newProduct)

	HtmxSoftRedirect(w, "/products", "#main-content")
}

func (h *ProductHandler) HandlePutProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating an existing product")
	r.ParseForm()

	updatedProduct, err := parseProductForm(r.PostForm)
	handleHttpError(w, err, 500)
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	handleHttpError(w, err, 500)

	updatedProduct.Id = id
	h.db.Update(*updatedProduct)

	HtmxSoftRedirect(w, "/products", "#main-content")
}

func parseProductForm(form url.Values) (*db.Product, error) {
	unitPrice, err := strconv.ParseFloat(form.Get("unitPrice"), 64)
	if err != nil {
		return nil, err
	}
	return &db.Product{Name: form.Get("name"), Description: form.Get("description"), UnitPrice: unitPrice}, nil
}
