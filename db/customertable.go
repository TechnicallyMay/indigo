package db

import (
	"database/sql"
	"log"
)

type Customer struct {
    Id int64
    FirstName string
    LastName string
    Email string
}

type CustomerTable struct {
    db *sql.DB
}

var instance *CustomerTable
func InitCustomerTable(db *sql.DB) *CustomerTable {
    log.Println("Initializing Customer table.")
    if instance != nil {
        log.Println("Customer table already initialized, returning existing instance.")
        return instance
    }

    log.Println("Ensuring customer table exists.")
    instance = &CustomerTable{db: db}
    _, err := instance.db.Exec(`
        CREATE TABLE IF NOT EXISTS customer(
            id INTEGER PRIMARY KEY,
            first_name TEXT NOT NULL,
            last_name TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE
        )
    `)

    if err != nil {
        log.Fatal("Error while creating customer table.", err)
    }

    log.Println("Customer table successfully initialized.")
    return instance
}

func (h *CustomerTable) Get(id int64) Customer {
    row:= h.db.QueryRow("SELECT id, first_name, last_name, email FROM customer WHERE id = (?)", id)

	var cust Customer
	if err:= row.Scan(&cust.Id, &cust.FirstName, &cust.LastName, &cust.Email); err != nil {
		log.Fatal("Error when getting customer.", err)
	}

    return cust
}

func (h *CustomerTable) List() []Customer {
    rows, err := h.db.Query("SELECT id, first_name, last_name, email FROM customer")
    if err != nil {
        log.Fatal("Error when listing customers.", err)
    }

    customers := make([]Customer, 0)
    for rows.Next() {
        var cust Customer
        if err := rows.Scan(&cust.Id, &cust.FirstName, &cust.LastName, &cust.Email); err != nil {
            log.Fatal("Error when listing customers.", err)
        }

        customers = append(customers, cust)
    }

    log.Printf("Found rows %v", customers)
    return customers
}

func (h *CustomerTable) Add(cust Customer) uint16 {
    log.Println("Adding new customer.")
    if cust.Id != 0 {
        log.Fatal("Tried to add an existing customer, updates are not yet supported.")
    }

    res, err := h.db.Exec(`INSERT INTO customer (first_name, last_name, email) VALUES (?, ?, ?)`,
        cust.FirstName, cust.LastName, cust.Email)

    if err != nil {
        log.Fatal("Error when inserting new customer.", err)
    }

    log.Println("Successfully added new customer.")
    newId, err := res.LastInsertId()

    if err != nil {
        log.Fatal("Error when retrieving id of new customer.", err)
    }

    return uint16(newId)
}

func (h *CustomerTable) Update(cust Customer) uint16 {
    log.Printf("Updating existing customer with id %v.\n", cust.Id)

    res, err := h.db.Exec(`UPDATE customer SET first_name = (?), last_name = (?), email = (?) WHERE customer.id = (?)`,
        cust.FirstName, cust.LastName, cust.Email, cust.Id)

    if err != nil {
        log.Fatal("Error when updating customer.", err)
    }

	rowsAffected, err:= res.RowsAffected()

    if err != nil {
        log.Fatal("Error when determining if update was successful.", err)
    }

	if rowsAffected == 0 {
        log.Fatal("Attempted update didn't modify any rows for id.", err)
	}

    log.Println("Successfully updated customer.")

    return uint16(cust.Id)
}
