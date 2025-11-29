# Indigo

## Todo
- [x] Figure out major data structures
- [ ] Billing Tab
    - [ ] Create new Button => Goes to `Billing/{id}`
    - [ ] `Billing/{id}` page
        - [ ] Customer list
        - [ ] Show each included customer
        - [ ] Edit each invoice
    - [ ] Billing page show historical sessions
* [ ] Send email to actual customers
* [ ] Config file
* [ ] Docker?
* [ ] nixpkg

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

    IsSent    bool
    SentAt    int64
    InvoiceId int64
    DueDate   int64 // Should be the same as the invoice batch, but keeping it flexible
}
```

```go
type InvoiceBatch struct {
    Id                int64
    CreatedAt         int64

    DueDate           int64 // Default to whatever invoice sender does (next 15th?), can add some fanciness later
    FinishedSendingAt int64 // When the last invoice notification in the batch was sent successfully
    AllNotificationsSent   bool
}
```
