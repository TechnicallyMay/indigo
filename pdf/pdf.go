package pdf

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"codeberg.org/go-pdf/fpdf"
	"github.com/TechnicallyMay/indigo/db"
	"github.com/dustin/go-humanize"
)

// A4 = 210 x 297 mm
var pageH float64 = 297
var pageW float64 = 210

var paddingY float64 = 5

func MakeInvoicePdf(settings db.IndigoSettings, batch db.InvoiceBatch, inv db.Invoice, cust db.Customer, items []db.InvoiceItem, allProducts map[int64]db.Product) (string, error) {
	grandTotal := 0.0
	for _, it := range items {
		grandTotal += float64(it.Quantity) * allProducts[it.ProductId].UnitPrice
	}
	totalStr := "$" + humanize.FormatFloat("#,###.##", grandTotal)

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetHeaderFunc(func() { header(pdf, batch, inv, settings, totalStr) })
	pdf.SetFooterFunc(func() { footer(pdf, settings) })
	pdf.AliasNbPages("")

	pdf.AddPage()
	invItems(pdf, items, allProducts, totalStr)

	filename := randFileName(".pdf")
	file, err := os.CreateTemp(os.TempDir(), filename)
	err = pdf.OutputAndClose(file)

	return file.Name(), err
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

func header(pdf *fpdf.Fpdf, batch db.InvoiceBatch, inv db.Invoice, settings db.IndigoSettings, grandTotal string) {
	r, g, b, err := hexToRgb(settings.InvoiceColor)
	if err != nil {
		pdf.SetFillColor(r, g, b)
	} else {
		pdf.SetFillColor(237, 229, 149)
	}
	pdf.Rect(0, 0, 1000, 40, "F")

	// Left header (business details)
	pdf.SetY(paddingY)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(40, 10, settings.BusinessName, "0", 2, "0", false, 0, "0")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 5, settings.BusinessPhone, "0", 2, "0", false, 0, "0")
	pdf.CellFormat(40, 5, settings.BusinessAddr, "0", 2, "0", false, 0, "0")
	pdf.CellFormat(40, 5, fmt.Sprintf("%s %s, %v", settings.BusinessCity, settings.BusinessState, settings.BusinessZip), "0", 2, "0", false, 0, "0")

	// Right header (invoice details)
	right := pageW * 0.5
	pdf.SetY(paddingY + 10)
	pdf.SetX(right)
	pdf.CellFormat(0, 5, "Invoice #:", "0", 0, "L", false, 0, "0")
	pdf.CellFormat(0, 5, string(inv.Id), "0", 1, "R", false, 0, "0")

	pdf.SetX(right)
	pdf.CellFormat(0, 5, "Payment Due:", "0", 0, "L", false, 0, "0")
	pdf.CellFormat(0, 5, batch.GetDueDateStr(), "0", 1, "R", false, 0, "0")

	pdf.SetX(right)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 5, "Invoice Total:", "0", 0, "L", false, 0, "0")
	pdf.CellFormat(0, 5, grandTotal, "0", 2, "R", false, 0, "0")

	pdf.SetY(50)
}

func invItems(pdf *fpdf.Fpdf, items []db.InvoiceItem, allProducts map[int64]db.Product, grandTotal string) {
	colWidths := []float64{pageW * .3, pageW * .2, pageW * .2, pageW * .2}
	cols := []string{"Item", "Unit Price", "Quantity", "Total"}
	alignment := []string{"L", "R", "R", "R"}

	pdf.SetFont("Arial", "B", 12)
	for i := range colWidths {
		pdf.CellFormat(colWidths[i], 10, cols[i], "TB", 0, alignment[i], false, 0, "0")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 12)
	for _, item := range items {
		product := allProducts[item.ProductId]

		for i := range colWidths {
			var content string

			switch cols[i] {
			case "Item":
				content = product.Name
			case "Unit Price":
				content = "$" + humanize.FormatFloat("#,###.##", product.UnitPrice)
			case "Quantity":
				content = fmt.Sprintf("x%v", item.Quantity)
			case "Total":
				total := product.UnitPrice * float64(item.Quantity)
				content = "$" + humanize.FormatFloat("#,###.##", total)
			}

			pdf.CellFormat(colWidths[i], 8, content, "", 0, alignment[i], false, 0, "0")
		}
		pdf.Ln(-1)
	}

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, grandTotal, "T", 0, "R", false, 0, "0")
}

func footer(pdf *fpdf.Fpdf, settings db.IndigoSettings) {
	pdf.SetY(-15)
	pdf.SetFont("Arial", "", 12)
	pdf.Write(5, settings.InvoiceFooter)
}

func hexToRgb(hex string) (r, g, b int, error error) {
	values, error := strconv.ParseInt(hex, 16, 64)

	if error != nil {
		return
	}

	r = int(values >> 16)
	g = int((values >> 8) & 0xFF)
	b = int(values & 0xFF)

	return
}
