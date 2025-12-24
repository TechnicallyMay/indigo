package db

import (
	"database/sql"
	"errors"
	"log"
)

type InvoiceItem struct {
	InvoiceId      int64
	ProductId      int64
	ProductVersion *int64
	Quantity       int64
}

type InvoiceItemWithProduct struct {
	Item    InvoiceItem
	Product Product
}

func (i *InvoiceItemWithProduct) GetItemSubtotal() float64 {
	return float64(i.Item.Quantity) * i.Product.UnitPrice
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
			product_version INTEGER NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 0,

			FOREIGN KEY(invoice_id) REFERENCES invoice(id)
			FOREIGN KEY(product_id, product_version) REFERENCES product(id, version)

			PRIMARY KEY (invoice_id, product_id)
        );
    `)

	if err != nil {
		log.Fatal("Error while creating InvoiceItem table.", err)
	}

	log.Println("InvoiceItem table successfully initialized.")
	return invoiceItemInstance
}

func (t *InvoiceItemTable) Get(invId int64, prodId int64) (*InvoiceItemWithProduct, error) {
	row := t.db.QueryRow(`
		SELECT it.invoice_id, it.product_id, it.product_version, it.quantity, p.id, p.version, p.created_at, p.name, p.description, p.unit_price
		FROM invoice_item as it
		INNER JOIN product as p ON it.product_id = p.id AND it.product_version = p.version
		WHERE it.invoice_id = (?) AND it.product_id = (?);`, invId, prodId)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var it InvoiceItem
	var pr Product
	if err := row.Scan(&it.InvoiceId, &it.ProductId, &it.ProductVersion, &it.Quantity, &pr.Id, &pr.Version, &pr.CreatedAt, &pr.Name, &pr.Description, &pr.UnitPrice); err != nil {
		return nil, err
	}

	return &InvoiceItemWithProduct{Item: it, Product: pr}, nil
}

func (t *InvoiceItemTable) Add(item InvoiceItem) error {
	_, err := t.db.Exec(`
	    INSERT INTO invoice_item(invoice_id, product_id, product_version, quantity)
		VALUES ((?),(?),(?),(?));`, item.InvoiceId, item.ProductId, item.ProductVersion, item.Quantity)

	return err
}

func (t *InvoiceItemTable) Update(item InvoiceItem) error {
	_, err := t.db.Exec(`
	    UPDATE invoice_item SET quantity = (?)
		WHERE invoice_id = (?) AND product_id = (?)`, item.Quantity, item.InvoiceId, item.ProductId)

	return err
}

func (t *InvoiceItemTable) List(invoiceId int64) ([]InvoiceItemWithProduct, error) {
	rows, err := t.db.Query(`
		SELECT it.invoice_id, it.product_id, it.product_version, it.quantity, p.id, p.version, p.created_at, p.name, p.description, p.unit_price
		FROM invoice_item as it
		INNER JOIN product as p ON it.product_id = p.id AND it.product_version = p.version
		WHERE it.invoice_id = (?);`, invoiceId)

	items := make([]InvoiceItemWithProduct, 0)
	if err != nil {
		return items, err
	}

	for rows.Next() {
		if rows.Err() != nil {
			return items, rows.Err()
		}

		var it InvoiceItem
		var pr Product
		if err := rows.Scan(&it.InvoiceId, &it.ProductId, &it.ProductVersion, &it.Quantity, &pr.Id, &pr.Version, &pr.CreatedAt, &pr.Name, &pr.Description, &pr.UnitPrice); err != nil {
			return items, err
		}

		items = append(items, InvoiceItemWithProduct{Item: it, Product: pr})
	}

	return items, nil
}

func (t InvoiceItemTable) Delete(invId int64, prodId int64) error {
	res, err := t.db.Exec(`
		DELETE FROM invoice_item
		WHERE invoice_id = (?)
		AND product_id = (?);`, invId, prodId)

	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if cnt != 1 {
		return errors.New("When deleting an invoice item, " + string(cnt) + " rows were effected. 1 was expected")
	}

	return nil
}
