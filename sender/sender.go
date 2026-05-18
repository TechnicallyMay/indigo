package sender

import (
	"context"
	"fmt"
	"log/slog"

	"strings"
	"time"

	"github.com/TechnicallyMay/indigo/appsettings"
	"github.com/TechnicallyMay/indigo/db"
	"github.com/TechnicallyMay/indigo/mail"
	"github.com/TechnicallyMay/indigo/pdf"

	"github.com/cenkalti/backoff/v5"
)

type InvoiceSender struct {
	MailClient *mail.SmtpClient
}

func (s *InvoiceSender) SendInvoice(settings db.IndigoSettings, smtpSettings appsettings.SmtpSettings, batch db.InvoiceBatch, inv db.Invoice, cust db.Customer, items []db.InvoiceItemWithProduct) error {
	slog.Info("Sending invoice", "customer first name", cust.FirstName, "customer last name", cust.LastName)
	pdf, err := pdf.MakeInvoicePdf(settings, batch, inv, cust, items)
	if err != nil {
		slog.Error("There was an error making the invoice pdf", "error", err)
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

		From: smtpSettings.From,

		Subject: batch.NotificationSubject,
		Body:    batch.NotificationDescription,

		AttachmentFilePath: pdf,
		AttachmentFileName: "invoice.pdf",
	}

	operation := func() (string, error) {
		err := s.MailClient.SendMail(mail)
		if err == nil {
			return "", nil
		}

		slog.Error("There was an error while sending the invoice", "error", err)

		if strings.Contains(err.Error(), "try again after ") {
			datePart := strings.Split(err.Error(), "try again after ")[1]
			nextTryAt, err := time.Parse(time.RFC1123, datePart)

			fmt.Println(datePart)

			if err != nil {
				return "", backoff.Permanent(err)
			}

			waitTime := time.Until(nextTryAt)
			slog.Info("Error was a throttling error, waiting to retry", "wait time secs", waitTime.Seconds())
			return "", backoff.RetryAfter(int(waitTime.Seconds()) + 2)
		}

		return "", backoff.Permanent(err)
	}

	_, err = backoff.Retry(context.TODO(), operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))

	if err == nil {
		slog.Info("Successfully sent invoice")
	} else {
		slog.Error("Unrecoverable error while sending invoice", "error", err)
	}

	return err
}
