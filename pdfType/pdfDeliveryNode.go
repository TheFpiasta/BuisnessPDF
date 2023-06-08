package pdfType

import (
	"SimpleInvoice/generator"
	"encoding/json"
	errorsWithStack "github.com/go-errors/errors"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog"
	"io"
	"net/http"
)

type DeliveryNode struct {
	data          deliveryNodeRequestData
	meta          PdfMeta
	logger        *zerolog.Logger
	printErrStack bool
	pdfGen        *generator.PDFGenerator
}

type deliveryNodeRequestData struct {
	SenderAddress   FullPersonInfo `json:"senderAddress"`
	ReceiverAddress FullPersonInfo `json:"receiverAddress"`
	SenderInfo      SenderInfo     `json:"senderInfo"`
	DeliveryMeta    struct {
		DeliveryNodeNumber string `json:"deliveryNodeNumber"`
		DeliveryDate       string `json:"deliveryDate"`
		CustomerNumber     string `json:"customerNumber"`
	} `json:"deliveryMeta"`
	DeliveryNodeTexts struct {
		OpeningText  string `json:"openingText"`
		HeadlineText string `json:"headlineText"`
		ClosingText  string `json:"closingText"`
		Agb          string `json:"agb"`
	} `json:"deliveryNodeTexts"`
	DeliveryItems []struct {
		PositionNumber string  `json:"positionNumber"`
		Quantity       float64 `json:"quantity"`
		Description    string  `json:"description"`
		Unit           string  `json:"unit"`
	} `json:"deliveryItems"`
}

func NewDeliveryNode(logger *zerolog.Logger) *DeliveryNode {
	return &DeliveryNode{
		data: deliveryNodeRequestData{},
		meta: PdfMeta{
			Margin: pdfMargin{
				Left:   25,
				Right:  20,
				Top:    45,
				Bottom: 0,
			},
			Font: pdfFont{
				FontName:    "openSans",
				SizeDefault: 10,
				SizeSmall:   8,
				SizeLarge:   15,
			},
		},
		logger:        logger,
		printErrStack: logger.GetLevel() <= zerolog.DebugLevel,
	}
}

func (d *DeliveryNode) SetDataFromRequest(request *http.Request) (err error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			d.LogError(err)
		}
	}(request.Body)

	err = json.NewDecoder(request.Body).Decode(&d.data)
	if err != nil {
		return err
	}

	err = d.validateData()
	if err != nil {
		d.data = deliveryNodeRequestData{}
		return err
	}

	return nil
}

func (d *DeliveryNode) validateData() (err error) {
	return err
}

func (d *DeliveryNode) LogError(err error) {
	var errStr string

	if _, ok := err.(*errorsWithStack.Error); ok && d.printErrStack {
		errStr = err.(*errorsWithStack.Error).ErrorStack()
	} else {
		errStr = err.Error()
	}

	d.logger.Error().Msgf(errStr)
}

func (d *DeliveryNode) GeneratePDF() (*gofpdf.Fpdf, error) {
	d.logger.Debug().Msg("generate delivery node")

	lineColor := generator.Color{R: 200, G: 200, B: 200}

	pdfGen, err := generator.NewPDFGenerator(generator.MetaData{
		FontName:     "OpenSans",
		FontGapY:     1.3,
		FontSize:     d.meta.Font.SizeDefault,
		MarginLeft:   d.meta.Margin.Left,
		MarginTop:    d.meta.Margin.Top,
		MarginRight:  d.meta.Margin.Right,
		MarginBottom: d.meta.Margin.Bottom,
		Unit:         "mm",
	}, false, d.logger)

	if err != nil {
		return nil, err
	}

	d.pdfGen = pdfGen

	if d.data.SenderInfo.MimeLogoUrl != "" {
		d.printMimeImg()
	}

	d.printAddressee(lineColor)
	d.printMetaData(pdfGen, lineColor)
	d.printHeadlineAndOpeningText(pdfGen)
	//i.printInvoiceTable(pdfGen)
	//TODO unterschriften einfügen
	d.printClosingText(pdfGen)
	letterFooter(d.pdfGen, d.meta, d.data.SenderInfo, d.data.SenderAddress, lineColor)

	return pdfGen.GetPdf(), pdfGen.GetError()
}

