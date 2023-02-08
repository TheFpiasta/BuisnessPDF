package invoice

import (
	"SimpleInvoice/generator"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog"
	"io"
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
		logger:      logger,
		textFont:    "openSans",
		marginLeft:  25,
		marginRight: 20,
		marginTop:   45,
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

func (iv *Invoice) GeneratePDF() (*gofpdf.Fpdf, error) {

	iv.logger.Debug().Msg("Endpoint Hit: pdfPage")

	lineColor := generator.Color{R: 200, G: 200, B: 200}
	const defaultFontSize = 11
	const smallFontSize = 8
	const headerFontSize = 16

	pdfGen := generator.NewPDFGenerator(generator.MetaData{
		LineHeight:  5,
		FontName:    "openSans",
		FontGapY:    1.5,
		FontSize:    defaultFontSize,
		MarginLeft:  iv.marginLeft,
		MarginTop:   iv.marginTop,
		MarginRight: iv.marginRight,
		Unit:        "mm",
	})

	pageWidth, _ := pdfGen.GetPdf().GetPageSize()

	err := pdfGen.PlaceImgOnPosXY("https://cdn.pictro.de/logosIcons/stack-one_logo_vector_white_small.png", 153, 20)
	if err != nil {
		return iv.pdf, err
	}

	//Anschrift Empfänger
	pdfGen.SetCursor(pageWidth-iv.marginRight, 61)

	pdfGen.PrintLnPdfText("Mein Name Gmbh", "", "R")
	pdfGen.PrintLnPdfText("Meine Paulaner-Str. 99", "", "R")
	pdfGen.PrintLnPdfText("Meine Str Zusatz", "", "R")
	pdfGen.PrintLnPdfText("04109 Leipzig", "", "R")

	//Anschrift Sender small
	pdfGen.SetCursor(iv.marginLeft, 61)
	pdfGen.SetFontSize(smallFontSize)
	pdfGen.PrintPdfText("Firmen Name Gmbh, Paulaner-Str. 99, 04109 Leipzig", "", "L")
	pdfGen.SetFontSize(defaultFontSize)

	//Anschrift Sender
	pdfGen.SetCursor(iv.marginLeft, 70)
	pdfGen.PrintLnPdfText("Firmen Name Gmbh", "", "L")
	pdfGen.PrintLnPdfText("Frau Musterfrau", "", "L")
	pdfGen.PrintLnPdfText("Paulaner-Str. 99", "", "L")
	pdfGen.PrintLnPdfText("Str. Zusatz", "", "L")
	pdfGen.PrintLnPdfText("04109 Leipzig", "", "L")

	//MetaData Infos Rechnung in 2 Spalten
	pdfGen.SetCursor(iv.marginLeft+100, 100)
	pdfGen.PrintLnPdfText("Kundennummer:", "", "L")
	pdfGen.PrintLnPdfText("Rechnungsnummer:", "", "L")
	pdfGen.PrintLnPdfText("Datum:", "", "L")

	pdfGen.SetCursor(iv.marginLeft+140, 100)
	pdfGen.PrintLnPdfText("KD83383", "", "L")
	pdfGen.PrintLnPdfText("RE20230002", "", "L")
	pdfGen.PrintLnPdfText("23.04.2023", "", "L")

	//Überschrift
	pdfGen.DrawLine(iv.marginLeft, 120, pageWidth-iv.marginRight, 120, lineColor, 0)
	pdfGen.SetCursor(iv.marginLeft, 122)
	pdfGen.SetFontSize(headerFontSize)
	pdfGen.PrintPdfText("Rechnung - 4", "b", "L")
	pdfGen.SetFontSize(defaultFontSize)

	//Tabelle
	//type InvoiceItemData struct {
	//	PositionNumber uint
	//	Quantity       float64
	//	Unit           string
	//	Description    string
	//	SinglePrice    float64
	//	NetPrice       float64
	//}

	getCellWith := func(percent float64) float64 {
		maxSavePrintingWidth, _ := pdfGen.GetPdf().GetPageSize()
		maxSavePrintingWidth = maxSavePrintingWidth - pdfGen.GetMarginLeft() - pdfGen.GetMarginRight()

		return (percent * maxSavePrintingWidth) / 100.0
	}

	pdfGen.SetCursor(iv.marginLeft, 200)
	pdfGen.PrintInvoiceTable(
		[]string{"Position", "Anzahl", "Beschreibung", "USt", "Einzelpreis", "Netto"},
		[]float64{getCellWith(11), getCellWith(11), getCellWith(40), getCellWith(8), getCellWith(15), getCellWith(15)},
		[][]string{
			{"1", "50,00", "Softwareentwicklung", "0%", "40,00€", "2.000,00€"},
			{"2", "25,00", "agiles Software-Testing,\n System-Monitoring", "0%", "30,00€", "750,00€"},
		},
		[][2]string{
			{"Zwischensumme", "2.000,00€"},
			{"USt. 19%", "0€"},
			{"USt. 7%", "0€"},
			{"Gesamtbetrag", "2.000,00€"},
		},
		[3]float64{getCellWith(60), getCellWith(25), getCellWith(15)},
		[]string{"LM", "LM", "LM", "LM", "RM", "RM"},
	)

	return pdfGen.GetPdf(), pdfGen.GetError()
}

func (iv *Invoice) handleError(err error, msg string) (responseErr error) {
	iv.logger.Error().Msgf(err.Error())
	return fmt.Errorf("ERROR: %s", msg)
}
