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
	data             deliveryNodeRequestData
	meta             PdfMeta
	logger           *zerolog.Logger
	printErrStack    bool
	pdfGen           *generator.PDFGenerator
	defaultLineColor generator.Color
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
		logger:           logger,
		printErrStack:    logger.GetLevel() <= zerolog.DebugLevel,
		defaultLineColor: generator.Color{R: 200, G: 200, B: 200},
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

	pdfGen, err := generator.NewPDFGenerator(generator.MetaData{
		FontName:     "OpenSans",
		FontGapY:     1.3,
		FontSize:     d.meta.Font.SizeDefault,
		MarginLeft:   d.meta.Margin.Left,
		MarginTop:    d.meta.Margin.Top,
		MarginRight:  d.meta.Margin.Right,
		MarginBottom: d.meta.Margin.Bottom,
		Unit:         "mm",
	}, false, d.logger, func() {

	}, func(isLastPage bool) {

	})

	if err != nil {
		return nil, err
	}

	d.pdfGen = pdfGen
	d.pdfGen.NextPage()

	if d.data.SenderInfo.MimeLogoUrl != "" {
		d.printMimeImg()
	}

	d.printAddressee()
	d.printMetaData(pdfGen)
	d.printHeadlineAndOpeningText(pdfGen)
	d.printDeliveryTable(pdfGen)
	d.printClosingText(pdfGen)
	d.printSignatureSection(pdfGen)
	d.printFooter(pdfGen)

	return pdfGen.GetPdf(), pdfGen.GetError()
}

func (d *DeliveryNode) printMimeImg() {
	pageWidth, _ := d.pdfGen.GetPdf().GetPageSize()
	mimeImg(d.pdfGen, d.data.SenderInfo.MimeLogoUrl, pageWidth-d.meta.Margin.Right, 15, d.data.SenderInfo.MimeLogoScale)
}

func (d *DeliveryNode) printAddressee() {
	//pageWidth, _ := i.pdfGen.GetPdf().GetPageSize()
	//i.pdfGen.DrawLine(i.meta.Margin.Left, i.meta.Margin.Top, pageWidth-i.meta.Margin.Right, i.meta.Margin.Top, lineColor, 0)

	letterAddressSenderSmall(d.pdfGen, getAddressLine(d.data.SenderAddress), d.meta.Margin.Left, 49, d.meta.Font.SizeSmall)
	d.pdfGen.SetFontSize(d.meta.Font.SizeDefault)

	letterReceiverAddress(d.pdfGen, d.data.ReceiverAddress, d.meta.Margin.Left, 56)
}

