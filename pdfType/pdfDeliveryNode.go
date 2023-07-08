package pdfType

import (
	"SimpleInvoice/generator"
	din5008a "SimpleInvoice/norms/letter/din-5008-a"
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
	footerStartY  float64
}

type deliveryNodeRequestData struct {
	SenderAddress   din5008a.FullAdresse `json:"senderAddress"`
	ReceiverAddress din5008a.FullAdresse `json:"receiverAddress"`
	SenderInfo      SenderInfo           `json:"senderInfo"`
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
				Top:    din5008a.AddressSenderTextStartY,
				Bottom: 0,
			},
			Font: pdfFont{
				FontName:    "openSans",
				SizeDefault: din5008a.FontSize10,
				SizeSmall:   din5008a.FontSizeSender8,
				SizeLarge:   din5008a.FontSize10 + 5,
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
	//TODO implement
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

	pdfGen, err := generator.NewPDFGenerator(
		generator.MetaData{
			FontName:         "OpenSans",
			FontGapY:         1.3,
			FontSize:         d.meta.Font.SizeDefault,
			MarginLeft:       d.meta.Margin.Left,
			MarginTop:        d.meta.Margin.Top,
			MarginRight:      d.meta.Margin.Right,
			MarginBottom:     d.meta.Margin.Bottom,
			Unit:             "mm",
			DefaultLineWidth: 0.4,
			DefaultLineColor: generator.Color{R: 162, G: 162, B: 162},
		},
		false,
		d.logger,
		func() {
			d.printHeader()
		},
		func(isLastPage bool) {
			d.printFooter()
		},
	)

	if err != nil {
		return nil, err
	}

	d.pdfGen = pdfGen
	d.pdfGen.NewPage()

	d.doGeneratePdf()

	return d.pdfGen.GetPdf(), d.pdfGen.GetError()
}

func (d *DeliveryNode) doGeneratePdf() {
	var infoData []din5008a.InfoData
	infoData = append(infoData, din5008a.InfoData{Name: "Kundennummer:", Value: d.data.DeliveryMeta.CustomerNumber})
	infoData = append(infoData, din5008a.InfoData{Name: "Liefernummer:", Value: d.data.DeliveryMeta.DeliveryNodeNumber})
	infoData = append(infoData, din5008a.InfoData{Name: "Datum:", Value: d.data.DeliveryMeta.DeliveryDate})

	din5008a.FullAddressesAndInfoPart(d.pdfGen, d.data.SenderAddress, d.data.ReceiverAddress, infoData)

	din5008a.Body(d.pdfGen, func() {
		d.printHeadlineAndOpeningText()
		d.printDeliveryTable()
		d.printClosingText()
		d.printSignatureSection()
	})

	din5008a.PageNumbering(d.pdfGen, d.footerStartY)
}

func (d *DeliveryNode) printHeadlineAndOpeningText() {
	//Überschrift
	d.pdfGen.SetFontSize(d.meta.Font.SizeLarge)
	d.pdfGen.PrintLnPdfText(d.data.DeliveryNodeTexts.HeadlineText+" "+d.data.DeliveryMeta.DeliveryNodeNumber, "b", "L")

	//opening
	d.pdfGen.SetFontSize(din5008a.FontSize10)
	d.pdfGen.SetFontGapY(din5008a.FontGab10)
	d.pdfGen.NewLine(din5008a.BodyStartX)
	d.pdfGen.PrintLnPdfText(d.data.DeliveryNodeTexts.OpeningText, "", "L")
}

func (d *DeliveryNode) printDeliveryTable() {
	var items = [][]string{{}}

	for _, item := range d.data.DeliveryItems {
		items = append(items,
			[]string{
				item.PositionNumber,
				germanNumber(int(item.Quantity)) + " " + item.Unit,
				item.Description,
				"",
			},
		)
	}

	var headerCells = []string{"Pos", "Anzahl", "Beschreibung", "Notiz"}
	var columnPercent = []float64{7, 18, 40, 35}
	var columnWidth = getColumnWithFromPercentage(d.pdfGen, columnPercent)
	var headerCellAlign = []string{"LM", "LM", "LM", "LM"}
	var bodyCellAlign = []string{"LM", "LM", "LM", "LM"}

	d.pdfGen.NewLine(din5008a.BodyStartX)
	d.pdfGen.SetFontSize(din5008a.FontSize10)
	d.pdfGen.SetFontGapY(din5008a.FontGab10)
	d.pdfGen.PrintTableHeader(headerCells, columnWidth, headerCellAlign)
	d.pdfGen.PrintTableBody(items, columnWidth, bodyCellAlign)

}

func (d *DeliveryNode) printClosingText() {
	d.pdfGen.SetFontSize(din5008a.FontSize10)
	d.pdfGen.SetFontGapY(din5008a.FontGab10)
	d.pdfGen.NewLine(din5008a.BodyStartX)
	d.pdfGen.PrintLnPdfText(d.data.DeliveryNodeTexts.Agb, "", "L")
	d.pdfGen.NewLine(din5008a.BodyStartX)
	d.pdfGen.PrintLnPdfText(d.data.DeliveryNodeTexts.ClosingText, "", "L")
}

