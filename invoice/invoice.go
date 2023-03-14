package invoice

import (
	"SimpleInvoice/generator"
	"errors"
	"fmt"
	errorsWithStack "github.com/go-errors/errors"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"net/http"
	"net/url"
	"strconv"
)

type Invoice struct {
	pdfData         pdfInvoiceData
	logger          *zerolog.Logger
	textFont        string
	marginLeft      float64
	marginRight     float64
	marginTop       float64
	marginBottom    float64
	printErrStack   bool
	defaultFontSize float64
	smallFontSize   float64
	headerFontSize  float64
}
type pdfInvoiceData struct {
	SenderAddress struct {
		FullForename string `json:"fullForename"`
		FullSurname  string `json:"fullSurname"`
		CompanyName  string `json:"companyName"`
		Supplement   string `json:"supplement"`
		Address      struct {
			Road             string `json:"road"`
			HouseNumber      string `json:"houseNumber"`
			StreetSupplement string `json:"streetSupplement"`
			ZipCode          string `json:"zipCode"`
			CityName         string `json:"cityName"`
			Country          string `json:"country"`
			CountryCode      string `json:"countryCode"`
		} `json:"address"`
	} `json:"senderAddress"`
	ReceiverAddress struct {
		FullForename string `json:"fullForename"`
		FullSurname  string `json:"fullSurname"`
		CompanyName  string `json:"companyName"`
		Supplement   string `json:"supplement"`
		Address      struct {
			Road             string `json:"road"`
			HouseNumber      string `json:"houseNumber"`
			StreetSupplement string `json:"streetSupplement"`
			ZipCode          string `json:"zipCode"`
			CityName         string `json:"cityName"`
			Country          string `json:"country"`
			CountryCode      string `json:"countryCode"`
		} `json:"address"`
	} `json:"receiverAddress"`
	SenderInfo struct {
		Phone         string  `json:"phone"`
		Web           string  `json:"web"`
		Email         string  `json:"email"`
		MimeLogoUrl   string  `json:"mimeLogoUrl"`
		MimeLogoScale float64 `json:"mimeLogoScale"`
		Iban          string  `json:"iban"`
		Bic           string  `json:"bic"`
		TaxNumber     string  `json:"taxNumber"`
		BankName      string  `json:"bankName"`
	} `json:"senderInfo"`
	InvoiceMeta struct {
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

func New(logger *zerolog.Logger) (iv *Invoice) {

	iv = &Invoice{
		logger:        logger,
		textFont:      "openSans",
		marginLeft:    25,
		marginRight:   20,
		marginTop:     45,
		marginBottom:  0,
		printErrStack: logger.GetLevel() <= zerolog.DebugLevel,

		pdfData:         pdfInvoiceData{},
		defaultFontSize: 10,
		smallFontSize:   8,
		headerFontSize:  15,
	}

	return iv
}

func (iv *Invoice) SetJsonInvoiceData(request *http.Request) (err error) {
	err = iv.parseJsonData(request)
	if err != nil {
		iv.LogError(err)
		return errors.New("data parsing Failed")
	}

	err = iv.validateJsonData()
	if err != nil {
		iv.pdfData = pdfInvoiceData{}
		iv.LogError(err)
		return errors.New("data parsing Failed")
	}

	return nil
}

func germanNumber(n float64) string {
	p := message.NewPrinter(language.German)
	return p.Sprintf("%.2f", n)
}

func (iv *Invoice) GeneratePDF() (*gofpdf.Fpdf, error) {

	iv.logger.Debug().Msg("Endpoint Hit: pdfPage")

	lineColor := generator.Color{R: 200, G: 200, B: 200}

	pdfGen, err := generator.NewPDFGenerator(generator.MetaData{
		FontName:     "openSans",
		FontGapY:     1.3,
		FontSize:     iv.defaultFontSize,
		MarginLeft:   iv.marginLeft,
		MarginTop:    iv.marginTop,
		MarginRight:  iv.marginRight,
		MarginBottom: iv.marginBottom,
		Unit:         "mm",
	}, false)

	if err != nil {
		iv.LogError(err)
		return nil, err
	}

	iv.printMimeImg(pdfGen)
	iv.printAddressee(pdfGen, lineColor)
	iv.printMetaData(pdfGen, lineColor)
	iv.printHeadlineAndOpeningText(pdfGen)
	iv.printInvoiceTable(pdfGen)
	iv.printClosingText(pdfGen)
	iv.printFooter(pdfGen, lineColor)

	if pdfGen.GetError() != nil {
		iv.LogError(pdfGen.GetError())
	}

	return pdfGen.GetPdf(), pdfGen.GetError()
}

func (iv *Invoice) printMimeImg(pdfGen *generator.PDFGenerator) {
	urlStruct, err := url.Parse(iv.pdfData.SenderInfo.MimeLogoUrl)
	if err != nil {
		iv.logger.Error().Msg(err.Error())
		pdfGen.SetError(errorsWithStack.New(err.Error()))
		return
	}

	pageWidth, _ := pdfGen.GetPdf().GetPageSize()
	pdfGen.SetUnsafeCursor(pageWidth-iv.marginRight, 15)
	pdfGen.PlaceMimeImageFromUrl(urlStruct, iv.pdfData.SenderInfo.MimeLogoScale, "R")
	if pdfGen.GetError() != nil {
		iv.logger.Error().Msg(err.Error())
		return
	}
}

func (iv *Invoice) printAddressee(pdfGen *generator.PDFGenerator, lineColor generator.Color) {
	pageWidth, _ := pdfGen.GetPdf().GetPageSize()
	pdfGen.DrawLine(iv.marginLeft, iv.marginTop, pageWidth-iv.marginRight, iv.marginTop, lineColor, 0)

	//Anschrift Sender small
	pdfGen.SetCursor(iv.marginLeft, 49)
	pdfGen.SetFontSize(iv.smallFontSize)
	addressSenderSmallText := fmt.Sprintf("%s,%s %s, %s %s %s",
		iv.pdfData.SenderAddress.CompanyName,
		iv.pdfData.SenderAddress.Address.Road,
		iv.pdfData.SenderAddress.Address.HouseNumber,
		iv.pdfData.SenderAddress.Address.CountryCode,
		iv.pdfData.SenderAddress.Address.ZipCode,
		iv.pdfData.SenderAddress.Address.CityName)
	pdfGen.PrintPdfText(addressSenderSmallText, "", "L")
	pdfGen.SetFontSize(iv.defaultFontSize)

	//Anschrift Empfänger
	pdfGen.SetCursor(iv.marginLeft, 56)
	pdfGen.PrintLnPdfText(iv.pdfData.ReceiverAddress.CompanyName, "", "L")
	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", iv.pdfData.ReceiverAddress.FullForename, iv.pdfData.ReceiverAddress.FullSurname),
		"", "L")
	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", iv.pdfData.ReceiverAddress.Address.Road, iv.pdfData.ReceiverAddress.Address.HouseNumber),
		"", "L")
	if iv.pdfData.ReceiverAddress.Address.StreetSupplement != "" {
		pdfGen.PrintLnPdfText(iv.pdfData.ReceiverAddress.Address.StreetSupplement, "", "L")
	}
	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", iv.pdfData.ReceiverAddress.Address.ZipCode, iv.pdfData.ReceiverAddress.Address.CityName),
		"", "L")
}

