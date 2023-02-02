package invoice

import (
	"errors"
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

	iv.pdf.SetFont(iv.textFont, "", iv.textSize)
	iv.pdf.SetMargins(iv.marginLeft, iv.marginTop, iv.marginRight)
	iv.pdf.SetHomeXY()
	//iv.pdf.AliasNbPages("{entute}")

	iv.pdf.AddPage()
}

// printPdfText
//
//	text		the text to write
//	styleStr	"" default, "l" light, "i" italic, "b" bold, "m" medium
//	textSize	the text size
//	alignStr	"L" right, "C" center, "R" right
func (iv *Invoice) printPdfText(text string, styleStr string, textSize float64, alignStr string) {
	iv.pdf.SetFont(iv.textFont, styleStr, textSize)
	pageWidth, _ := iv.pdf.GetPageSize()
	saveWriteArea := pageWidth - iv.marginLeft - iv.pdf.GetX()
	_, lineHeight := iv.pdf.GetFontSize()

	switch alignStr {
	case "L":
		iv.pdf.Cell(saveWriteArea/2, lineHeight, text)
	case "R":
		stringWidth := iv.pdf.GetStringWidth(text) + 2
		x := iv.pdf.GetX()

		iv.pdf.SetX(x - stringWidth)
		iv.pdf.Cell(stringWidth, lineHeight, text)
		iv.pdf.SetX(x - stringWidth)
	case "C":
	default:
		iv.pdf.SetError(errors.New("can't interpret the given text align code"))
	}
}

// printLnPdfText
//
//	 prints a line with line break
//
//		text		the text to print
//		styleStr	"" default, "l" light, "i" italic, "b" bold, "m" medium
//		textSize	the text size
//		alignStr	"L" right, "C" center, "R" right
func (iv *Invoice) printLnPdfText(text string, styleStr string, textSize float64, alignStr string) {
	currentX := iv.pdf.GetX()
	iv.printPdfText(text, styleStr, textSize, alignStr)
	iv.newLine(currentX)
}

func (iv *Invoice) newLine(oldX float64) {
	_, lineHeight := iv.pdf.GetFontSize()
	currentY := iv.pdf.GetY() + lineHeight + iv.fontGapY
	iv.pdf.SetXY(oldX, currentY)
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

	return pdf.Error()

}
