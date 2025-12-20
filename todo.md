# Indigo

## Todo
### In Progress
- [x] Save invoice batch metadata
- [x] Join billing page into 1 page, edit vs not-edit
- [ ] oob button rendering when not Htmx-request
- [ ] Invoice notifications
- [ ] Send email to actual customers
### Short List

### Must
- [ ] When invoices are sent, need to finalize the invoice item version
- [ ] Deleted flags for versioned records
- [ ] Product dropdown properly re-filter after adding item
- [ ] Add all customers to batch in a single action
- [ ] Ability to delete draft batches
- [ ] Don't create a new batch every time you navigate to `/new`
  - Maybe just check if there already exists an empty batch
* [ ] Config file

### Should
- [ ] Remove weird save button from invocie batch
- [ ] Save when starting an invoice batch?
- [ ] Customers Tweaks
  - [ ] I think I need a "enabled" flag
  - [ ] Delete button (soft delete obvs)
- [ ] Fix content scrolling main screen. Limit widgets to bottom padding, scroll widget
- [ ] Delete an invoice batch (if in draft status)
- [ ] Account for timezone of browser
* [ ] nixpkg
- [ ] Cleanup
  - [ ] Error handling
  - [ ] Rename things to make more sense (handlers)
  - [ ] println vs log

### Could
- [ ] Use pico modal for confirmation
* [ ] Docker?

### Archive
- [x] Select Due Date
- [x] Write invoice subject & body
- [x] Add a route for rendering viewInvoiceItem
- [x] Move the "Add a new Customer" button on the invoice page to the header
- [x] Some CSS container for scrolling the main content while still having footer
  - [x] Maybe need to split into body/footer instead of just content template
     - [x] Would need to define in each page
- [x] Billing Tab
    - [x] `Billing/{id}` page
        - [x] Customer list
        - [x] Show each included customer
        - [x] Edit each invoice
        - [x] Remove invoice
    - [x] Billing page list existing sessions
    - [x] Properly redirect from billing/new to billing/{new-id}
    - [x] Do invoice batches need a state?
      - Maybe invoice notifications also have a state. If they error, we have to create a new notification. Error column added to notification
    - [x] Add states for invoice batches
    - [x] Billing database handler
    - [x] Create new Button => Goes to `Billing/{id}`

- [x] Figure out major data structures
- [x] Invoice Table

## Notes
### Invoice Batch States
The point would be to better indicate the cumulative state of the invoices belonging to the batch. If all were sent successfully, the batch is successful. If only some were successfully sent, it would be some kind of partial state.

* Draft
* Sent
* Failed
* PartialFailure

## Design

### Billing
#### Session
* Creating a session starts empty
* When a customer is added, a new invoice is created for them
* Customer list shows a clear indicator of how many registered customers **are not** included
* Copy other session
  * Copies all customers from selected session.
  * Creates a new invoice for each customer
  * Copies data from previous invoice to new
* Session page shows a button "copy most recent"
* History page has a "copy" button for each past sent session
* Sending one invoice is still a session
* Session controls due date in a relative fashion to enable copying
     * Big indicator for what the due date for the current session is
     * Options
         * x date of next month
         * x date of this month
         * Coming x date
         * in x days

#### Data Structures

```go
type Customer struct {
	Id        int64
  Version   int64 // Updates just insert new records. Primary key = id + version. On GET, return latest version.
  CreatedAt int64

	FirstName string
	LastName  string
	Email     string
}
```

```go
type Product struct {
	Id          int64
    Version     int64
    CreatedAt   int64

	Name        string
    Description string
    UnitCost    float32
}
```

```go
type Invoice struct {
    Id              int64
    CreatedAt       int64

    BatchId         int64
    CustomerId      int64
    CustomerVersion int64 // Useful for records, like seeing which email was active at the time
    IsPaid          bool

    // Invoice doesn't have concept of due date. That's a property of the notification.
}
```

```go
type InvoiceItem struct {
	InvoiceId      int64
	ProductId      int64
  ProductVersion int64
  Quantity       int64
}
```

```go
type InvoiceNotification struct {
    Id        int64
    CreatedAt int64

    State     State
    SentAt    int64
    InvoiceId int64
    DueDate   int64 // Should be the same as the invoice batch, but keeping it flexible
}
```

```go
tpe InvoiceBatch struct {
    Id                     int64
    CreatedAt              int64

    DueDate                int64 // Default to whatever invoice sender does (next 15th?), can add some fanciness later
    FinishedSendingAt      int64 // When the last invoice notification in the batch was sent successfully
    AllNotificationsSent   bool
}
```
