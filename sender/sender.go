package sender

import (
	"fmt"

	"github.com/TechnicallyMay/indigo/db"
	"github.com/TechnicallyMay/indigo/mail"
)

type InvoiceSender struct {
	MailClient *mail.SmtpClient
}

func (s *InvoiceSender) SendInvoice(batch db.InvoiceBatch, cust db.Customer, items []db.InvoiceItem, allProducts map[int64]db.Product) error {
	fmt.Println("-----------------")
	fmt.Println("Sending invoice to", cust.FirstName, "@", cust.Email)
	fmt.Println("Due Date", batch.DueDate)
	fmt.Println("Items")
	for _, item := range items {
		fmt.Println("\t", allProducts[item.ProductId].Name, "x"+fmt.Sprint(item.Quantity))
	}
	fmt.Println()

	// pdf.MakeInvoicePdf()
	return nil
}
