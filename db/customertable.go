package db

import (
	"database/sql"
	"log"
)

type Customer struct {
	Id        int64
	Version   int64
	CreatedAt int64

	FirstName string
	LastName  string
	Email     string
}

type CustomerTable struct {
	db *sql.DB
}

var customerTableInstance *CustomerTable

func InitCustomerTable(db *sql.DB) *CustomerTable {
	log.Println("Initializing Customer table.")
	if customerTableInstance != nil {
		log.Println("Customer table already initialized, returning existing instance.")
		return customerTableInstance
	}

	log.Println("Ensuring customer table exists.")
	customerTableInstance = &CustomerTable{db: db}
	_, err := customerTableInstance.db.Exec(`
        CREATE TABLE IF NOT EXISTS customer(
            id INTEGER NOT NULL,
			version INTEGER NOT NULL,
			created_at INTEGER NOT NULL,

            first_name TEXT NOT NULL,
            last_name TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE,
			PRIMARY KEY (id, version)
        );
    `)

	if err != nil {
		log.Fatal("Error while creating customer table.", err)
	}

	log.Println("Customer table successfully initialized.")
	return customerTableInstance
}

func (h *CustomerTable) Get(id int64) Customer {
	row := h.db.QueryRow(`
		SELECT id, version, created_at, first_name, last_name, email 
		FROM customer 
		WHERE id = (?)
		ORDER BY version DESC
		LIMIT 1;`, id)

	var cust Customer
	if err := row.Scan(&cust.Id, &cust.Version, &cust.CreatedAt, &cust.FirstName, &cust.LastName, &cust.Email); err != nil {
		log.Fatal("Error when getting customer.", err)
	}

	return cust
}

func (h *CustomerTable) GetByInvoiceBatch(bId int64) []Customer {
	rows, err := h.db.Query(`
		SELECT customer.id, customer.version, customer.created_at, customer.first_name, customer.last_name, customer.email
		FROM customer
		INNER JOIN (
			SELECT id as innerId, MAX(version) as maxVersion
			FROM customer
			GROUP BY id
		) ON customer.id = innerId
		INNER JOIN invoice ON customer.id = customer_id
		WHERE customer.version = maxVersion AND invoice.batch_id = (?);`, bId)

	if err != nil {
		log.Fatal("Error when listing customers.", err)
	}

	return parseCustomerRows(rows)
}

func (h *CustomerTable) List() []Customer {
	rows, err := h.db.Query(`
		SELECT id, version, created_at, first_name, last_name, email
		FROM customer
		INNER JOIN (
			SELECT id as innerId, MAX(version) as maxVersion
			FROM customer
			GROUP BY id
		) ON id = innerId
		WHERE version = maxVersion;
		`)

	if err != nil {
		log.Fatal("Error when listing customers.", err)
	}

	return parseCustomerRows(rows)
}

func (h *CustomerTable) Add(cust Customer) int64 {
	log.Println("Adding new customer.")
	if cust.Id != 0 {
		log.Fatal("Tried to add an existing customer, updates are not yet supported.")
	}

	res, err := h.db.Exec(`
		INSERT INTO customer (id, version, created_at, first_name, last_name, email) 
		VALUES (
			(SELECT COALESCE(MAX(id) + 1, 0) from customer), 
			0, strftime('%s', 'now'), ?, ?, ?);`,
		cust.FirstName, cust.LastName, cust.Email)

	if err != nil {
		log.Fatal("Error when inserting new customer. ", err)
	}

	log.Println("Successfully added new customer.")
	newId, err := res.LastInsertId()

	if err != nil {
		log.Fatal("Error when retrieving id of new customer.", err)
	}

	return newId
}

func (h *CustomerTable) Update(cust Customer) int64 {
	log.Printf("Updating existing customer with id %v.\n", cust.Id)

	res, err := h.db.Exec(`
		INSERT INTO customer (id, version, created_at, first_name, last_name, email) 
		VALUES (?, (SELECT MAX(version) + 1 FROM customer WHERE id == (?)), strftime('%s', 'now'), ?, ?, ?);`, cust.Id, cust.Id, cust.FirstName, cust.LastName, cust.Email)

	if err != nil {
		log.Fatal("Error when updating customer.", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatal("Error when determining if update was successful.", err)
	}

	if rowsAffected == 0 {
		log.Fatal("Attempted update didn't modify any rows for id.", err)
	}

	log.Println("Successfully updated customer.")

	return cust.Id
}

func parseCustomerRows(rows *sql.Rows) []Customer {
	customers := make([]Customer, 0)
	for rows.Next() {
		if rows.Err() != nil {
			log.Fatal("Error when listing customers.", rows.Err())
		}
		var cust Customer
		if err := rows.Scan(&cust.Id, &cust.Version, &cust.CreatedAt, &cust.FirstName, &cust.LastName, &cust.Email); err != nil {
			log.Fatal("Error when listing customers.", err)
		}

		customers = append(customers, cust)
	}

	return customers
}