func (iv *Invoice) printMetaData(pdfGen *generator.PDFGenerator, lineColor generator.Color) {
	pdfGen.SetFontSize(iv.defaultFontSize)
	pdfGen.DrawLine(iv.marginLeft+98, 56, iv.marginLeft+98, 80, lineColor, 0)
	pdfGen.SetCursor(iv.marginLeft+100, 56)
	pdfGen.PrintLnPdfText("Kundennummer:", "", "L")
	pdfGen.PrintLnPdfText("Rechnungsnummer:", "", "L")
	pdfGen.PrintLnPdfText("Datum:", "", "L")

	pdfGen.SetCursor(iv.marginLeft+140, 56)
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceMeta.CustomerNumber, "", "L")
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceMeta.InvoiceNumber, "", "L")
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceMeta.InvoiceDate, "", "L")
}

func (iv *Invoice) printHeadlineAndOpeningText(pdfGen *generator.PDFGenerator) {
	//Überschrift
	pdfGen.SetCursor(iv.marginLeft, 100)
	pdfGen.SetFontSize(iv.headerFontSize)
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceBody.HeadlineText+" "+iv.pdfData.InvoiceMeta.InvoiceNumber, "b", "L")

	//opening
	pdfGen.SetFontSize(iv.defaultFontSize)
	pdfGen.NewLine(pdfGen.GetMarginLeft())
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceBody.OpeningText, "", "L")
}

