package pdfType

import (
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
}

type tableAttachmentRequestData struct {
	TableHeader       []string   `json:"tableHeader"`
	TableData         [][]string `json:"tableData"`
	ColumnPercentages []float64  `json:"columnPercentages"`
}

func NewTableAttachment(logger *zerolog.Logger) *TableAttachment {
	t := &TableAttachment{
		data:          tableAttachmentRequestData{},
		logger:        logger,
		printErrStack: logger.GetLevel() <= zerolog.DebugLevel,
	}

	return t
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
