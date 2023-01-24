package invoice

import (
	"encoding/json"
	"io"
)

func (iv *Invoice) parseJsonData(jsonData io.ReadCloser) (err error) {
	return json.NewDecoder(jsonData).Decode(&iv.pdfData)
}

func (iv *Invoice) validateJsonData() (err error) {
	return nil
}
