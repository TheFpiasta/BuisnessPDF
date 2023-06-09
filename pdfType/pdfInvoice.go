package pdfType

import (
	"SimpleInvoice/generator"
	"encoding/json"
	"fmt"
	errorsWithStack "github.com/go-errors/errors"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"strconv"
)

type Invoice struct {
	data             invoiceRequestData
	meta             PdfMeta
	logger           *zerolog.Logger
	printErrStack    bool
	pdfGen           *generator.PDFGenerator
	defaultLineColor generator.Color
}

type invoiceRequestData struct {
	SenderAddress   FullPersonInfo `json:"senderAddress"`
	ReceiverAddress FullPersonInfo `json:"receiverAddress"`
	SenderInfo      SenderInfo     `json:"senderInfo"`
	InvoiceMeta     struct {
		InvoiceNumber  string `json:"invoiceNumber"`
		InvoiceDate    string `json:"invoiceDate"`
		CustomerNumber string `json:"customerNumber"`
	} `json:"invoiceMeta"`
	InvoiceBody struct {
		OpeningText     string `json:"openingText"`
		ServiceTimeText string `json:"serviceTimeText"`
		HeadlineText    string `json:"headlineText"`
		ClosingText     string `json:"closingText"`
		UstNotice       string `json:"ustNotice"`
		InvoicedItems   []struct {
			PositionNumber string  `json:"positionNumber"`
			Quantity       float64 `json:"quantity"`
			Unit           string  `json:"unit"`
			Description    string  `json:"description"`
			SinglePrice    int     `json:"singlePrice"`
			Currency       string  `json:"currency"`
			TaxRate        int     `json:"taxRate"`
		} `json:"invoicedItems"`
	} `json:"invoiceBody"`
}

func NewInvoice(logger *zerolog.Logger) *Invoice {
	return &Invoice{
		data: invoiceRequestData{},
		meta: PdfMeta{
			Margin: pdfMargin{
				Left:   25,
				Right:  20,
				Top:    45,
				Bottom: 0,
			},
			Font: pdfFont{
				FontName:    "openSans",
				SizeDefault: 10,
				SizeSmall:   8,
				SizeLarge:   15,
			},
		},
		logger:           logger,
		printErrStack:    logger.GetLevel() <= zerolog.DebugLevel,
		defaultLineColor: generator.Color{R: 200, G: 200, B: 200},
	}
}

func (i *Invoice) SetDataFromRequest(request *http.Request) (err error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			i.LogError(err)
		}
	}(request.Body)

	err = json.NewDecoder(request.Body).Decode(&i.data)
	if err != nil {
		return err
	}

	err = i.validateData()
	if err != nil {
		i.data = invoiceRequestData{}
		return err
	}

	return nil
}

func (i *Invoice) validateData() (err error) {
	//TODO implement
	return err
}

func (i *Invoice) LogError(err error) {
	var errStr string

	if _, ok := err.(*errorsWithStack.Error); ok && i.printErrStack {
		errStr = err.(*errorsWithStack.Error).ErrorStack()
	} else {
		errStr = err.Error()
	}

	i.logger.Error().Msgf(errStr)
}

func (i *Invoice) GeneratePDF() (*gofpdf.Fpdf, error) {
	i.logger.Debug().Msg("generate invoice")

	pdfGen, err := generator.NewPDFGenerator(generator.MetaData{
		FontName:     "OpenSans",
		FontGapY:     1.3,
		FontSize:     i.meta.Font.SizeDefault,
		MarginLeft:   i.meta.Margin.Left,
		MarginTop:    i.meta.Margin.Top,
		MarginRight:  i.meta.Margin.Right,
		MarginBottom: i.meta.Margin.Bottom,
		Unit:         "mm",
	}, false, i.logger)

	if err != nil {
		return nil, err
	}

	i.pdfGen = pdfGen

	if i.data.SenderInfo.MimeLogoUrl != "" {
		i.printMimeImg()
	}

	i.printAddressee()
	i.printMetaData(pdfGen)
	i.printHeadlineAndOpeningText(pdfGen)
	i.printInvoiceTable(pdfGen)
	i.printClosingText(pdfGen)
	i.printFooter(pdfGen)
	i.printFooter(pdfGen)

	return pdfGen.GetPdf(), pdfGen.GetError()
}

func (i *Invoice) printMimeImg() {
	pageWidth, _ := i.pdfGen.GetPdf().GetPageSize()
	mimeImg(i.pdfGen, i.data.SenderInfo.MimeLogoUrl, pageWidth-i.meta.Margin.Right, 15, i.data.SenderInfo.MimeLogoScale)
}

