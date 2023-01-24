package invoice

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog"
	"io"
	"time"
)

type Invoice struct {
	pdfData  invoicePdf
	pdf      *gofpdf.Fpdf
	logger   *zerolog.Logger
	textFont string
}

type invoicePdf struct {
	senderInfo      addressInfo
	receiverInfo    addressInfo
	invoiceMeta     invoiceMeta
	logoURL         string
	openingText     string
	serviceTimeText string
	headline        string
	closingText     string
	positions       []service
}

type address struct {
	Addressee string `json:"addressee"`
	ZipCode   string `json:"zipCode"`
	CityName  string `json:"cityName"`
}

type addressInfo struct {
	Name       string  `json:""`
	NameSecond string  `json:""`
	Address    address `json:""`
	Phone      string  `json:""`
	Email      string  `json:""`
}

type invoiceMeta struct {
	invoiceNumber  string
	invoiceDate    time.Time
	customerNumber string
}

type service struct {
	position int
	name     string
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
		iv.pdfData = invoicePdf{}
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