func (d *DeliveryNode) printSignatureSection() {
	const signatureHeight = 50
	d.pdfGen.NewLine(din5008a.BodyStartX)
	d.pdfGen.NewLine(din5008a.BodyStartX)
	d.pdfGen.NewLine(din5008a.BodyStartX)
	_, y := d.pdfGen.GetCursor()
	var startSignatureSectionOnPosY = y //230
	var startSupplierX = din5008a.BodyStartX
	var startCustomerX = ((din5008a.BodyStopX - din5008a.BodyStartX) / 2) + din5008a.BodyStartX
	var senderSignatureName string

	if d.data.SenderAddress.CompanyName != "" {
		senderSignatureName = d.data.SenderAddress.CompanyName
	} else {
		senderSignatureName = "Lieferant"
	}

	d.printSignaturePart(senderSignatureName, startSupplierX, startSignatureSectionOnPosY)
	d.printSignaturePart("Kunde", startCustomerX, startSignatureSectionOnPosY)

}

func (d *DeliveryNode) printSignaturePart(headText string, startX float64, startY float64) {
	const contentWidth = (din5008a.BodyStopX - din5008a.BodyStartX) / 2
	const marginLeft = 22.5
	const nameLength = contentWidth - marginLeft
	const dateLength = nameLength / 3
	const gabLength = 5
	const signatureLength = (nameLength/3)*2 - gabLength

	var cY float64

	// name
	d.pdfGen.SetCursor(startX, startY)
	d.pdfGen.DrawLine(startX, startY, startX+nameLength, startY)
	_, cY = d.pdfGen.GetCursor()
	d.pdfGen.SetCursor(startX, cY+1)
	d.pdfGen.SetFontSize(d.meta.Font.SizeSmall)
	d.pdfGen.PrintPdfText(headText, "b", "L")
	d.pdfGen.PrintLnPdfText("(Name)", "", "L")

	d.pdfGen.SetFontSize(d.meta.Font.SizeDefault)
	d.pdfGen.NewLine(startX)
	d.pdfGen.NewLine(startX)
	d.pdfGen.NewLine(startX)

	//date & signature
	_, cY = d.pdfGen.GetCursor()
	var dateEndX = startX + dateLength
	var signatureStartX = startX + dateLength + gabLength
	var signatureEndX = startX + dateLength + gabLength + signatureLength
	d.pdfGen.DrawLine(startX, cY, dateEndX, cY)
	d.pdfGen.DrawLine(signatureStartX, cY, signatureEndX, cY)
	_, cY = d.pdfGen.GetCursor()
	d.pdfGen.SetCursor(startX, cY+1)
	d.pdfGen.SetFontSize(d.meta.Font.SizeSmall)
	d.pdfGen.PrintPdfText("Datum", "", "L")
	_, cY = d.pdfGen.GetCursor()
	d.pdfGen.SetCursor(signatureStartX, cY)
	d.pdfGen.PrintPdfText("Unterschrift", "", "L")
	d.pdfGen.SetFontSize(din5008a.FontGab10)
}

func (d *DeliveryNode) printFooter() {
	footerStartY, err := din5008a.Footer(d.printFooterContent, d.pdfGen)

	if err != nil {
		d.pdfGen.SetError(err)
	}

	if d.footerStartY == 0 {
		d.footerStartY = footerStartY
	}
}

func (d *DeliveryNode) printFooterContent(maxFooterHeight float64) (footerStartY float64) {
	// calculate height
	var currentStartX float64
	var currentY float64
	d.pdfGen.SetUnsafeCursor(din5008a.BodyStartX, maxFooterHeight)
	d.pdfGen.PreviousLine(din5008a.BodyStartX)
	d.pdfGen.PreviousLine(din5008a.BodyStartX)
	d.pdfGen.PreviousLine(din5008a.BodyStartX)
	_, currentY = d.pdfGen.GetCursor()
	footerStartY = currentY

	currentStartX = din5008a.BodyStartX
	d.pdfGen.SetCursor(currentStartX, footerStartY)
	d.pdfGen.NewLine(currentStartX)
	d.pdfGen.PrintPdfText(d.data.SenderInfo.Web, "", "L")

	currentStartX = ((din5008a.BodyStopX - din5008a.BodyStartX) / 2) + din5008a.BodyStartX
	d.pdfGen.SetCursor(currentStartX, footerStartY)
	d.pdfGen.NewLine(currentStartX)
	d.pdfGen.PrintPdfText(d.data.SenderInfo.Phone, "", "C")

	currentStartX = din5008a.BodyStopX
	d.pdfGen.SetCursor(currentStartX, footerStartY)
	d.pdfGen.NewLine(currentStartX)
	d.pdfGen.PrintPdfText(d.data.SenderInfo.Email, "", "R")

	return footerStartY
}

func (d *DeliveryNode) printHeader() {
	if d.data.SenderInfo.MimeLogoUrl != "" {
		din5008a.MimeImageHeader(d.pdfGen, d.data.SenderInfo.MimeLogoUrl)
	}
}
