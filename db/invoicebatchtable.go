package db

import (
	"database/sql"
	"log"
)

type InvoiceBatchState int

const (
	Draft InvoiceBatchState = iota
	Sent
	Failed
	PartialFailure
)

func (state InvoiceBatchState) String() string {
	switch state {
	case Draft:
		return "Draft"
	case Sent:
		return "Sent"
	case Failed:
		return "Failed"
	case PartialFailure:
		return "PartialFailure"
	default:
		return "Unknown"
	}
}

type InvoiceBatch struct {
	Id        int64
	CreatedAt int64

	State             InvoiceBatchState
	DueDate           int64 // Default to whatever invoice sender does (next 15th?), can add some fanciness later
	FinishedSendingAt int64 // When the last invoice notification in the batch was sent successfully
}

type InvoiceBatchTable struct {
	db *sql.DB
}

var invoiceBatchInstance *InvoiceBatchTable

func InitInvoiceBatchTable(db *sql.DB) *InvoiceBatchTable {
	log.Println("Initializing InvoiceBatch table.")
	if invoiceBatchInstance != nil {
		log.Println("InvoiceBatch table already initialized, returning existing instance.")
		return invoiceBatchInstance
	}

	log.Println("Ensuring Invoice Batch table exists.")
	invoiceBatchInstance = &InvoiceBatchTable{db: db}
	_, err := invoiceBatchInstance.db.Exec(`
        CREATE TABLE IF NOT EXISTS invoice_batch(
            id INTEGER PRIMARY KEY,
			created_at INTEGER NOT NULL,

			state INTEGER NOT NULL DEFAULT 0,
            due_date INTEGER NOT NULL,
			finished_sending_at INTEGER
        );
    `)

	if err != nil {
		log.Fatal("Error while creating Invoice Batch table.", err)
	}

	log.Println("InvoiceBatch table successfully initialized.")
	return invoiceBatchInstance
}

func (h *InvoiceBatchTable) Get(id int64) InvoiceBatch {
	row := h.db.QueryRow(`
		SELECT id, created_at, state, due_date, finished_sending_at
		FROM invoice_batch 
		WHERE id = (?);`, id)

	var batch InvoiceBatch
	if err := row.Scan(&batch.Id, &batch.CreatedAt, &batch.State, &batch.DueDate, &batch.FinishedSendingAt); err != nil {
		log.Fatal("Error when getting invoice batch.", err)
	}

	return batch
}

func (h *InvoiceBatchTable) List() []InvoiceBatch {
	rows, err := h.db.Query(`
		SELECT id, created_at, state, due_date, finished_sending_at 
		FROM invoice_batch;`)
	if err != nil {
		log.Fatal("Error when listing invoice batches.", err)
	}

	batches := make([]InvoiceBatch, 0)
	for rows.Next() {
		if rows.Err() != nil {
			log.Fatal("Error when listing invoice batches.", err)
		}
		var batch InvoiceBatch
		if err := rows.Scan(&batch.Id, &batch.CreatedAt, &batch.State, &batch.DueDate, &batch.FinishedSendingAt); err != nil {
			log.Fatal("Error when listing invoice batches.", err)
		}

		batches = append(batches, batch)
	}

	return batches
}

func (h *InvoiceBatchTable) Add(batch InvoiceBatch) int64 {
	log.Println("Adding new invoice batch.")
	if batch.Id != 0 {
		log.Fatal("Tried to add an existing invoice batch, updates are not yet supported.")
	}

	res, err := h.db.Exec(`
		INSERT INTO invoice_batch(created_at, state, due_date, finished_sending_at) 
		VALUES (strftime('%s', 'now'), ?, ?, ?)`, batch.State, batch.DueDate, batch.FinishedSendingAt)

	if err != nil {
		log.Fatal("Error when inserting new invoice batch. ", err)
	}

	log.Println("Successfully added new invoice batch.")
	newId, err := res.LastInsertId()

	if err != nil {
		log.Fatal("Error when retrieving id of new invoice batch.", err)
	}

	return newId
}

func (h *InvoiceBatchTable) Update(batch InvoiceBatch) int64 {
	log.Printf("Updating existing invoice batch with id %v.\n", batch.Id)

	res, err := h.db.Exec(`
		UPDATE invoice_batch SET created_at = (?), state = (?), due_date = (?), finished_sending_at = (?)
		WHERE id = (?);`, batch.Id, batch.CreatedAt, batch.State, batch.DueDate, batch.FinishedSendingAt)

	if err != nil {
		log.Fatal("Error when updating invoice batch.", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatal("Error when determining if update was successful.", err)
	}

	if rowsAffected == 0 {
		log.Fatal("Attempted update didn't modify any rows for id.", err)
	}

	log.Println("Successfully updated invoice batch.")

	return batch.Id
}
