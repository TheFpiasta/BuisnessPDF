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
	"net/url"
	"strconv"
)

type Invoice struct {
	data          invoiceRequestData
	meta          pdfMeta
	logger        *zerolog.Logger
	printErrStack bool
}

type invoiceRequestData struct {
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
	//todo rename
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

func NewInvoice(logger *zerolog.Logger) *Invoice {
	return &Invoice{
		data: invoiceRequestData{},
		meta: pdfMeta{
			margin: pdfMargin{
				left:   25,
				right:  20,
				top:    45,
				bottom: 0,
			},
			font: pdfFont{
				fontName:    "openSans",
				sizeDefault: 10,
				sizeSmall:   8,
				SizeLarge:   15,
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
	//todo implement
	return err
}

func (i *Invoice) GeneratePDF() (*gofpdf.Fpdf, error) {
	i.logger.Debug().Msg("Endpoint Hit: pdfPage")

	lineColor := generator.Color{R: 200, G: 200, B: 200}

	pdfGen, err := generator.NewPDFGenerator(generator.MetaData{
		FontName:     "OpenSans",
		FontGapY:     1.3,
		FontSize:     i.meta.font.sizeDefault,
		MarginLeft:   i.meta.margin.left,
		MarginTop:    i.meta.margin.top,
		MarginRight:  i.meta.margin.right,
		MarginBottom: i.meta.margin.bottom,
		Unit:         "mm",
	}, false, i.logger)

	if err != nil {
		return nil, err
	}

	if i.data.SenderInfo.MimeLogoUrl != "" {
		i.printMimeImg(pdfGen)
	}

	i.printAddressee(pdfGen, lineColor)
	i.printMetaData(pdfGen, lineColor)
	i.printHeadlineAndOpeningText(pdfGen)
	i.printInvoiceTable(pdfGen)
	i.printClosingText(pdfGen)
	i.printFooter(pdfGen, lineColor)

	return pdfGen.GetPdf(), pdfGen.GetError()
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

func (i *Invoice) printMimeImg(pdfGen *generator.PDFGenerator) {
	urlStruct, err := url.Parse(i.data.SenderInfo.MimeLogoUrl)
	if err != nil {
		pdfGen.SetError(errorsWithStack.New(err.Error()))
		return
	}

	pageWidth, _ := pdfGen.GetPdf().GetPageSize()
	pdfGen.SetUnsafeCursor(pageWidth-i.meta.margin.right, 15)
	pdfGen.PlaceMimeImageFromUrl(urlStruct, i.data.SenderInfo.MimeLogoScale, "R")
}

func (i *Invoice) printAddressee(pdfGen *generator.PDFGenerator, lineColor generator.Color) {
	pageWidth, _ := pdfGen.GetPdf().GetPageSize()
	pdfGen.DrawLine(i.meta.margin.left, i.meta.margin.top, pageWidth-i.meta.margin.right, i.meta.margin.top, lineColor, 0)

	//Anschrift Sender small
	pdfGen.SetCursor(i.meta.margin.left, 49)
	pdfGen.SetFontSize(i.meta.font.sizeSmall)

	var addressSenderSmallText = ""

	addressSenderSmallText += i.data.SenderAddress.CompanyName
	if i.data.SenderAddress.CompanyName != "" && (i.data.SenderAddress.FullForename != "" || i.data.SenderAddress.FullSurname != "") {
		addressSenderSmallText += ", "
	}

	addressSenderSmallText += i.data.SenderAddress.FullForename
	if i.data.SenderAddress.FullSurname != "" {
		addressSenderSmallText += " "
	}
	addressSenderSmallText += i.data.SenderAddress.FullSurname

	addressSenderSmallText += fmt.Sprintf(" - %s %s",
		i.data.SenderAddress.Address.Road,
		i.data.SenderAddress.Address.HouseNumber,
	)

	if i.data.SenderAddress.Address.StreetSupplement != "" {
		addressSenderSmallText += ", "
		addressSenderSmallText += i.data.SenderAddress.Address.StreetSupplement
	}

	addressSenderSmallText += fmt.Sprintf(", %s %s %s",
		i.data.SenderAddress.Address.CountryCode,
		i.data.SenderAddress.Address.ZipCode,
		i.data.SenderAddress.Address.CityName,
	)

	pdfGen.PrintPdfText(addressSenderSmallText, "", "L")
	pdfGen.SetFontSize(i.meta.font.sizeDefault)

	//Anschrift Empfänger
	pdfGen.SetCursor(i.meta.margin.left, 56)
	if i.data.ReceiverAddress.CompanyName != "" {
		pdfGen.PrintLnPdfText(i.data.ReceiverAddress.CompanyName, "", "L")

	}
	if i.data.ReceiverAddress.FullForename != "" || i.data.ReceiverAddress.FullSurname != "" {
		pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", i.data.ReceiverAddress.FullForename, i.data.ReceiverAddress.FullSurname),
			"", "L")
	}
	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", i.data.ReceiverAddress.Address.Road, i.data.ReceiverAddress.Address.HouseNumber),
		"", "L")
	if i.data.ReceiverAddress.Address.StreetSupplement != "" {
		pdfGen.PrintLnPdfText(i.data.ReceiverAddress.Address.StreetSupplement, "", "L")
	}
	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", i.data.ReceiverAddress.Address.ZipCode, i.data.ReceiverAddress.Address.CityName),
		"", "L")
}

func (i *Invoice) printMetaData(pdfGen *generator.PDFGenerator, lineColor generator.Color) {
	pdfGen.SetFontSize(i.meta.font.sizeDefault)
	pdfGen.DrawLine(i.meta.margin.left+98, 56, i.meta.margin.left+98, 80, lineColor, 0)
	pdfGen.SetCursor(i.meta.margin.left+100, 56)
	pdfGen.PrintLnPdfText("Kundennummer:", "", "L")
	pdfGen.PrintLnPdfText("Rechnungsnummer:", "", "L")
	pdfGen.PrintLnPdfText("Datum:", "", "L")

	pdfGen.SetCursor(i.meta.margin.left+140, 56)
	pdfGen.PrintLnPdfText(i.data.InvoiceMeta.CustomerNumber, "", "L")
	pdfGen.PrintLnPdfText(i.data.InvoiceMeta.InvoiceNumber, "", "L")
	pdfGen.PrintLnPdfText(i.data.InvoiceMeta.InvoiceDate, "", "L")
}

func (i *Invoice) printHeadlineAndOpeningText(pdfGen *generator.PDFGenerator) {
	//Überschrift
	pdfGen.SetCursor(i.meta.margin.left, 100)
	pdfGen.SetFontSize(i.meta.font.SizeLarge)
	pdfGen.PrintLnPdfText(i.data.InvoiceBody.HeadlineText+" "+i.data.InvoiceMeta.InvoiceNumber, "b", "L")

	//opening
	pdfGen.SetFontSize(i.meta.font.sizeDefault)
	pdfGen.NewLine(pdfGen.GetMarginLeft())
	pdfGen.PrintLnPdfText(i.data.InvoiceBody.OpeningText, "", "L")
}

func (i *Invoice) printInvoiceTable(pdfGen *generator.PDFGenerator) {
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

	pdfGen.NewLine(i.meta.margin.left)
	pdfGen.SetFontSize(i.meta.font.sizeSmall)
	pdfGen.PrintLnPdfText(i.data.InvoiceBody.ServiceTimeText, "i", "L")
	pdfGen.SetFontSize(i.meta.font.sizeDefault)

	pdfGen.PrintTableHeader(headerCells, columnWidth, headerCellAlign)
	pdfGen.PrintTableBody(invoicedItems, columnWidth, bodyCellAlign)
	pdfGen.PrintTableFooter(summaryCells, summaryColumnWidths, summaryCellAlign)
}

func (i *Invoice) printClosingText(pdfGen *generator.PDFGenerator) {
	pdfGen.SetFontSize(i.meta.font.sizeDefault)
	pdfGen.NewLine(i.meta.margin.left)
	pdfGen.NewLine(i.meta.margin.left)
	pdfGen.NewLine(i.meta.margin.left)
	pdfGen.PrintLnPdfText(i.data.InvoiceBody.ClosingText, "", "L")
	pdfGen.NewLine(i.meta.margin.left)
	pdfGen.NewLine(i.meta.margin.left)
	pdfGen.PrintLnPdfText(i.data.InvoiceBody.UstNotice, "", "L")
}

func (i *Invoice) printFooter(pdfGen *generator.PDFGenerator, lineColor generator.Color) {
	pageWidth, _ := pdfGen.GetPdf().GetPageSize()

	pdfGen.SetFontSize(i.meta.font.sizeSmall)
	pdfGen.DrawLine(i.meta.margin.left, 261, pageWidth-i.meta.margin.right, 261, lineColor, 0)
	pdfGen.SetCursor(i.meta.margin.left, 264)
	pdfGen.PrintLnPdfText(i.data.SenderInfo.Web, "", "L")
	pdfGen.PrintLnPdfText(i.data.SenderInfo.Phone, "", "L")
	pdfGen.PrintLnPdfText(i.data.SenderInfo.Email, "", "L")
	pdfGen.SetCursor(105, 264)
	pdfGen.PrintLnPdfText(i.data.SenderAddress.CompanyName, "", "C")
	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", i.data.SenderAddress.Address.Road, i.data.SenderAddress.Address.HouseNumber), "", "C")
	pdfGen.PrintLnPdfText(i.data.SenderAddress.Address.ZipCode+" "+i.data.SenderAddress.Address.CityName, "", "C")
	pdfGen.PrintLnPdfText(i.data.SenderInfo.TaxNumber, "", "C")
	pdfGen.SetCursor(190, 264)
	pdfGen.PrintLnPdfText(i.data.SenderInfo.BankName, "", "R")
	pdfGen.PrintLnPdfText(i.data.SenderInfo.Iban, "", "R")
	pdfGen.PrintLnPdfText(i.data.SenderInfo.Bic, "", "R")
	pdfGen.DrawLine(i.meta.margin.left, 282, pageWidth-i.meta.margin.right, 282, lineColor, 0)
	pdfGen.SetFontSize(i.meta.font.sizeDefault)

	pdfGen.SetCursor(pageWidth/2, 285)
	pdfGen.SetFontSize(i.meta.font.sizeSmall)
	pdfGen.PrintLnPdfText("Seite 1 von 1", "", "C")
	pdfGen.SetFontSize(i.meta.font.sizeDefault)
}
