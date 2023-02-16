package invoice

import (
	"SimpleInvoice/generator"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"math"
	"net/url"
	"strconv"
)

type Invoice struct {
	pdfData       invoicePdfData
	pdf           *gofpdf.Fpdf
	logger        *zerolog.Logger
	textFont      string
	lineHeight    float64
	textSize      float64
	textSizeSmall float64
	marginLeft    float64
	marginRight   float64
	marginTop     float64
	marginBottom  float64
	fontGapY      float64
}

type invoicePdfData struct {
	SenderAddress   addressInfo `json:"senderAddress"`
	ReceiverAddress addressInfo `json:"receiverAddress"`
	SenderInfo      senderInfo  `json:"senderInfo"`
	InvoiceMeta     invoiceMeta `json:"invoiceMeta"`
	InvoiceBody     invoiceBody `json:"InvoiceBody"`
}

type addressInfo struct {
	FullForename string  `json:"fullForename"`
	FullSurname  string  `json:"fullSurname"`
	CompanyName  string  `json:"companyName"`
	Supplement   string  `json:"supplement"`
	Address      address `json:"address"`
}

type address struct {
	Road             string `json:"road"`
	HouseNumber      string `json:"houseNumber"`
	StreetSupplement string `json:"streetSupplement"`
	ZipCode          string `json:"zipCode"`
	CityName         string `json:"cityName"`
	Country          string `json:"country"`
	CountryCode      string `json:"countryCode"`
}

type senderInfo struct {
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	LogSvgo   string `json:"logoSvg"`
	Iban      string `json:"iban"`
	Bic       string `json:"bic"`
	TaxNumber string `json:"taxNumber"`
	BankName  string `json:"bankName"`
}

type invoiceMeta struct {
	InvoiceNumber  string `json:"invoiceNumber"`
	InvoiceDate    string `json:"invoiceDate"`
	CustomerNumber string `json:"customerNumber"`
}

type invoiceBody struct {
	OpeningText     string          `json:"openingText"`
	ServiceTimeText string          `json:"serviceTimeText"`
	HeadlineText    string          `json:"headlineText"`
	ClosingText     string          `json:"closingText"`
	UstNotice       string          `json:"ustNotice"`
	InvoicedItems   []invoicedItems `json:"invoicedItems"`
	InvoicedSum     invoicedSum     `json:"invoicedSum"`
}

type invoicedItems struct {
	PositionNumber int     `json:"positionNumber"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	Description    string  `json:"description"`
	SinglePrice    float64 `json:"singlePrice"`
	NetPrice       float64 `json:"netPrice"`
}

type invoicedSum struct {
	OverallPriceNet   float64 `json:"overallPriceNet"`
	OverallPriceGross float64 `json:"overallPriceGross"`
	OverallTaxes      float64 `json:"overallTaxes"`
	TaxesPercentage   float64 `json:"taxesPercentage"`
	Currency          string  `json:"currency"`
}

func New(logger *zerolog.Logger) (iv *Invoice) {

	iv = &Invoice{
		logger:       logger,
		textFont:     "openSans",
		marginLeft:   25,
		marginRight:  20,
		marginTop:    45,
		marginBottom: 0,
	}

	return iv
}

func (iv *Invoice) SetJsonInvoiceData(jsonData io.ReadCloser) (err error) {
	err = iv.parseJsonData(jsonData)
	if err != nil {
		return iv.handleError(err, "Parsing Data Failed!")
	}

	err = iv.validateJsonData()
	if err != nil {
		iv.pdfData = invoicePdfData{}
		return iv.handleError(err, "Incorrect data!")
	}

	return nil
}
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
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

	pdfGen := generator.NewPDFGenerator(generator.MetaData{
		LineHeight:   5,
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
	pdfGen.PrintPdfText("Firmen Name Gmbh, Paulaner-Str. 99, 04109 Leipzig", "", "L")
	pdfGen.SetFontSize(defaultFontSize)

	//Anschrift Empfänger
	pdfGen.SetCursor(iv.marginLeft, 56)
	pdfGen.PrintLnPdfText("Firmen Name Gmbh", "", "L")
	pdfGen.PrintLnPdfText("Frau Musterfrau", "", "L")
	pdfGen.PrintLnPdfText("Paulaner-Str. 99", "", "L")
	pdfGen.PrintLnPdfText("Str. Zusatz", "", "L")
	pdfGen.PrintLnPdfText("04109 Leipzig", "", "L")

	//MetaData Infos Rechnung in 2 Spalten
	pdfGen.SetCursor(iv.marginLeft+100, 93)
	pdfGen.PrintLnPdfText("Kundennummer:", "", "L")
	pdfGen.PrintLnPdfText("Rechnungsnummer:", "", "L")
	pdfGen.PrintLnPdfText("Datum:", "", "L")

	pdfGen.SetCursor(iv.marginLeft+140, 93)
	pdfGen.PrintLnPdfText("KD83383", "", "L")
	pdfGen.PrintLnPdfText("RE20230002", "", "L")
	pdfGen.PrintLnPdfText("23.04.2023", "", "L")

	//Überschrift
	pdfGen.DrawLine(iv.marginLeft, 108, pageWidth-iv.marginRight, 108, lineColor, 0)
	pdfGen.SetCursor(iv.marginLeft, 112)
	pdfGen.SetFontSize(headerFontSize)
	pdfGen.PrintLnPdfText("Rechnung - 4", "b", "L")
	pdfGen.SetFontSize(defaultFontSize)
	pdfGen.NewLine(pdfGen.GetMarginLeft())

	pdfGen.PrintLnPdfText("Sehr geehrter Herr Mustermann,", "", "L")
	pdfGen.NewLine(pdfGen.GetMarginLeft())

	pdfGen.PrintLnPdfText("lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor \n     invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.", "", "L")
	pdfGen.NewLine(pdfGen.GetMarginLeft())

	pdfGen.SetFontSize(smallFontSize)
	pdfGen.PrintLnPdfText("Leistungszeitraum: 01.01.1999 - 01.01.2023", "i", "L")
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
	}

	for _, product := range productList {
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
	//var bodyText = [][]string{
	//	{"1", "40,50", "45,00€", "h", "agiles Software-Testing, System-Monitoring, \n Programmierung", "19%", "2.000,00€"},
	//	{"1", "50,00", "40,00€", "h", "Softwareentwicklung", "19%", "1.000,00€"},
	//}
	var bodyCellAlign = []string{"LM", "LM", "LM", "LM", "RM", "RM"}
	var summaryCells = [][]string{
		{"", "Zwischensumme", "2.000,00€"},
		{"", "USt. 19%", "0€"},
		{"", "USt. 7%", "0€"},
		{"", "Gesamtbetrag", "2.000,00€"},
	}
	var summaryColumnWidths = []float64{getCellWith(60), getCellWith(25), getCellWith(15)}
	var summaryCellAlign = []string{"LM", "LM", "RM"}

	pdfGen.PrintTableHeader(headerCells, columnWidth)
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

	return pdfGen.GetPdf(), pdfGen.GetError()
}

func (iv *Invoice) handleError(err error, msg string) (responseErr error) {
	iv.logger.Error().Msgf(err.Error())
	return fmt.Errorf("ERROR: %s", msg)
}
