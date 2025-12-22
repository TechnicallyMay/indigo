package sender

import (
	"github.com/TechnicallyMay/indigo/db"
	"github.com/TechnicallyMay/indigo/mail"
	"github.com/TechnicallyMay/indigo/pdf"
)

type InvoiceSender struct {
	MailClient *mail.SmtpClient
}

func (s *InvoiceSender) SendInvoice(batch db.InvoiceBatch, cust db.Customer, items []db.InvoiceItem, allProducts map[int64]db.Product) error {
	pdf, err := pdf.MakeInvoicePdf(batch, cust, items, allProducts)
	if err != nil {
		return err
	}

	mail := mail.Mail{
		To:  []string{cust.Email},
		Bcc: []string{"masonwells01@pm.me"},

		From: "postmaster@sandbox6799805213174515aa97c14ef663c50e.mailgun.org", // TODO

		Subject: batch.NotificationSubject,
		Body:    batch.NotificationDescription,

		AttachmentFilePath: pdf,
		AttachmentFileName: "invoice.pdf",
	}

	err = s.MailClient.SendMail(mail)

	return err
}
