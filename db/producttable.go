package db

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

type productQuery struct {
}

func CreateProductQuery() productQuery {
	return productQuery{}
}

type Product struct {
	Id        int64
	Version   int64
	CreatedAt int64

	Name        string
	Description string
	UnitPrice   float64
}

type ProductTable struct {
	db *sql.DB
}

var productInstance *ProductTable

func InitProductTable(db *sql.DB) *ProductTable {
	log.Println("Initializing Product table.")
	if productInstance != nil {
		log.Println("Product table already initialized, returning existing instance.")
		return productInstance
	}

	log.Println("Ensuring Product table exists.")
	productInstance = &ProductTable{db: db}
	_, err := productInstance.db.Exec(`
        CREATE TABLE IF NOT EXISTS product(
            id INTEGER NOT NULL,
			version INTEGER NOT NULL,
			created_at INTEGER NOT NULL,

            name TEXT NOT NULL,
            description TEXT,
			unit_price REAL NOT NULL DEFAULT 0,
			PRIMARY KEY (id, version)
        );
    `)

	if err != nil {
		log.Fatal("Error while creating Product table.", err)
	}

	log.Println("Product table successfully initialized.")
	return productInstance
}

func (h *ProductTable) Get(id int64) (Product, error) {
	row := h.db.QueryRow(`
		SELECT id, version, created_at, name, description, unit_price
		FROM product
		WHERE id = (?)
		ORDER BY version DESC
		LIMIT 1;`, id)

	var product Product

	err := row.Scan(&product.Id, &product.Version, &product.CreatedAt, &product.Name, &product.Description, &product.UnitPrice)

	return product, err
}

func (h *ProductTable) QueryRow(query productQuery) (*Product, error) {
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

func (h *ProductTable) Query(query productQuery) ([]Product, error) {
	filters := make([]string, 0)
	params := make([]any, 0)

	// if query.BatchId != -1 {
	// 	filters = append(filters, "batch_id = (?)")
	// 	params = append(params, query.BatchId)
	// }

	queryStr := `
		SELECT id, version, created_at, name, description, unit_price
		FROM product
		INNER JOIN (
			SELECT id as innerId, MAX(version) as maxVersion
			FROM product as p2
			GROUP BY id
		) ON id = innerId AND version = maxVersion
	`

	if len(filters) > 0 {
		queryStr += " WHERE " + strings.Join(filters, " AND ")
	}

	rows, err := h.db.Query(queryStr, params...)
	if err != nil {
		return make([]Product, 0), err
	}
	return parseRows(rows)
}

func (h *ProductTable) List() ([]Product, error) {
	return h.Query(CreateProductQuery())
}

func (h *ProductTable) Add(product Product) (int64, error) {
	if product.Id != 0 {
		return 0, errors.New("Tried to add an product with an id")
	}

	res, err := h.db.Exec(`
		INSERT INTO product (id, version, created_at, name, description, unit_price) 
		VALUES (
			(SELECT COALESCE(MAX(id) + 1, 0) from product), 
			0, strftime('%s', 'now'), ?, ?, ?);`, product.Name, product.Description, product.UnitPrice)

	if err != nil {
		return 0, err
	}

	log.Println("Successfully added a new product.")
	return res.LastInsertId()
}

func (h *ProductTable) Update(product Product) error {
	log.Println("Updating existing product with id", product.Id)
	res, err := h.db.Exec(`
		INSERT INTO product (id, version, created_at, name, description, unit_price) 
		VALUES (?, (SELECT MAX(version) + 1 FROM product WHERE id == (?)), strftime('%s', 'now'), ?, ?, ?);`, product.Id, product.Id, product.Name, product.Description, product.UnitPrice)

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

	log.Println("Successfully updated product.")
	return nil
}

func parseRows(rows *sql.Rows) ([]Product, error) {
	products := make([]Product, 0)
	for rows.Next() {
		if rows.Err() != nil {
			return products, rows.Err()
		}

		var product Product
		if err := rows.Scan(&product.Id, &product.Version, &product.CreatedAt, &product.Name, &product.Description, &product.UnitPrice); err != nil {
			return products, err
		}

		products = append(products, product)
	}

	return products, nil
}