func (i *Invoice) printAddressee() {
	//pageWidth, _ := i.pdfGen.GetPdf().GetPageSize()
	//i.pdfGen.DrawLine(i.meta.Margin.Left, i.meta.Margin.Top, pageWidth-i.meta.Margin.Right, i.meta.Margin.Top, lineColor, 0)

	letterAddressSenderSmall(i.pdfGen, getAddressLine(i.data.SenderAddress), i.meta.Margin.Left, 49, i.meta.Font.SizeSmall)
	i.pdfGen.SetFontSize(i.meta.Font.SizeDefault)

	letterReceiverAddress(i.pdfGen, i.data.ReceiverAddress, i.meta.Margin.Left, 56)
}

func (i *Invoice) printMetaData(pdfGen *generator.PDFGenerator) {
	pdfGen.SetFontSize(i.meta.Font.SizeDefault)
	pdfGen.DrawLine(i.meta.Margin.Left+98, 56, i.meta.Margin.Left+98, 80, i.defaultLineColor, 0)
	pdfGen.SetCursor(i.meta.Margin.Left+100, 56)
	pdfGen.PrintLnPdfText("Kundennummer:", "", "L")
	pdfGen.PrintLnPdfText("Rechnungsnummer:", "", "L")
	pdfGen.PrintLnPdfText("Datum:", "", "L")

	pdfGen.SetCursor(i.meta.Margin.Left+140, 56)
	pdfGen.PrintLnPdfText(i.data.InvoiceMeta.CustomerNumber, "", "L")
	pdfGen.PrintLnPdfText(i.data.InvoiceMeta.InvoiceNumber, "", "L")
	pdfGen.PrintLnPdfText(i.data.InvoiceMeta.InvoiceDate, "", "L")
}

func (i *Invoice) printHeadlineAndOpeningText(pdfGen *generator.PDFGenerator) {
	//Überschrift
	pdfGen.SetCursor(i.meta.Margin.Left, 100)
	pdfGen.SetFontSize(i.meta.Font.SizeLarge)
	pdfGen.PrintLnPdfText(i.data.InvoiceBody.HeadlineText+" "+i.data.InvoiceMeta.InvoiceNumber, "b", "L")

	//opening
	pdfGen.SetFontSize(i.meta.Font.SizeDefault)
	pdfGen.NewLine(pdfGen.GetMarginLeft())
	pdfGen.PrintLnPdfText(i.data.InvoiceBody.OpeningText, "", "L")
}

func (i *Invoice) printInvoiceTable(pdfGen *generator.PDFGenerator) {
	var invoicedItems = [][]string{{}}

	type taxSumType struct {
		taxName string
		taxSum  float64
	}

	var netSum float64

	var taxSums []taxSumType

	for _, product := range i.data.InvoiceBody.InvoicedItems {
		netSum += product.Quantity * (float64(product.SinglePrice) / float64(100))

		//check if taxRate already exists
		var taxSumExists = false
		for i, taxSum := range taxSums {
			if taxSum.taxName == strconv.Itoa(product.TaxRate)+"%" {
				taxSums[i].taxSum += product.Quantity * (float64(product.SinglePrice) / float64(100)) * (float64(product.TaxRate) / float64(100))
				taxSumExists = true
			}
		}
		if !taxSumExists {
			taxSums = append(taxSums, taxSumType{taxName: strconv.Itoa(product.TaxRate) + "%",
				taxSum: product.Quantity * (float64(product.SinglePrice) / float64(100)) * (float64(product.TaxRate) / float64(100))})
		}

		invoicedItems = append(invoicedItems,
			[]string{
				product.PositionNumber,
				germanNumber(product.Quantity) + " " + product.Unit,
				germanNumber(float64(product.SinglePrice)/float64(100)) + "€",
				product.Description,
				strconv.Itoa(product.TaxRate) + "%",
				germanNumber(product.Quantity*(float64(product.SinglePrice)/float64(100))) + "€",
			})
	}

	var headerCells = []string{"Pos", "Anzahl", "Preis", "Beschreibung", "USt", "Netto"}
	var columnPercent = []float64{6, 10, 10, 54, 8, 12}
	var columnWidth = getColumnWithFromPercentage(pdfGen, columnPercent)

	var headerCellAlign = []string{"LM", "LM", "LM", "LM", "RM", "RM"}
	var bodyCellAlign = []string{"LM", "LM", "LM", "LM", "RM", "RM"}
	var summaryCells = [][]string{
		{"", "Zwischensumme", germanNumber(netSum) + "€"},
	}
	//summaryCells append taxSums
	for _, taxSum := range taxSums {
		//append only if txSum is not 0
		if taxSum.taxSum != 0 {
			summaryCells = append(summaryCells, []string{"", taxSum.taxName, germanNumber(taxSum.taxSum) + "€"})
		}
	}

	var totalTax = 0.
	for _, taxSum := range taxSums {
		totalTax += taxSum.taxSum
	}

	//add last row with total sum, calculated from netSum plus each taxSum
	summaryCells = append(summaryCells, []string{"", "Gesamtbetrag", germanNumber(totalTax+netSum) + "€"})

	var summaryColumnPercent = []float64{60, 25, 15}
	var summaryColumnWidths = getColumnWithFromPercentage(pdfGen, summaryColumnPercent)
	var summaryCellAlign = []string{"LM", "LM", "RM"}

	pdfGen.NewLine(i.meta.Margin.Left)
	pdfGen.SetFontSize(i.meta.Font.SizeSmall)
	pdfGen.PrintLnPdfText(i.data.InvoiceBody.ServiceTimeText, "i", "L")
	pdfGen.SetFontSize(i.meta.Font.SizeDefault)

	pdfGen.PrintTableHeader(headerCells, columnWidth, headerCellAlign)
	pdfGen.PrintTableBody(invoicedItems, columnWidth, bodyCellAlign)
	pdfGen.PrintTableFooter(summaryCells, summaryColumnWidths, summaryCellAlign)
}

