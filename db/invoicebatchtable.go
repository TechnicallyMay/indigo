package db

import (
	"database/sql"
	"log"
	"time"
)

type InvoiceBatchState int

const (
	Draft InvoiceBatchState = iota
	Sending
	Sent
	Failed
	PartialFailure
)

func (state InvoiceBatchState) String() string {
	switch state {
	case Draft:
		return "Draft"
	case Sending:
		return "Sending"
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

	DueDate                 int64 // Default to whatever invoice sender does (next 15th?), can add some fanciness later
	NotificationSubject     string
	NotificationDescription string

	State             InvoiceBatchState
	FinishedSendingAt int64 // When the last invoice notification in the batch was sent successfully
}

func (b *InvoiceBatch) GetDueDateStr() string {
	dueDate := time.Unix(b.DueDate, 0)
	return dueDate.Format(time.DateOnly)
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

            due_date INTEGER NOT NULL,
			notification_subject TEXT,
			notification_description TEXT,

			state INTEGER NOT NULL DEFAULT 0,
			finished_sending_at INTEGER
        );
    `)

	if err != nil {
		log.Fatal("Error while creating Invoice Batch table.", err)
	}

	log.Println("InvoiceBatch table successfully initialized.")
	return invoiceBatchInstance
}

func (h *InvoiceBatchTable) Get(id int64) (InvoiceBatch, error) {
	row := h.db.QueryRow(`
		SELECT id, created_at, due_date, notification_subject, notification_description, state, finished_sending_at
		FROM invoice_batch 
		WHERE id = (?);`, id)

	var batch InvoiceBatch
	err := row.Scan(&batch.Id, &batch.CreatedAt, &batch.DueDate, &batch.NotificationSubject, &batch.NotificationDescription, &batch.State, &batch.FinishedSendingAt)

	return batch, err
}

func (h *InvoiceBatchTable) List() []InvoiceBatch {
	rows, err := h.db.Query(`
		SELECT id, created_at, due_date, notification_subject, notification_description, state, finished_sending_at
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
		if err := rows.Scan(&batch.Id, &batch.CreatedAt, &batch.DueDate, &batch.NotificationSubject, &batch.NotificationDescription, &batch.State, &batch.FinishedSendingAt); err != nil {
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
		INSERT INTO invoice_batch(created_at, due_date, notification_subject, notification_description, state, finished_sending_at) 
		VALUES (strftime('%s', 'now'), ?, ?, ?, ?, ?)`, batch.DueDate, batch.NotificationSubject, batch.NotificationDescription, batch.State, batch.FinishedSendingAt)

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
		UPDATE invoice_batch SET due_date = (?), notification_subject = (?), notification_description = (?), state = (?), finished_sending_at = (?)
		WHERE id = (?);`, batch.DueDate, batch.NotificationSubject, batch.NotificationDescription, batch.State, batch.FinishedSendingAt, batch.Id)

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

func (h *InvoiceBatchTable) GetInvoiceTotalsByCust(batchId int64) (map[int64]float64, error) {
	result := make(map[int64]float64, 0)
	rows, err := h.db.Query(`
		SELECT inv.customer_id, SUM(it.quantity * p.unit_price) FROM invoice AS inv
		INNER JOIN invoice_item AS it ON it.invoice_id = inv.id
		INNER JOIN product AS p ON it.product_id = p.id AND it.product_version = p.version
		WHERE inv.batch_id = (?)
		GROUP BY inv.id, inv.customer_id
	`, batchId)

	if err != nil {
		return result, err
	}

	for rows.Next() {
		if rows.Err() != nil {
			return result, rows.Err()
		}
		var custId int64
		var total float64
		if err := rows.Scan(&custId, &total); err != nil {
			return result, err
		}

		result[custId] = total
	}

	return result, nil
}
