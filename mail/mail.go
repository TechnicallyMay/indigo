package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Mail struct {
	To                 []string
	Bcc                []string
	From               string
	Subject            string
	Body               string
	AttachmentFileName string
	AttachmentFilePath string
}

type SmtpClient struct {
	Host string
	Port int
	Auth smtp.Auth
}

func (client SmtpClient) SendMail(mail Mail) error {
	mailData := buildMailData(mail)
	// Cc works by including the addresses in the SendMail call but not in the message headers
	allRecipients := slices.Concat(mail.To, mail.Bcc)
	return smtp.SendMail(client.getFullAddress(), client.Auth, mail.From, allRecipients, mailData)
}

func readFile(fileName string) []byte {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func buildMailData(mail Mail) []byte {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "From: %s\r\n", mail.From)
	fmt.Fprintf(&buf, "To: %s\r\n", strings.Join(mail.To, ","))
	fmt.Fprintf(&buf, fmt.Sprintf("Subject: %s\r\n", mail.Subject))

	boundary := "doesthismatter"
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&buf, fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
	fmt.Fprintf(&buf, fmt.Sprintf("\r\n--%s\r\n", boundary))
	fmt.Fprintf(&buf, "Content-Type: text/plain;\r\n")
	fmt.Fprintf(&buf, fmt.Sprintf("\r\n%s", mail.Body))

	fmt.Fprintf(&buf, fmt.Sprintf("\r\n--%s\r\n", boundary))
	fmt.Fprintf(&buf, "Content-Type: text/plain;\r\n")
	fmt.Fprintf(&buf, "Content-Transfer-Encoding: base64\r\n")
	fmt.Fprintf(&buf, "Content-Disposition: attachment; filename="+mail.AttachmentFileName+"\r\n")
	fmt.Fprintf(&buf, "Content-ID: <"+mail.AttachmentFilePath+">\r\n\r\n")

	data := readFile(mail.AttachmentFilePath)
	b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(b, data)
	buf.Write(b)
	fmt.Fprintf(&buf, fmt.Sprintf("\r\n--%s", boundary))
	fmt.Fprintf(&buf, "--")

	return buf.Bytes()
}

func (client SmtpClient) getFullAddress() string {
	return client.Host + ":" + strconv.Itoa(client.Port)
}
