package invoice

import (
	"SimpleInvoice/generator"
	"errors"
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
	pdfData       pdfInvoiceData
	pdf           *gofpdf.Fpdf
	logger        *zerolog.Logger
	textFont      string
	marginLeft    float64
	marginRight   float64
	marginTop     float64
	marginBottom  float64
	printErrStack bool
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
		Phone     string `json:"phone"`
		Email     string `json:"email"`
		LogoSvg   string `json:"logoSvg"`
		Iban      string `json:"iban"`
		Bic       string `json:"bic"`
		TaxNumber string `json:"taxNumber"`
		BankName  string `json:"bankName"`
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

		pdfData: pdfInvoiceData{},
		pdf:     nil,
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
	const defaultFontSize = 10
	const smallFontSize = 8
	const headerFontSize = 15

	pdfGen, err := generator.NewPDFGenerator(generator.MetaData{
		FontName:     "openSans",
		FontGapY:     1.3,
		FontSize:     defaultFontSize,
		MarginLeft:   iv.marginLeft,
		MarginTop:    iv.marginTop,
		MarginRight:  iv.marginRight,
		MarginBottom: iv.marginBottom,
		Unit:         "mm",
	}, false)

	pageWidth, _ := pdfGen.GetPdf().GetPageSize()

	urlStruct, err := url.Parse("https://cdn.pictro.de/logosIcons/stack-one_logo_vector_white_small.png")
	if err != nil {
		return iv.pdf, err
	}

	err = pdfGen.PlaceMimeImageFromUrl(urlStruct, 153, 20, 0.5)
	if err != nil {
		return iv.pdf, err
	}

	//Anschrift Sender small
	pdfGen.SetCursor(iv.marginLeft, 49)
	pdfGen.SetFontSize(smallFontSize)
	pdfGen.PrintPdfText(
		iv.pdfData.SenderAddress.CompanyName+","+
			iv.pdfData.SenderAddress.Address.Road+" "+
			iv.pdfData.SenderAddress.Address.HouseNumber+", "+
			iv.pdfData.SenderAddress.Address.CountryCode+" "+
			iv.pdfData.SenderAddress.Address.ZipCode+" "+
			iv.pdfData.SenderAddress.Address.CityName, "", "L")
	pdfGen.SetFontSize(defaultFontSize)

	//Anschrift Empfänger
	pdfGen.SetCursor(iv.marginLeft, 56)
	pdfGen.PrintLnPdfText(iv.pdfData.SenderAddress.CompanyName, "", "L")
	pdfGen.PrintLnPdfText(iv.pdfData.SenderAddress.FullForename+" "+iv.pdfData.SenderAddress.FullSurname,
		"", "L")
	pdfGen.PrintLnPdfText(iv.pdfData.SenderAddress.Address.Road+" "+iv.pdfData.SenderAddress.Address.HouseNumber,
		"", "L")
	pdfGen.PrintLnPdfText(iv.pdfData.SenderAddress.Address.StreetSupplement, "", "L")
	pdfGen.PrintLnPdfText(iv.pdfData.SenderAddress.Address.ZipCode+" "+iv.pdfData.SenderAddress.Address.CityName,
		"", "L")

	//MetaData Infos Rechnung in 2 Spalten
	pdfGen.DrawLine(iv.marginLeft+98, 56, iv.marginLeft+98, 80, lineColor, 0)
	pdfGen.SetCursor(iv.marginLeft+100, 56)
	pdfGen.PrintLnPdfText("Kundennummer:", "", "L")
	pdfGen.PrintLnPdfText("Rechnungsnummer:", "", "L")
	pdfGen.PrintLnPdfText("Datum:", "", "L")

	pdfGen.SetCursor(iv.marginLeft+140, 56)
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceMeta.CustomerNumber, "", "L")
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceMeta.InvoiceNumber, "", "L")
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceMeta.InvoiceDate, "", "L")

	//Überschrift
	pdfGen.SetCursor(iv.marginLeft, 100)
	pdfGen.SetFontSize(headerFontSize)
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceBody.HeadlineText+" "+iv.pdfData.InvoiceMeta.InvoiceNumber, "b", "L")
	pdfGen.SetFontSize(defaultFontSize)
	pdfGen.NewLine(pdfGen.GetMarginLeft())

	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceBody.OpeningText, "", "L")

	pdfGen.SetFontSize(smallFontSize)
	pdfGen.PrintLnPdfText(iv.pdfData.InvoiceBody.ServiceTimeText, "i", "L")
	pdfGen.SetFontSize(defaultFontSize)

	getCellWith := func(percent float64) float64 {
		maxSavePrintingWidth, _ := pdfGen.GetPdf().GetPageSize()
		maxSavePrintingWidth = maxSavePrintingWidth - pdfGen.GetMarginLeft() - pdfGen.GetMarginRight()

		return (percent * maxSavePrintingWidth) / 100.0
	}

	type productStruct struct {
		iterator    string
		count       float64
		unit        string
		price       int
		description string
		taxRate     int
	}
	var bodyText = [][]string{{}}

	var productList = []productStruct{
		{"1", 40.5, "h", 4500, "agiles Software-Testing, System-Monitoring, \n Programmierung", 19},
		{"2", 19, "h", 6500, "agiles Software-Testing", 7},
		{"2", 19, "h", 6500, "agiles Software-Testing", 7},
		{"2", 19, "h", 6500, "agiles Software-Testing", 7},
		{"1", 40.5, "h", 4500, "agiles Software-Testing, System-Monitoring, \n Programmierung", 19},
		{"2", 19, "h", 6500, "agiles Software-Testing", 7},
		{"2", 19, "h", 6500, "agiles Software-Testing", 0},
	}

	type taxSumType struct {
		taxName string
		taxSum  float64
	}

	var netSum = 0.

	var taxSums []taxSumType

	for _, product := range productList {
		netSum += product.count * (float64(product.price) / float64(100))

		//check if taxRate already exists
		var taxSumExists = false
		for i, taxSum := range taxSums {
			if taxSum.taxName == strconv.Itoa(product.taxRate)+"%" {
				taxSums[i].taxSum += product.count * (float64(product.price) / float64(100)) * (float64(product.taxRate) / float64(100))
				taxSumExists = true
			}
		}
		if !taxSumExists {
			taxSums = append(taxSums, taxSumType{taxName: strconv.Itoa(product.taxRate) + "%",
				taxSum: product.count * (float64(product.price) / float64(100)) * (float64(product.taxRate) / float64(100))})
		}

		bodyText = append(bodyText,
			[]string{
				product.iterator,
				germanNumber(product.count) + " " + product.unit,
				germanNumber(float64(product.price)/float64(100)) + "€",
				product.description,
				strconv.Itoa(product.taxRate) + "%",
				germanNumber(product.count * (float64(product.price) / float64(100))),
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
	summaryCells = append(summaryCells, []string{"", "Gesamtbetrag", germanNumber(float64(totalTax)+netSum) + "€"})

	var summaryColumnWidths = []float64{getCellWith(60), getCellWith(25), getCellWith(15)}
	var summaryCellAlign = []string{"LM", "LM", "RM"}

	pdfGen.PrintTableHeader(headerCells, columnWidth, headerCellAlign)
	pdfGen.PrintTableBody(bodyText, columnWidth, bodyCellAlign)
	pdfGen.PrintTableFooter(summaryCells, summaryColumnWidths, summaryCellAlign)

	pdfGen.SetFontSize(smallFontSize)
	pdfGen.DrawLine(iv.marginLeft, 261, pageWidth-iv.marginRight, 261, lineColor, 0)
	pdfGen.SetCursor(iv.marginLeft, 264)
	pdfGen.PrintLnPdfText("https://stack-one.tech", "", "L")
	pdfGen.PrintLnPdfText("015154897208", "", "L")
	pdfGen.PrintLnPdfText("carsten@stack-one.de", "", "L")
	pdfGen.SetCursor(105, 264)
	pdfGen.PrintLnPdfText("stack1 GmbH", "", "C")
	pdfGen.PrintLnPdfText("Floßplatz 24", "", "C")
	pdfGen.PrintLnPdfText("04107 Leipzig", "", "C")
	pdfGen.PrintLnPdfText("UstID: DE3498754987", "", "C")
	pdfGen.SetCursor(190, 264)
	pdfGen.PrintLnPdfText("Sparkasse Leipzig", "", "R")
	pdfGen.PrintLnPdfText("DE55 8605 5592 1090 3143 33", "", "R")
	pdfGen.PrintLnPdfText("BIC: WELADE8LXXX", "", "R")
	pdfGen.DrawLine(iv.marginLeft, 282, pageWidth-iv.marginRight, 282, lineColor, 0)
	pdfGen.SetFontSize(defaultFontSize)

	pdfGen.SetCursor(pageWidth/2, 285)
	pdfGen.SetFontSize(smallFontSize)
	pdfGen.PrintLnPdfText("Seite 1 von 1", "", "C")
	pdfGen.SetFontSize(defaultFontSize)

	if pdfGen.GetError() != nil {
		iv.LogError(pdfGen.GetError())
	}
	return pdfGen.GetPdf(), pdfGen.GetError()
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
