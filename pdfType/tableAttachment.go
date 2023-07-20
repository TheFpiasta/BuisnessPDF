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
	TimeInfo          string     `json:"timeInfo"`
	TableHeader       []string   `json:"tableHeader"`
	TableData         [][]string `json:"tableData"`
	ColumnPercentages []float64  `json:"columnPercentages"`
}

func NewTableAttachment(logger *zerolog.Logger) *TableAttachment {
	return &TableAttachment{
		data:          tableAttachmentRequestData{},
		logger:        logger,
		printErrStack: logger.GetLevel() <= zerolog.DebugLevel,
		pdfGen:        nil,
		footerStartY:  din5008a.Height - 30,
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
			FontName:         "",
			FontGapY:         0,
			FontSize:         0,
			MarginLeft:       0,
			MarginTop:        0,
			MarginRight:      0,
			MarginBottom:     0,
			Unit:             "",
			DefaultLineWidth: 0,
			DefaultLineColor: generator.Color{},
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

	//TODO implement me
	panic("implement me")
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

	din5008a.PageNumbering(t.pdfGen, t.footerStartY)
}

func (t *TableAttachment) printHeadline() {
	t.pdfGen.SetFontSize(din5008a.FontSize10 + 5)
	x, y := t.pdfGen.GetCursor()

	//todo is this DIN conform or how to design the second page???
	y = din5008a.HeaderStopY + 30
	t.pdfGen.SetCursor(x, y)
	t.pdfGen.PrintLnPdfText("Anhang", "b", "L")
	t.pdfGen.SetFontSize(din5008a.FontSize10)
	t.pdfGen.NewLine(din5008a.BodyStartX)
}

func (t *TableAttachment) printTimeInfo() {
	t.pdfGen.NewLine(din5008a.BodyStartX)
	t.pdfGen.SetFontSize(din5008a.FontSizeSender8)
	t.pdfGen.PrintLnPdfText(t.data.TimeInfo, "i", "L")
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

	for i := 0; i < len(t.data.TableData); i++ {
		cellAlign = append(cellAlign, "LM")
	}

	t.pdfGen.PrintTableHeader(t.data.TableHeader, columnWidth, cellAlign)
	t.pdfGen.PrintTableBody(t.data.TableData, columnWidth, cellAlign)
}
