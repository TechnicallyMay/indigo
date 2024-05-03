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

func (h *CustomerTable) List() {
    rows, err := h.db.Query("SELECT id, first_name, last_name, email FROM customer")
    if err != nil {
        log.Fatal("Error when listing customers.", err)
    }

    customers := make([]Customer, 0)
    for rows.Next() {
        var cust Customer
        if err := rows.Scan(&cust.Id, &cust.FirstName, &cust.LastName, &cust.LastName); err != nil {
            log.Fatal("Error when listing customers.", err)
        }

        customers = append(customers, cust)
    }

    log.Printf("Found rows %v", customers)
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
