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

func (iv *Invoice) GeneratePDF() (pdf *gofpdf.Fpdf, err error) {

	lineColor := color{200, 200, 200}
	iv.logger.Debug().Msg("Endpoint Hit: pdfPage")

	iv.newPDF()
	iv.writePdfText("TEST", "", iv.textSize, "L")

	err = iv.placeImgOnPosXY("https://cdn.pictro.de/logosIcons/stack-one_logo_vector_white_small.png", 153, 20)

	iv.pdf.SetXY(25, 51)

	iv.writePdfText("Firmen Name Gmbh, Paulaner-Str. 99, 04109 Leipzig", "", iv.textSizeSmall, "R")

	iv.writePdfText("Firmen Name Gmbh", "", iv.textSize, "L")
	iv.writePdfText("Frau Musterfrau", "", iv.textSize, "L")
	iv.writePdfText("Paulaner-Str. 99", "", iv.textSize, "L")
	iv.writePdfText("04109 Leipzig", "", iv.textSize, "L")

	iv.drawLine(25, 94, 186, 94, lineColor)

	iv.pdf.SetXY(25, 96)
	iv.writePdfText("Rechnung - 4", "b", 16, "L")

	pageWith, _ := iv.pdf.GetPageSize()
	iv.pdf.SetXY(185, pageWith)

	iv.writePdfText("Kundennummer: KD83383", "", 11, "R")
	iv.writePdfText("Rechnungsnummer: RE20230002", "", 11, "R")
	iv.writePdfText("Dateum: 23.04.2023", "", 11, "R")

	iv.drawLine(25, 111, 186, 111, lineColor)

	iv.pdf.SetXY(25, 117)
	iv.writePdfText("ciĝas ĉe paĝo Vielen Dank für Ihr Vertrauen!\nHiermit stelle ich Ihnen die folgenden Positionen in Rechnung.", "", 11, "L")

	return iv.pdf, iv.pdf.Error()
}

func (iv *Invoice) handleError(err error, msg string) (responseErr error) {
	iv.logger.Error().Msgf(err.Error())
	return fmt.Errorf("ERROR: %s", msg)
}
