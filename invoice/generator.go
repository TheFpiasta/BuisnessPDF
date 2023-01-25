package invoice

import (
	"github.com/jung-kurt/gofpdf"
	"net/http"
)

func (iv *Invoice) newPDF() {
	iv.pdf = gofpdf.New("P", "mm", "A4", "")
	//iv.pdf.AddUTF8Font("dejavu", "", example.FontFile("DejaVuSansCondensed.ttf"))
	iv.pdf.AddPage()
}

func (iv *Invoice) drawPdfTextCell(posX float64, posY float64, text string, styleStr string, textSize float64, elementWith float64, elementHeight float64) {
	pdf := iv.pdf
	pdf.SetXY(posX, posY)
	pdf.SetFont(iv.textFont, styleStr, textSize)
	pdf.Cell(elementWith, elementHeight, text)
}

func (iv *Invoice) drawPdfTextRightAligned(posX float64, posY float64, text string, styleStr string, textSize float64, elementWith float64, elementHeight float64) {
	pdf := iv.pdf
	pdf.SetFont(iv.textFont, styleStr, textSize)
	stringWidth := pdf.GetStringWidth(text) + 2
	pdf.SetXY(posX-stringWidth, posY)
	pdf.Cell(elementWith, elementHeight, text)
}

func (iv *Invoice) drawLine(x1 float64, y1 float64, x2 float64, y2 float64, color color) {
	pdf := iv.pdf
	pdf.SetDrawColor(int(color.r), int(color.g), int(color.b))
	pdf.Line(x1, y1, x2, y2)
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
		pdf.Image(logoUrl, float64(posX), float64(posY), imgWd/2, imgHt/2, false, tp, 0, "")
	}

	return nil

}
