package pdfType

import (
	"SimpleInvoice/generator"
	din5008a "SimpleInvoice/norms/letter/din-5008-a"
	"encoding/json"
	"errors"
	errorsWithStack "github.com/go-errors/errors"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog"
	"io"
	"net/http"
)

type TableAttachment struct {
	data          tableAttachmentRequestData
	logger        *zerolog.Logger
	printErrStack bool
	pdfGen        *generator.PDFGenerator
	footerStartY  float64
}

type tableAttachmentRequestData struct {
	Headline          string     `json:"headline"`
	TableInfo         string     `json:"tableInfo"`
	TableHeader       []string   `json:"tableHeader"`
	TableData         [][]string `json:"tableData"`
	ColumnPercentages []float64  `json:"columnPercentages"`
	PageNumberPrefix  string     `json:"pageNumberPrefix"`
}

func NewTableAttachment(logger *zerolog.Logger) *TableAttachment {
	return &TableAttachment{
		data:          tableAttachmentRequestData{},
		logger:        logger,
		printErrStack: logger.GetLevel() <= zerolog.DebugLevel,
		pdfGen:        nil,
		footerStartY:  din5008a.Height - 5,
	}
}

func (t *TableAttachment) SetDataFromRequest(request *http.Request) (err error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.LogError(err)
		}
	}(request.Body)

	err = json.NewDecoder(request.Body).Decode(&t.data)
	if err != nil {
		return err
	}

	err = t.validateData()
	if err != nil {
		t.data = tableAttachmentRequestData{}
		return err
	}

	return nil
}

func (t *TableAttachment) GeneratePDF() (*gofpdf.Fpdf, error) {

	t.logger.Debug().Msg("generate table attachment")

	pdfGen, err := generator.NewPDFGenerator(
		generator.MetaData{
			FontName:         "OpenSans",
			FontGapY:         1.3,
			FontSize:         din5008a.FontSize10,
			MarginLeft:       din5008a.BodyStartX,
			MarginTop:        din5008a.AddressSenderTextStartY,
			MarginRight:      din5008a.Width - din5008a.BodyStopX,
			MarginBottom:     0,
			Unit:             "mm",
			DefaultLineWidth: 0.4,
			DefaultLineColor: generator.Color{R: 162, G: 162, B: 162},
		},
		false,
		t.logger,
		func() {

		},
		func(isLastPage bool) {

		},
	)

	if err != nil {
		return nil, err
	}

	t.pdfGen = pdfGen
	t.pdfGen.NewPage()

	t.doGenerate()

	return t.pdfGen.GetPdf(), t.pdfGen.GetError()
}

func (t *TableAttachment) LogError(err error) {
	var errStr string

	if _, ok := err.(*errorsWithStack.Error); ok && t.printErrStack {
		errStr = err.(*errorsWithStack.Error).ErrorStack()
	} else {
		errStr = err.Error()
	}

	t.logger.Error().Msgf(errStr)
}

func (t *TableAttachment) validateData() (err error) {
	//TODO implement me
	return err
}

func (t *TableAttachment) doGenerate() {

	din5008a.Body(t.pdfGen, func() {
		t.printHeadline()
		t.printTimeInfo()
		t.printTable()
	})

	din5008a.PageNumberingCustom(t.data.PageNumberPrefix, t.pdfGen, t.footerStartY, false)
}

func (t *TableAttachment) printHeadline() {
	t.pdfGen.SetFontSize(din5008a.FontSize10 + 5)
	x, y := t.pdfGen.GetCursor()

	//todo is this DIN conform or how to design the second page???
	y = din5008a.HeaderStopY + 5
	t.pdfGen.SetCursor(x, y)
	t.pdfGen.PrintLnPdfText(t.data.Headline, "b", "L")
	t.pdfGen.SetFontSize(din5008a.FontSize10)
	t.pdfGen.NewLine(din5008a.BodyStartX)
}

func (t *TableAttachment) printTimeInfo() {
	t.pdfGen.NewLine(din5008a.BodyStartX)
	t.pdfGen.SetFontSize(din5008a.FontSizeSender8)
	t.pdfGen.PrintLnPdfText(t.data.TableInfo, "i", "L")
	t.pdfGen.SetFontSize(din5008a.FontSize10)
}

func (t *TableAttachment) printTable() {
	// check, if ColumnPercentages is nearly 100%
	// small inaccuracy is allowed
	// todo delete this and add it to request data validation
	// todo rewrite to correct 100% validation
	// ---
	var pFull float64
	for _, percentage := range t.data.ColumnPercentages {
		pFull += percentage
	}
	if pFull < 99.9 || pFull > 100.1 {
		t.LogError(errors.New("sum of ColumnPercentages out of range"))
		return
	}
	// ---

	var columnWidth = getColumnWithFromPercentage(t.pdfGen, t.data.ColumnPercentages)
	var cellAlign []string

	for i := 0; i <= len(t.data.TableData); i++ {
		cellAlign = append(cellAlign, "LM")
	}

	t.pdfGen.PrintTableHeader(t.data.TableHeader, columnWidth, cellAlign)
	t.pdfGen.PrintTableBody(t.data.TableData, columnWidth, cellAlign)
}
