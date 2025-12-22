package db

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

type invoiceNotificationQuery struct {
	InvoiceId int64
	// BatchId int64 // TODO
}

func CreateInvoiceNotificationQuery() invoiceNotificationQuery {
	return invoiceNotificationQuery{
		InvoiceId: -1,
	}
}

type InvoiceNotification struct {
	Id     int64
	SentAt int64

	Successful bool
	Error      string

	InvoiceId int64
	DueDate   int64
}

type InvoiceNotificationTable struct {
	db *sql.DB
}

var invoiceNotificationInstance *InvoiceNotificationTable

func InitInvoiceNotificationTable(db *sql.DB) *InvoiceNotificationTable {
	log.Println("Initializing InvoiceNotification table.")
	if invoiceNotificationInstance != nil {
		log.Println("InvoiceNotification table already initialized, returning existing instance.")
		return invoiceNotificationInstance
	}

	log.Println("Ensuring InvoiceNotification table exists.")
	invoiceNotificationInstance = &InvoiceNotificationTable{db: db}
	_, err := invoiceNotificationInstance.db.Exec(`
        CREATE TABLE IF NOT EXISTS invoice_notification(
			id INTEGER PRIMARY KEY,
			invoice_id INTEGER NOT NULL,

			sent_at INTEGER NOT NULL,
			due_date INTEGER NOT NULL,

			successful INTEGER NOT NULL,
			error TEXT NULL,

			FOREIGN KEY(invoice_id) REFERENCES invoice(id)
        );
    `)

	if err != nil {
		log.Fatal("Error while creating InvoiceNotification table.", err)
	}

	log.Println("InvoiceNotification table successfully initialized.")
	return invoiceNotificationInstance
}

func (h *InvoiceNotificationTable) Query(query invoiceNotificationQuery) ([]InvoiceNotification, error) {
	filters := make([]string, 0)
	params := make([]any, 0)

	if query.InvoiceId != -1 {
		filters = append(filters, "invoice_id = (?)")
		params = append(params, query.InvoiceId)
	}

	queryStr := `
		SELECT id, invoice_id, sent_at, due_date, successful, error
		FROM invoice_notification
	`

	if len(filters) > 0 {
		queryStr += " WHERE " + strings.Join(filters, " AND ")
	}

	rows, err := h.db.Query(queryStr, params...)
	if err != nil {
		return make([]InvoiceNotification, 0), err
	}
	return parseNotificationRows(rows)
}

func (h *InvoiceNotificationTable) Add(not InvoiceNotification) (int64, error) {
	if not.Id != 0 {
		return 0, errors.New("Tried to add an invoiceNotification with an id")
	}

	res, err := h.db.Exec(`
		INSERT INTO invoice_notification(invoice_id, sent_at, due_date, successful, error)
		VALUES ((?), (?), (?), (?), (?));`, not.InvoiceId, not.SentAt, not.DueDate, not.Successful, not.Error)

	if err != nil {
		return 0, err
	}

	log.Println("Successfully added a new invoiceNotification.")
	return res.LastInsertId()
}

func parseNotificationRows(rows *sql.Rows) ([]InvoiceNotification, error) {
	nots := make([]InvoiceNotification, 0)
	for rows.Next() {
		if rows.Err() != nil {
			return nots, rows.Err()
		}

		var not InvoiceNotification
		if err := rows.Scan(&not.Id, &not.InvoiceId, &not.SentAt, &not.DueDate, &not.Successful, &not.Error); err != nil {
			return nots, err
		}

		nots = append(nots, not)
	}

	return nots, nil
}