func (d *DeliveryNode) printMimeImg() {
	pageWidth, _ := d.pdfGen.GetPdf().GetPageSize()
	mimeImg(d.pdfGen, d.data.SenderInfo.MimeLogoUrl, pageWidth-d.meta.Margin.Right, 15, d.data.SenderInfo.MimeLogoScale)
}

func (d *DeliveryNode) printAddressee(lineColor generator.Color) {
	//pageWidth, _ := i.pdfGen.GetPdf().GetPageSize()
	//i.pdfGen.DrawLine(i.meta.Margin.Left, i.meta.Margin.Top, pageWidth-i.meta.Margin.Right, i.meta.Margin.Top, lineColor, 0)

	letterAddressSenderSmall(d.pdfGen, getAddressLine(d.data.SenderAddress), d.meta.Margin.Left, 49, d.meta.Font.SizeSmall)
	d.pdfGen.SetFontSize(d.meta.Font.SizeDefault)

	letterReceiverAddress(d.pdfGen, d.data.ReceiverAddress, d.meta.Margin.Left, 56)
}

func (d *DeliveryNode) printMetaData(pdfGen *generator.PDFGenerator, lineColor generator.Color) {
	pdfGen.SetFontSize(d.meta.Font.SizeDefault)
	pdfGen.DrawLine(d.meta.Margin.Left+98, 56, d.meta.Margin.Left+98, 80, lineColor, 0)
	pdfGen.SetCursor(d.meta.Margin.Left+100, 56)
	pdfGen.PrintLnPdfText("Kundennummer:", "", "L")
	pdfGen.PrintLnPdfText("Liefernummer:", "", "L")
	pdfGen.PrintLnPdfText("Datum:", "", "L")

	pdfGen.SetCursor(d.meta.Margin.Left+140, 56)
	pdfGen.PrintLnPdfText(d.data.DeliveryMeta.CustomerNumber, "", "L")
	pdfGen.PrintLnPdfText(d.data.DeliveryMeta.DeliveryNodeNumber, "", "L")
	pdfGen.PrintLnPdfText(d.data.DeliveryMeta.DeliveryDate, "", "L")
}

func (d *DeliveryNode) printHeadlineAndOpeningText(pdfGen *generator.PDFGenerator) {
	//Überschrift
	pdfGen.SetCursor(d.meta.Margin.Left, 100)
	pdfGen.SetFontSize(d.meta.Font.SizeLarge)
	pdfGen.PrintLnPdfText(d.data.DeliveryNodeTexts.HeadlineText+" "+d.data.DeliveryMeta.DeliveryNodeNumber, "b", "L")

	//opening
	pdfGen.SetFontSize(d.meta.Font.SizeDefault)
	pdfGen.NewLine(pdfGen.GetMarginLeft())
	pdfGen.PrintLnPdfText(d.data.DeliveryNodeTexts.OpeningText, "", "L")
}

//TODO implelment table

func (d *DeliveryNode) printClosingText(pdfGen *generator.PDFGenerator) {
	pdfGen.SetFontSize(d.meta.Font.SizeDefault)
	pdfGen.NewLine(d.meta.Margin.Left)
	pdfGen.NewLine(d.meta.Margin.Left)
	pdfGen.NewLine(d.meta.Margin.Left)
	pdfGen.PrintLnPdfText(d.data.DeliveryNodeTexts.ClosingText, "", "L")
	pdfGen.NewLine(d.meta.Margin.Left)
	pdfGen.NewLine(d.meta.Margin.Left)
	pdfGen.PrintLnPdfText(d.data.DeliveryNodeTexts.Agb, "", "L")
}