func (i *Invoice) printClosingText(pdfGen *generator.PDFGenerator) {
	pdfGen.SetFontSize(i.meta.Font.SizeDefault)
	pdfGen.NewLine(i.meta.Margin.Left)
	pdfGen.NewLine(i.meta.Margin.Left)
	pdfGen.NewLine(i.meta.Margin.Left)
	pdfGen.PrintLnPdfText(i.data.InvoiceBody.ClosingText, "", "L")
	pdfGen.NewLine(i.meta.Margin.Left)
	pdfGen.NewLine(i.meta.Margin.Left)
	pdfGen.PrintLnPdfText(i.data.InvoiceBody.UstNotice, "", "L")
}

func (i *Invoice) printFooter(pdfGen *generator.PDFGenerator) {
	const startAtY = 261
	const startPageNumberY = 282
	const gabY = 3

	pageWidth, _ := pdfGen.GetPdf().GetPageSize()

	pdfGen.SetFontSize(i.meta.Font.SizeSmall)
	pdfGen.DrawLine(i.meta.Margin.Left, startAtY, pageWidth-i.meta.Margin.Right, startAtY, i.defaultLineColor, 0)

	pdfGen.SetCursor(i.meta.Margin.Left, startAtY+gabY)
	pdfGen.PrintLnPdfText(i.data.SenderInfo.Web, "", "L")
	pdfGen.PrintLnPdfText(i.data.SenderInfo.Phone, "", "L")
	pdfGen.PrintLnPdfText(i.data.SenderInfo.Email, "", "L")

	pdfGen.SetCursor(pageWidth/2, startAtY+gabY)
	pdfGen.PrintLnPdfText(i.data.SenderAddress.CompanyName, "", "C")
	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", i.data.SenderAddress.Address.Road, i.data.SenderAddress.Address.HouseNumber), "", "C")
	pdfGen.PrintLnPdfText(i.data.SenderAddress.Address.ZipCode+" "+i.data.SenderAddress.Address.CityName, "", "C")
	pdfGen.PrintLnPdfText(i.data.SenderInfo.TaxNumber, "", "C")

	pdfGen.SetCursor(pageWidth-i.meta.Margin.Right, startAtY+gabY)
	pdfGen.PrintLnPdfText(i.data.SenderInfo.BankName, "", "R")
	pdfGen.PrintLnPdfText(i.data.SenderInfo.Iban, "", "R")
	pdfGen.PrintLnPdfText(i.data.SenderInfo.Bic, "", "R")

	pdfGen.DrawLine(i.meta.Margin.Left, startPageNumberY, pageWidth-i.meta.Margin.Right, startPageNumberY, i.defaultLineColor, 0)
	pdfGen.SetCursor(pageWidth/2, startPageNumberY+gabY)
	pdfGen.PrintLnPdfText("Seite 1 von 1", "", "C")
	pdfGen.SetFontSize(i.meta.Font.SizeDefault)
}
