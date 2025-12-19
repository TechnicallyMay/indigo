package pdf

import (
	"fmt"
	"math/rand"

	"codeberg.org/go-pdf/fpdf"
)

func MakeInvoicePdf() (string, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "you owe me so much money")

	filename := randFileName(".pdf")
	err := pdf.OutputFileAndClose(filename)
	return filename, err
}

func randFileName(extension string) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := 24
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b) + extension
}