func (iv *Invoice) printInvoiceTable(pdfGen *generator.PDFGenerator) {
	getCellWith := func(percent float64) float64 {
		maxSavePrintingWidth, _ := pdfGen.GetPdf().GetPageSize()
		maxSavePrintingWidth = maxSavePrintingWidth - pdfGen.GetMarginLeft() - pdfGen.GetMarginRight()

		return (percent * maxSavePrintingWidth) / 100.0
	}

	var invoicedItems = [][]string{{}}

	type taxSumType struct {
		taxName string
		taxSum  float64
	}

	var netSum float64

	var taxSums []taxSumType

	for _, product := range iv.pdfData.InvoiceBody.InvoicedItems {
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
				germanNumber(product.Quantity * (float64(product.SinglePrice) / float64(100))),
			})
	}

	var headerCells = []string{"Pos", "Anzahl", "Preis", "Beschreibung", "USt", "Netto"}
	var columnWidth = []float64{getCellWith(6), getCellWith(10), getCellWith(10), getCellWith(54), getCellWith(8), getCellWith(12)}

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

	var summaryColumnWidths = []float64{getCellWith(60), getCellWith(25), getCellWith(15)}
	var summaryCellAlign = []string{"LM", "LM", "RM"}

	pdfGen.NewLine(iv.marginLeft)
	pdfGen.SetFontSize(iv.smallFontSize)
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceBody.ServiceTimeText, "i", "L")
	pdfGen.SetFontSize(iv.defaultFontSize)

	pdfGen.PrintTableHeader(headerCells, columnWidth, headerCellAlign)
	pdfGen.PrintTableBody(invoicedItems, columnWidth, bodyCellAlign)
	pdfGen.PrintTableFooter(summaryCells, summaryColumnWidths, summaryCellAlign)
}

func (iv *Invoice) printClosingText(pdfGen *generator.PDFGenerator) {
	pdfGen.SetFontSize(iv.defaultFontSize)
	pdfGen.NewLine(iv.marginLeft)
	pdfGen.NewLine(iv.marginLeft)
	pdfGen.NewLine(iv.marginLeft)
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceBody.ClosingText, "", "L")
	pdfGen.NewLine(iv.marginLeft)
	pdfGen.NewLine(iv.marginLeft)
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceBody.UstNotice, "", "L")
}

func (iv *Invoice) printFooter(pdfGen *generator.PDFGenerator, lineColor generator.Color) {
	pageWidth, _ := pdfGen.GetPdf().GetPageSize()

	pdfGen.SetFontSize(iv.smallFontSize)
	pdfGen.DrawLine(iv.marginLeft, 261, pageWidth-iv.marginRight, 261, lineColor, 0)
	pdfGen.SetCursor(iv.marginLeft, 264)
	pdfGen.PrintLnPdfText(iv.pdfData.SenderInfo.Web, "", "L")
	pdfGen.PrintLnPdfText(iv.pdfData.SenderInfo.Phone, "", "L")
	pdfGen.PrintLnPdfText(iv.pdfData.SenderInfo.Email, "", "L")
	pdfGen.SetCursor(105, 264)
	pdfGen.PrintLnPdfText(iv.pdfData.SenderAddress.CompanyName, "", "C")
	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", iv.pdfData.SenderAddress.Address.Road, iv.pdfData.SenderAddress.Address.HouseNumber), "", "C")
	pdfGen.PrintLnPdfText(iv.pdfData.SenderAddress.Address.ZipCode+" "+iv.pdfData.SenderAddress.Address.CityName, "", "C")
	pdfGen.PrintLnPdfText(iv.pdfData.SenderInfo.TaxNumber, "", "C")
	pdfGen.SetCursor(190, 264)
	pdfGen.PrintLnPdfText(iv.pdfData.SenderInfo.BankName, "", "R")
	pdfGen.PrintLnPdfText(iv.pdfData.SenderInfo.Iban, "", "R")
	pdfGen.PrintLnPdfText(iv.pdfData.SenderInfo.Bic, "", "R")
	pdfGen.DrawLine(iv.marginLeft, 282, pageWidth-iv.marginRight, 282, lineColor, 0)
	pdfGen.SetFontSize(iv.defaultFontSize)

	pdfGen.SetCursor(pageWidth/2, 285)
	pdfGen.SetFontSize(iv.smallFontSize)
	pdfGen.PrintLnPdfText("Seite 1 von 1", "", "C")
	pdfGen.SetFontSize(iv.defaultFontSize)
}

func (iv *Invoice) LogError(err error) {
	var errStr string

	if _, ok := err.(*errorsWithStack.Error); ok && iv.printErrStack {
		errStr = err.(*errorsWithStack.Error).ErrorStack()
	} else {
		errStr = err.Error()
	}

	iv.logger.Error().Msgf(errStr)
}
