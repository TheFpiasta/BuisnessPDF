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

type TableAttachment struct {
	data          tableAttachmentRequestData
	logger        *zerolog.Logger
	printErrStack bool
	pdfGen        *generator.PDFGenerator
}

type tableAttachmentRequestData struct {
	TableHeader       []string   `json:"tableHeader"`
	TableData         [][]string `json:"tableData"`
	ColumnPercentages []float64  `json:"columnPercentages"`
}

func NewTableAttachment(logger *zerolog.Logger) *TableAttachment {
	return &TableAttachment{
		data:          tableAttachmentRequestData{},
		logger:        logger,
		printErrStack: logger.GetLevel() <= zerolog.DebugLevel,
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
		//TODO implement print info text
		//TODO implement print header
		//TODO implement print table
	})

	din5008a.PageNumbering(t.pdfGen, din5008a.Height-30)
}
