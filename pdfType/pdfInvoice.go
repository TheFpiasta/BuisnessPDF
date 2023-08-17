package pdfType

import (
	"SimpleInvoice/generator"
	din5008a "SimpleInvoice/norms/letter/din-5008-a"
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
	data          invoiceRequestData
	meta          PdfMeta
	logger        *zerolog.Logger
	printErrStack bool
	pdfGen        *generator.PDFGenerator
	footerStartY  float64
}

type invoiceRequestData struct {
	SenderAddress   din5008a.FullAdresse `json:"senderAddress"`
	ReceiverAddress din5008a.FullAdresse `json:"receiverAddress"`
	SenderInfo      SenderInfo           `json:"senderInfo"`
	InvoiceMeta     struct {
		InvoiceNumber  string            `json:"invoiceNumber"`
		InvoiceDate    string            `json:"invoiceDate"`
		CustomerNumber string            `json:"customerNumber"`
		CustomMetaData []CustomMetaDatum `json:"customMetaData"`
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

type CustomMetaDatum struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewInvoice(logger *zerolog.Logger) *Invoice {
	return &Invoice{
		data: invoiceRequestData{},
		meta: PdfMeta{
			Font: pdfFont{
				FontName:    "openSans",
				SizeDefault: din5008a.FontSize10,
				SizeSmall:   din5008a.FontSizeSender8,
				SizeLarge:   din5008a.FontSize10 + 5,
			},
		},
		logger:        logger,
		printErrStack: logger.GetLevel() <= zerolog.DebugLevel,
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

	pdfGen, err := generator.NewPDFGenerator(
		generator.MetaData{
			FontName:         "OpenSans",
			FontGapY:         1.3,
			FontSize:         i.meta.Font.SizeDefault,
			MarginLeft:       din5008a.BodyStartX,
			MarginTop:        din5008a.AddressSenderTextStartY,
			MarginRight:      din5008a.Width - din5008a.BodyStopX,
			MarginBottom:     0,
			Unit:             "mm",
			DefaultLineWidth: 0.4,
			DefaultLineColor: generator.Color{R: 162, G: 162, B: 162},
		},
		false,
		i.logger,
		func() {
			i.printHeader()
		},
		func(isLastPage bool) {
			i.printFooter()
		},
	)

	if err != nil {
		return nil, err
	}

	i.pdfGen = pdfGen
	i.pdfGen.NewPage()

	i.doGeneratePdf()

	return pdfGen.GetPdf(), pdfGen.GetError()
}

func (i *Invoice) doGeneratePdf() {
	var infoData []din5008a.InfoData
	infoData = append(infoData, din5008a.InfoData{Name: "Kundennummer:", Value: i.data.InvoiceMeta.CustomerNumber})
	infoData = append(infoData, din5008a.InfoData{Name: "Rechnungsnummer:", Value: i.data.InvoiceMeta.InvoiceNumber})
	infoData = append(infoData, din5008a.InfoData{Name: "Datum:", Value: i.data.InvoiceMeta.InvoiceDate})
	//TODO check length and throw error, if over din norm
	for _, datum := range i.data.InvoiceMeta.CustomMetaData {
		infoData = append(infoData, din5008a.InfoData{Name: datum.Name, Value: datum.Value})
	}

	din5008a.FullAddressesAndInfoPart(i.pdfGen, i.data.SenderAddress, i.data.ReceiverAddress, infoData)

	din5008a.Body(i.pdfGen, func() {
		i.printHeadlineAndOpeningText()
		i.printInvoiceTable()
		i.printClosingText()
	})

	din5008a.PageNumbering(i.pdfGen, i.footerStartY)
}

func (i *Invoice) printHeadlineAndOpeningText() {
	//Überschrift
	i.pdfGen.SetFontSize(i.meta.Font.SizeLarge)
	i.pdfGen.PrintLnPdfText(i.data.InvoiceBody.HeadlineText+" "+i.data.InvoiceMeta.InvoiceNumber, "b", "L")

	//opening
	i.pdfGen.SetFontSize(din5008a.FontSize10)
	i.pdfGen.SetFontGapY(din5008a.FontGab10)
	i.pdfGen.NewLine(din5008a.BodyStartX)
	i.pdfGen.PrintLnPdfText(i.data.InvoiceBody.OpeningText, "", "L")
}

func (i *Invoice) printInvoiceTable() {
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
	var columnWidth = getColumnWithFromPercentage(i.pdfGen, columnPercent)

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
	var summaryColumnWidths = getColumnWithFromPercentage(i.pdfGen, summaryColumnPercent)
	var summaryCellAlign = []string{"LM", "LM", "RM"}

	i.pdfGen.NewLine(din5008a.BodyStartX)
	i.pdfGen.SetFontSize(i.meta.Font.SizeSmall)
	i.pdfGen.PrintLnPdfText(i.data.InvoiceBody.ServiceTimeText, "i", "L")
	i.pdfGen.SetFontSize(i.meta.Font.SizeDefault)

	i.pdfGen.PrintTableHeader(headerCells, columnWidth, headerCellAlign)
	i.pdfGen.PrintTableBody(invoicedItems, columnWidth, bodyCellAlign)
	i.pdfGen.PrintTableFooter(summaryCells, summaryColumnWidths, summaryCellAlign)
}

func (i *Invoice) printClosingText() {
	i.pdfGen.SetFontSize(i.meta.Font.SizeDefault)
	i.pdfGen.NewLine(din5008a.BodyStartX)
	i.pdfGen.NewLine(din5008a.BodyStartX)
	i.pdfGen.NewLine(din5008a.BodyStartX)
	i.pdfGen.PrintLnPdfText(i.data.InvoiceBody.ClosingText, "", "L")
	i.pdfGen.NewLine(din5008a.BodyStartX)
	i.pdfGen.NewLine(din5008a.BodyStartX)
	i.pdfGen.PrintLnPdfText(i.data.InvoiceBody.UstNotice, "", "L")
}

func (i *Invoice) printFooter() {
	footerStartY, err := din5008a.Footer(i.printFooterContent, i.pdfGen)

	if err != nil {
		i.pdfGen.SetError(err)
	}

	if i.footerStartY == 0 {
		i.footerStartY = footerStartY
	}
}

func (i *Invoice) printFooterContent(maxFooterHeight float64) (footerStartY float64) {
	// calculate height
	var currentStartX float64
	var currentY float64
	i.pdfGen.SetUnsafeCursor(din5008a.BodyStartX, maxFooterHeight)
	i.pdfGen.PreviousLine(0)
	i.pdfGen.PreviousLine(0)
	i.pdfGen.PreviousLine(0)
	i.pdfGen.PreviousLine(0)
	_, currentY = i.pdfGen.GetCursor()
	footerStartY = currentY

	currentStartX = din5008a.BodyStartX
	i.pdfGen.SetCursor(currentStartX, footerStartY)
	i.pdfGen.PrintLnPdfText(i.data.SenderInfo.Web, "", "L")
	i.pdfGen.PrintLnPdfText(i.data.SenderInfo.Phone, "", "L")
	i.pdfGen.PrintLnPdfText(i.data.SenderInfo.Email, "", "L")

	currentStartX = ((din5008a.BodyStopX - din5008a.BodyStartX) / 2) + din5008a.BodyStartX
	i.pdfGen.SetCursor(currentStartX, footerStartY)
	i.pdfGen.PrintLnPdfText(i.data.SenderAddress.CompanyName, "", "C")
	i.pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", i.data.SenderAddress.Address.Road, i.data.SenderAddress.Address.HouseNumber), "", "C")
	i.pdfGen.PrintLnPdfText(i.data.SenderAddress.Address.ZipCode+" "+i.data.SenderAddress.Address.CityName, "", "C")
	i.pdfGen.PrintLnPdfText(i.data.SenderInfo.TaxNumber, "", "C")

	currentStartX = din5008a.BodyStopX
	i.pdfGen.SetCursor(currentStartX, footerStartY)
	i.pdfGen.PrintLnPdfText(i.data.SenderInfo.BankName, "", "R")
	i.pdfGen.PrintLnPdfText(i.data.SenderInfo.Iban, "", "R")
	i.pdfGen.PrintLnPdfText(i.data.SenderInfo.Bic, "", "R")

	return footerStartY
}

func (i *Invoice) printHeader() {
	if i.data.SenderInfo.MimeLogoUrl != "" {
		din5008a.MimeImageHeader(i.pdfGen, i.data.SenderInfo.MimeLogoUrl)
	}
}
