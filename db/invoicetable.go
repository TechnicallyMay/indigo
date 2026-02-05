package db

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

type invoiceQuery struct {
	BatchId         int64
	CustomerId      int64
	CustomerVersion int64
	IsPaid          int64
}

func CreateInvoiceQuery() invoiceQuery {
	return invoiceQuery{
		BatchId:         -1,
		CustomerId:      -1,
		CustomerVersion: -1,
		IsPaid:          -1,
	}
}

type Invoice struct {
	Id        int64
	CreatedAt int64

	BatchId         int64
	CustomerId      int64
	CustomerVersion int64 // Useful for records, like seeing which email was active at the time
	IsPaid          bool

	// Invoice doesn't have concept of due date. That's a property of the notification.
}

type InvoiceTable struct {
	db *sql.DB
}

var invoiceInstance *InvoiceTable

func InitInvoiceTable(db *sql.DB) *InvoiceTable {
	log.Println("Initializing Invoice table.")
	if invoiceInstance != nil {
		log.Println("Invoice table already initialized, returning existing instance.")
		return invoiceInstance
	}

	log.Println("Ensuring Invoice table exists.")
	invoiceInstance = &InvoiceTable{db: db}
	_, err := invoiceInstance.db.Exec(`
        CREATE TABLE IF NOT EXISTS invoice(
            id INTEGER PRIMARY KEY,
			created_at INTEGER NOT NULL,

			batch_id INTEGER NOT NULL,
			customer_id INTEGER NOT NULL,
			customer_version INTEGER NOT NULL,

			is_paid INTEGER NOT NULL DEFAULT 0,

			FOREIGN KEY(batch_id) REFERENCES invoice_batch(id)
			FOREIGN KEY(customer_id) REFERENCES customer(id)
        );
    `)

	if err != nil {
		log.Fatal("Error while creating Invoice table.", err)
	}

	log.Println("Invoice table successfully initialized.")
	return invoiceInstance
}

func (h *InvoiceTable) Get(id int64) (Invoice, error) {
	row := h.db.QueryRow(`
		SELECT id, created_at, batch_id, customer_id, customer_version, is_paid
		FROM invoice
		WHERE id = (?);`, id)

	var invoice Invoice

	err := row.Scan(&invoice.Id, &invoice.CreatedAt, &invoice.BatchId, &invoice.CustomerId, &invoice.CustomerVersion, &invoice.IsPaid)

	return invoice, err
}

func (h *InvoiceTable) QueryRow(query invoiceQuery) (*Invoice, error) {
	rows, err := h.Query(query)
	if err != nil {
		return nil, err
	}

	if len(rows) > 1 {
		return nil, errors.New("Found more than one row while querying for exactly one")
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return &rows[0], nil

}

func (h *InvoiceTable) Query(query invoiceQuery) ([]Invoice, error) {
	filters := make([]string, 0)
	params := make([]any, 0)

	if query.BatchId != -1 {
		filters = append(filters, "batch_id = (?)")
		params = append(params, query.BatchId)
	}

	if query.CustomerId != -1 {
		filters = append(filters, "customer_id = (?)")
		params = append(params, query.CustomerId)
	}

	if query.CustomerVersion != -1 {
		filters = append(filters, "customer_version = (?)")
		params = append(params, query.CustomerVersion)
	}

	if query.IsPaid != -1 {
		filters = append(filters, "is_paid = (?)")
		params = append(params, query.CustomerVersion)
	}

	queryStr := `
		SELECT id, created_at, batch_id, customer_id, customer_version, is_paid
		FROM invoice
	`

	if len(filters) > 0 {
		queryStr += " WHERE " + strings.Join(filters, " AND ")
	}

	rows, err := h.db.Query(queryStr, params...)
	if err != nil {
		return make([]Invoice, 0), err
	}
	return parseInvoiceRows(rows)
}

func (h *InvoiceTable) List() ([]Invoice, error) {
	return h.Query(CreateInvoiceQuery())
}

func (h *InvoiceTable) Add(invoice Invoice) (int64, error) {
	if invoice.Id != 0 {
		return 0, errors.New("Tried to add an invoice with an id")
	}

	res, err := h.db.Exec(`
		INSERT INTO invoice(created_at, batch_id, customer_id, customer_version, is_paid)
		VALUES (strftime('%s', 'now'), ?, ?, ?, ?)`, invoice.BatchId, invoice.CustomerId, invoice.CustomerVersion, invoice.IsPaid)

	if err != nil {
		return 0, err
	}

	log.Println("Successfully added a new invoice.")
	newId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (h *InvoiceTable) Update(invoice Invoice) error {
	res, err := h.db.Exec(`
		UPDATE invoice SET created_at = (?), batch_id = (?), customer_id = (?), customer_version = (?), is_paid = (?)
		WHERE id = (?)`, invoice.BatchId, invoice.CustomerId, invoice.CustomerVersion, invoice.IsPaid)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatal("Error when determining if update was successful.", err)
	}

	if rowsAffected == 0 {
		log.Fatal("Attempted update didn't modify any rows for id.", err)
	}

	log.Println("Successfully updated invoice.")
	return nil
}

func (h InvoiceTable) Delete(id int64) error {
	// TODO: Soft delete needed if a notification has already been sent
	// TODO: For now don't allow deleting if state is not draft
	res, err := h.db.Exec(`
		BEGIN TRANSACTION;

		DELETE FROM invoice_item
		WHERE invoice_id = (?);

		DELETE FROM invoice
		WHERE id = (?);

		COMMIT;`, id, id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Attempted update didn't modify any rows for id.")
	}

	return nil
}

func parseInvoiceRows(rows *sql.Rows) ([]Invoice, error) {
	invoices := make([]Invoice, 0)
	for rows.Next() {
		if rows.Err() != nil {
			return invoices, rows.Err()
		}

		var invoice Invoice
		if err := rows.Scan(&invoice.Id, &invoice.CreatedAt, &invoice.BatchId, &invoice.CustomerId, &invoice.CustomerVersion, &invoice.IsPaid); err != nil {
			return invoices, err
		}

		invoices = append(invoices, invoice)
	}

	return invoices, nil
}
