package handlers

import (
	"net/http"

	"github.com/TechnicallyMay/indigo/db"
)

type SettingsHandler struct {
	db db.SettingsTable
}

var settingsHandlerInstance *SettingsHandler

func NewSettingsHandler(db db.SettingsTable) *SettingsHandler {
	if settingsHandlerInstance == nil {
		settingsHandlerInstance = &SettingsHandler{db: db}
	}

	return settingsHandlerInstance
}

func (h SettingsHandler) HandleGetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.db.GetAll()
	handleHttpError(w, err, 500)

	renderTemplate(w, r, newRenderOpts("settings", settings, "content", "header"))
}

func (h SettingsHandler) HandlePutSettings(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	handleHttpError(w, err, 400)

	settings := db.IndigoSettings{
		BusinessName:  r.FormValue("businessName"),
		BusinessPhone: r.FormValue("businessPhone"),
		BusinessAddr:  r.FormValue("businessAddress"),
		BusinessCity:  r.FormValue("businessCity"),
		BusinessState: r.FormValue("businessState"),
		BusinessZip:   r.FormValue("businessZip"),
		InvoiceColor:  r.FormValue("invoiceColor"),
		InvoiceFooter: r.FormValue("invoiceFooter"),
		EmailBccs:     r.FormValue("emailBccs"),
	}

	err = h.db.SetAll(settings)
	handleHttpError(w, err, 500)

	renderTemplate(w, r, newRenderOpts("settings", settings, "content"))
}
