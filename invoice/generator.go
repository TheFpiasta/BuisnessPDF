package invoice

import (
	"github.com/jung-kurt/gofpdf"
	"net/http"
)

func (iv *Invoice) newPDF() {
	iv.pdf = gofpdf.New("P", "mm", "A4", "")

	iv.pdf.AddUTF8Font(iv.textFont, "", "fonts/OpenSans-Regular.ttf")
	iv.pdf.AddUTF8Font(iv.textFont, "l", "fonts/OpenSans-Light.ttf")
	iv.pdf.AddUTF8Font(iv.textFont, "i", "fonts/OpenSans-Italic.ttf")
	iv.pdf.AddUTF8Font(iv.textFont, "b", "fonts/OpenSans-Bold.ttf")
	iv.pdf.AddUTF8Font(iv.textFont, "m", "fonts/OpenSans-Medium.ttf")

	iv.pdf.SetFont(iv.textFont, "", 12)
	iv.pdf.SetMargins(25, 45, 20)
	iv.pdf.SetHomeXY()
	//iv.pdf.AliasNbPages("{entute}")

	iv.pdf.AddPage()
}

// writePdfText
//
//	text		the text to write
//	styleStr	"" default, "l" light, "i" italic, "b" bold, "m" medium
//	textSize	the text size
//	alignStr	"L" right, "C" center, "R" right
func (iv *Invoice) writePdfText(text string, styleStr string, textSize float64, alignStr string) {
	//TODO refactor to cell!
	iv.pdf.SetFont(iv.textFont, styleStr, textSize)
	textWidthGlyf := iv.pdf.GetStringSymbolWidth(text)

	textWidthGlyf = textWidthGlyf / 64

	iv.pdf.WriteAligned(float64(textWidthGlyf), iv.lineHeight, text, alignStr)
	//iv.pdf.Cell(iv.pdf.GetStringWidth(text), elementHeight, text)
	iv.pdf.Ln(iv.lineHeight)
}

func (iv *Invoice) drawPdfTextRightAligned(posXRight float64, posY float64, text string, styleStr string, textSize float64, elementWith float64, elementHeight float64) {
	iv.pdf.SetFont(iv.textFont, styleStr, textSize)
	stringWidth := iv.pdf.GetStringWidth(text) + 2
	iv.pdf.SetXY(posXRight-stringWidth, posY)
	iv.pdf.WriteAligned(iv.pdf.GetStringWidth(text), iv.lineHeight, text, "R")
	iv.pdf.Cell(elementWith, elementHeight, text)
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