func (d *DeliveryNode) printMetaData(pdfGen *generator.PDFGenerator) {
	pdfGen.SetFontSize(d.meta.Font.SizeDefault)
	pdfGen.DrawLine(d.meta.Margin.Left+98, 56, d.meta.Margin.Left+98, 80, d.defaultLineColor, 0)
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

func (d *DeliveryNode) printDeliveryTable(pdfGen *generator.PDFGenerator) {
	var items = [][]string{{}}

	for _, item := range d.data.DeliveryItems {
		items = append(items,
			[]string{
				item.PositionNumber,
				germanNumber(int(item.Quantity)) + " " + item.Unit,
				item.Description,
				"",
			})
	}

	var headerCells = []string{"Pos", "Anzahl", "Beschreibung", "Notiz"}
	var columnPercent = []float64{7, 18, 40, 35}
	var columnWidth = getColumnWithFromPercentage(pdfGen, columnPercent)
	var headerCellAlign = []string{"LM", "LM", "LM", "LM"}
	var bodyCellAlign = []string{"LM", "LM", "LM", "LM"}

	pdfGen.NewLine(d.meta.Margin.Left)
	pdfGen.SetFontSize(d.meta.Font.SizeDefault)
	pdfGen.PrintTableHeader(headerCells, columnWidth, headerCellAlign)
	pdfGen.PrintTableBody(items, columnWidth, bodyCellAlign)

}

func (d *DeliveryNode) printClosingText(pdfGen *generator.PDFGenerator) {
	pdfGen.SetFontSize(d.meta.Font.SizeDefault)
	pdfGen.NewLine(d.meta.Margin.Left)
	pdfGen.NewLine(d.meta.Margin.Left)
	pdfGen.NewLine(d.meta.Margin.Left)
	pdfGen.PrintLnPdfText(d.data.DeliveryNodeTexts.Agb, "", "L")
	pdfGen.NewLine(d.meta.Margin.Left)
	pdfGen.NewLine(d.meta.Margin.Left)
	pdfGen.PrintLnPdfText(d.data.DeliveryNodeTexts.ClosingText, "", "L")
}

func (d *DeliveryNode) printSignatureSection(pdfGen *generator.PDFGenerator) {
	const startSignatureSectionOnPosY = 230

	pageWidth, _ := d.pdfGen.GetPdf().GetPageSize()
	var startSupplierX = d.meta.Margin.Left
	var startCustomerX = pageWidth / 2.0
	var senderSignatureName string

	if d.data.SenderAddress.CompanyName != "" {
		senderSignatureName = d.data.SenderAddress.CompanyName
	} else {
		senderSignatureName = "Lieferant"
	}

	d.printSignaturePart(pdfGen, senderSignatureName, startSupplierX, startSignatureSectionOnPosY, d.defaultLineColor)
	d.printSignaturePart(pdfGen, "Kunde", startCustomerX, startSignatureSectionOnPosY, d.defaultLineColor)

}

func (d *DeliveryNode) printSignaturePart(pdfGen *generator.PDFGenerator, headText string, startX float64, startY float64, lineColor generator.Color) {
	const nameLength = 60
	const dateLength = 20
	const gabLength = 5
	const signatureLength = 35

	var cY float64

	pdfGen.SetCursor(startX, startY)
	pdfGen.DrawLine(startX, startY, startX+nameLength, startY, lineColor, 0)
	_, cY = pdfGen.GetCursor()
	pdfGen.SetCursor(startX, cY+1)
	pdfGen.SetFontSize(d.meta.Font.SizeSmall)
	pdfGen.PrintPdfText(headText, "b", "L")
	pdfGen.PrintPdfText("(Name)", "", "L")
	pdfGen.NewLine(startX)
	pdfGen.SetFontSize(d.meta.Font.SizeDefault)

	pdfGen.NewLine(startX)
	pdfGen.NewLine(startX)
	pdfGen.NewLine(startX)

	_, cY = pdfGen.GetCursor()
	var dateEndX = startX + dateLength
	var signatureStartX = startX + dateLength + gabLength
	var signatureEndX = startX + dateLength + gabLength + signatureLength
	pdfGen.DrawLine(startX, cY, dateEndX, cY, lineColor, 0)
	pdfGen.DrawLine(signatureStartX, cY, signatureEndX, cY, lineColor, 0)
	_, cY = pdfGen.GetCursor()
	pdfGen.SetCursor(startX, cY+1)
	pdfGen.SetFontSize(d.meta.Font.SizeSmall)
	pdfGen.PrintPdfText("Datum", "", "L")
	_, cY = pdfGen.GetCursor()
	pdfGen.SetCursor(signatureStartX, cY)
	pdfGen.PrintPdfText("Unterschrift", "", "L")
}

func (d *DeliveryNode) printFooter(pdfGen *generator.PDFGenerator) {
	const startAtY = 273
	const startPageNumberY = 282
	const gabY = 3

	pageWidth, _ := pdfGen.GetPdf().GetPageSize()

	pdfGen.SetFontSize(d.meta.Font.SizeSmall)
	pdfGen.DrawLine(d.meta.Margin.Left, startAtY, pageWidth-d.meta.Margin.Right, startAtY, d.defaultLineColor, 0)

	pdfGen.SetCursor(d.meta.Margin.Left, startAtY+gabY)
	//pdfGen.PrintLnPdfText("Web", "", "L")
	pdfGen.PrintPdfText(d.data.SenderInfo.Web, "", "L")

	pdfGen.SetCursor(pageWidth/2, startAtY+gabY)
	//pdfGen.PrintLnPdfText("Tel", "", "C")
	pdfGen.PrintPdfText(d.data.SenderInfo.Phone, "", "C")

	pdfGen.SetCursor(pageWidth-d.meta.Margin.Right, startAtY+gabY)
	//pdfGen.PrintLnPdfText("E-Mail", "", "R")
	pdfGen.PrintPdfText(d.data.SenderInfo.Email, "", "R")

	pdfGen.DrawLine(d.meta.Margin.Left, startPageNumberY, pageWidth-d.meta.Margin.Right, startPageNumberY, d.defaultLineColor, 0)
	pdfGen.SetCursor(pageWidth/2, startPageNumberY+gabY)
	pdfGen.PrintLnPdfText("Seite 1 von 1", "", "C")
	pdfGen.SetFontSize(d.meta.Font.SizeDefault)
}
