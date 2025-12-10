package db

import (
	"database/sql"
	"log"
)

type InvoiceItem struct {
	InvoiceId      int64
	ProductId      int64
	ProductVersion int64
	Quantity       int64
}

type InvoiceItemTable struct {
	db *sql.DB
}

var invoiceItemInstance *InvoiceItemTable

func InitInvoiceItemTable(db *sql.DB) *InvoiceItemTable {
	log.Println("Initializing InvoiceItem table.")
	if invoiceItemInstance != nil {
		log.Println("InvoiceItem table already initialized, returning existing instance.")
		return invoiceItemInstance
	}

	log.Println("Ensuring InvoiceItem table exists.")
	invoiceItemInstance = &InvoiceItemTable{db: db}
	_, err := invoiceItemInstance.db.Exec(`
        CREATE TABLE IF NOT EXISTS invoice_item(
            invoice_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			product_version INTEGER NULL,
			quantity INTEGER NOT NULL DEFAULT 0,

			FOREIGN KEY(invoice_id) REFERENCES invoice(id)
			FOREIGN KEY(product_id) REFERENCES product(id)

			PRIMARY KEY (invoice_id, product_id)
        );
    `)

	if err != nil {
		log.Fatal("Error while creating InvoiceItem table.", err)
	}

	log.Println("InvoiceItem table successfully initialized.")
	return invoiceItemInstance
}

func (t *InvoiceItemTable) Add(item InvoiceItem) error {
	_, err := t.db.Exec(`
	    INSERT INTO invoice_item(invoice_id, product_id, quantity)
		VALUES ((?),(?),(?));`, item.InvoiceId, item.ProductId, item.Quantity)

	return err
}

func (t *InvoiceItemTable) Update(item InvoiceItem) error {
	_, err := t.db.Exec(`
	    UPDATE invoice_item SET quantity = (?)
		WHERE invoice_id = (?) AND product_id = (?)`, item.Quantity, item.ProductId)

	return err
}

// TODO: When finalizing an invoice, batch update the "product_version" to the latest
func (t *InvoiceItemTable) List(invoiceId int64) ([]InvoiceItem, error) {
	rows, err := t.db.Query(`
		SELECT invoice_id, product_id, product_version, quantity
		FROM invoice_item
		WHERE invoice_id = (?);`, invoiceId)

	items := make([]InvoiceItem, 0)
	if err != nil {
		return items, err
	}

	for rows.Next() {
		if rows.Err() != nil {
			return items, rows.Err()
		}

		var item InvoiceItem
		if err := rows.Scan(&item.InvoiceId, &item.ProductId, &item.ProductVersion, &item.Quantity); err != nil {
			return items, err
		}

		items = append(items, item)
	}

	return items, nil
}
