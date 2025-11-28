package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

func readFile(fileName string) []byte {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func buildMailData(mail Mail) []byte {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("From: %s\r\n", mail.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ",")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.Subject))

	boundary := "doesthismatter"
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/plain;\r\n")
	buf.WriteString(fmt.Sprintf("\r\n%s", mail.Body))

	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/plain;\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
	buf.WriteString("Content-Disposition: attachment; filename=" + mail.AttachmentFilePath + "\r\n")
	buf.WriteString("Content-ID: <" + mail.AttachmentFilePath + ">\r\n\r\n")

	data := readFile(mail.AttachmentFilePath)
	b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(b, data)
	buf.Write(b)
	buf.WriteString(fmt.Sprintf("\r\n--%s", boundary))
	buf.WriteString("--")

	return buf.Bytes()
}

type Mail struct {
	To                 []string
	From               string
	Subject            string
	Body               string
	AttachmentFilePath string
}

type SmtpClient struct {
	Host string
	Port int
	Auth smtp.Auth
}

func (client SmtpClient) SendMail(mail Mail) error {
	var mailData = buildMailData(mail)
	return smtp.SendMail(client.getFullAddress(), client.Auth, mail.From, mail.To, mailData)
}

func (client SmtpClient) getFullAddress() string {
	return client.Host + ":" + strconv.Itoa(client.Port)
}
