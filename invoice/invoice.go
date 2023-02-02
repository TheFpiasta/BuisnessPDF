package invoice

import (
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
}

type invoicedItems struct {
	PositionNumber    int     `json:"positionNumber"`
	Quantity          float64 `json:"quantity"`
	Unit              string  `json:"unit"`
	Description       string  `json:"description"`
	SinglePrice       float64 `json:"singlePrice"`
	SinglePriceNet    float64 `json:"singlePriceNet"`
	OverallPriceNet   float64 `json:"overallPriceNet"`
	OverallPriceGross float64 `json:"overallPriceGross"`
	OverallTaxes      float64 `json:"overallTaxes"`
	TaxesPercentage   float64 `json:"taxesPercentage"`
	Currency          string  `json:"currency"`
}

func New(logger *zerolog.Logger) (iv *Invoice) {
	iv = &Invoice{
		logger:        logger,
		textFont:      "openSans",
		lineHeight:    5,
		textSize:      11,
		textSizeSmall: 8,
		marginLeft:    25,
		marginRight:   20,
		marginTop:     45,
		fontGapY:      2,
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

type color struct {
	r uint8
	g uint8
	b uint8
}

func (iv *Invoice) GeneratePDF() (*gofpdf.Fpdf, error) {

	lineColor := color{200, 200, 200}

	iv.logger.Debug().Msg("Endpoint Hit: pdfPage")

	iv.newPDF()
	pageWidth, _ := iv.pdf.GetPageSize()

	err := iv.placeImgOnPosXY("https://cdn.pictro.de/logosIcons/stack-one_logo_vector_white_small.png", 153, 20)
	if err != nil {
		return iv.pdf, err
	}

	//Anschrift Empfänger
	iv.pdf.SetXY(pageWidth-iv.marginRight, 61)
	iv.printLnPdfText("Mein Name Gmbh", "", iv.textSize, "R")
	iv.printLnPdfText("Meine Paulaner-Str. 99", "", iv.textSize, "R")
	iv.printLnPdfText("Meine Str Zusatz", "", iv.textSize, "R")
	iv.printLnPdfText("04109 Leipzig", "", iv.textSize, "R")

	//Anschrift Sender small
	iv.pdf.SetXY(iv.marginLeft, 61)
	iv.printPdfText("Firmen Name Gmbh, Paulaner-Str. 99, 04109 Leipzig", "", iv.textSizeSmall, "L")

	//Anschrift Sender
	iv.pdf.SetXY(iv.marginLeft, 70)
	iv.printLnPdfText("Firmen Name Gmbh", "", iv.textSize, "L")
	iv.printLnPdfText("Frau Musterfrau", "", iv.textSize, "L")
	iv.printLnPdfText("Paulaner-Str. 99", "", iv.textSize, "L")
	iv.printLnPdfText("Str. Zusatz", "", iv.textSize, "L")
	iv.printLnPdfText("04109 Leipzig", "", iv.textSize, "L")

	//Meta Infos Rechnung in 2 Spalten
	iv.pdf.SetXY(iv.marginLeft+100, 100)
	iv.printLnPdfText("Kundennummer:", "", 11, "L")
	iv.printLnPdfText("Rechnungsnummer:", "", 11, "L")
	iv.printLnPdfText("Datum:", "", 11, "L")

	iv.pdf.SetXY(iv.marginLeft+140, 100)
	iv.printLnPdfText("KD83383", "", 11, "L")
	iv.printLnPdfText("RE20230002", "", 11, "L")
	iv.printLnPdfText("23.04.2023", "", 11, "L")

	//Überschrift
	iv.drawLine(iv.marginLeft, 120, pageWidth-iv.marginRight, 120, lineColor)
	iv.pdf.SetXY(iv.marginLeft, 122)
	iv.printPdfText("Rechnung - 4", "b", 16, "L")

	//Tabelle
	iv.pdf.SetFillColor(200, 200, 200)
	const colNumber = 5
	header := [colNumber]string{"No", "Description", "Quantity", "Unit Price ($)", "Price ($)"}
	colWidth := [colNumber]float64{10.0, 50.0, 40.0, 30.0, 30.0}
	lineHt := 10.0
	iv.pdf.SetXY(iv.marginLeft, iv.pdf.GetY()+iv.lineHeight+10.0)
	for colJ := 0; colJ < colNumber; colJ++ {
		iv.pdf.CellFormat(colWidth[colJ], lineHt, header[colJ], "1", 0, "CM", true, 0, "")
	}

	return iv.pdf, iv.pdf.Error()
}

func (iv *Invoice) handleError(err error, msg string) (responseErr error) {
	iv.logger.Error().Msgf(err.Error())
	return fmt.Errorf("ERROR: %s", msg)
}
