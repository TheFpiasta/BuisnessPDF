package invoice

import (
	"github.com/jung-kurt/gofpdf"
	"net/http"
)

func (iv *Invoice) newPDF() {
	iv.pdf = gofpdf.New("P", "mm", "A4", "")
	iv.pdf.AddPage()
}

func (iv *Invoice) setPdfText(posX float64, posY float64, text string, styleStr string, textSize float64, elementWith float64, elementHeight float64) {
	pdf := iv.pdf

	pdf.SetXY(posX, posY)
	pdf.SetFont(iv.textFont, styleStr, textSize)
	pdf.Cell(elementWith, elementHeight, text)
}

func (iv *Invoice) placeImgOnPosXY(logoUrl string, posX int, posY int) (err error) {
	var (
		rsp *http.Response
		tp  string
	)
	pdf := iv.pdf

	rsp, err = http.Get(logoUrl)
	if err != nil {
		pdf.SetError(err)
		return iv.handleError(err, "Can not get logo!")
	}

	tp = pdf.ImageTypeFromMime(rsp.Header["Content-Type"][0])
	infoPtr := pdf.RegisterImageReader(logoUrl, tp, rsp.Body)
	if pdf.Ok() {
		imgWd, imgHt := infoPtr.Extent()
		pdf.Image(logoUrl, float64(posX), float64(posY), imgWd, imgHt, false, tp, 0, "")
	}

	return nil

}
