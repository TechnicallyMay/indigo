package sender

import (
	"strings"

	"github.com/TechnicallyMay/indigo/db"
	"github.com/TechnicallyMay/indigo/mail"
	"github.com/TechnicallyMay/indigo/pdf"
)

type InvoiceSender struct {
	MailClient *mail.SmtpClient
}

func (s *InvoiceSender) SendInvoice(settings db.IndigoSettings, batch db.InvoiceBatch, inv db.Invoice, cust db.Customer, items []db.InvoiceItemWithProduct) error {
	pdf, err := pdf.MakeInvoicePdf(settings, batch, inv, cust, items)
	if err != nil {
		return err
	}

	bccs := make([]string, 0)
	for bcc := range strings.SplitSeq(settings.EmailBccs, ",") {
		bcc = strings.TrimSpace(bcc)
		if bcc != "" {
			bccs = append(bccs, bcc)
		}
	}

	mail := mail.Mail{
		To:  []string{cust.Email},
		Bcc: bccs,

		From: "postmaster@sandbox6799805213174515aa97c14ef663c50e.mailgun.org", // TODO

		Subject: batch.NotificationSubject,
		Body:    batch.NotificationDescription,

		AttachmentFilePath: pdf,
		AttachmentFileName: "invoice.pdf",
	}

	err = s.MailClient.SendMail(mail)

	return err
}
