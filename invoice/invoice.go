package invoice

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog"
	"io"
	"time"
)

type Invoice struct {
	pdfData  invoicePdfData
	pdf      *gofpdf.Fpdf
	logger   *zerolog.Logger
	textFont string
}

type invoicePdfData struct {
	SenderAddress   addressInfo `json:"senderAddress"`
	ReceiverAddress addressInfo `json:"receiverAddress"`
	SenderInfo      senderInfo  `json:"senderInfo"`
	InvoiceMeta     invoiceMeta `json:"invoiceMeta"`
	InvoiceBody     invoiceBody `json:"InvoiceBody"`
}

type invoiceBody struct {
	OpeningText     string          `json:"openingText"`
	ServiceTimeText string          `json:"serviceTimeText"`
	HeadlineText    string          `json:"headlineText"`
	ClosingText     string          `json:"closingText"`
	UstNotice       string          `json:"ustNotice"`
	InvoicedItems   []invoicedItems `json:"invoicedItems"`
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

type addressInfo struct {
	FullForename string  `json:"fullForename"`
	FullSurname  string  `json:"fullSurname"`
	NameSecond   string  `json:"nameSecond"`
	Address      address `json:"address"`
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
	InvoiceNumber  string    `json:"invoiceNumber"`
	InvoiceDate    time.Time `json:"invoiceDate"`
	CustomerNumber string    `json:"customerNumber"`
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
		logger:   logger,
		textFont: "Arial",
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

func (iv *Invoice) GeneratePDF() (pdf *gofpdf.Fpdf, err error) {

	iv.logger.Debug().Msg("Endpoint Hit: pdfPage")

	iv.newPDF()
	iv.setPdfText(100, 100, "LOGO", "b", 16, 40, 10)

	err = iv.placeImgOnPosXY("https://cdn.pictro.de/logosIcons/stack-one_logo_vector_white_small.png", 100, 20)

	return iv.pdf, err
}

func (iv *Invoice) handleError(err error, msg string) (responseErr error) {
	iv.logger.Error().Msgf(err.Error())
	return fmt.Errorf("ERROR: %s", msg)
}
