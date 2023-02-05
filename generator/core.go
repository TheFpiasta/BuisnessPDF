package generator

import (
	"errors"
	"github.com/jung-kurt/gofpdf"
	"net/http"
)

func NewPDFGenerator(data MetaData) (gen *PDFGenerator) {
	gen = new(PDFGenerator)

	pdf := gofpdf.New("P", data.Unit, "A4", "")

	pdf.AddUTF8Font(data.FontName, "", "fonts/OpenSans-Regular.ttf")
	pdf.AddUTF8Font(data.FontName, "l", "fonts/OpenSans-Light.ttf")
	pdf.AddUTF8Font(data.FontName, "i", "fonts/OpenSans-Italic.ttf")
	pdf.AddUTF8Font(data.FontName, "b", "fonts/OpenSans-Bold.ttf")
	pdf.AddUTF8Font(data.FontName, "m", "fonts/OpenSans-Medium.ttf")
	pdf.SetFont(data.FontName, "", data.FontSize)
	pdf.SetMargins(data.MarginLeft, data.MarginTop, data.MarginRight)
	pdf.SetHomeXY()
	//iv.pdf.AliasNbPages("{entute}")
	pdf.AddPage()

	gen.pdf = pdf
	gen.data = data

	return
}

func (core *PDFGenerator) SetCursor(x float64, y float64) {
	core.pdf.SetXY(x, y)
}

// PrintPdfText
//
//	text		the text to write
//	styleStr	"" default, "l" light, "i" italic, "b" bold, "m" medium
//	textSize	the text size
//	alignStr	"L" right, "C" center, "R" right
func (core *PDFGenerator) PrintPdfText(text string, styleStr string, alignStr string) {
	core.pdf.SetFont(core.data.FontName, styleStr, core.GetFontSize())
	pageWidth, _ := core.pdf.GetPageSize()
	saveWriteArea := pageWidth - core.data.MarginLeft - core.pdf.GetX()
	_, lineHeight := core.pdf.GetFontSize()

	switch alignStr {
	case "L":
		core.pdf.Cell(saveWriteArea/2, lineHeight, text)
	case "R":
		stringWidth := core.pdf.GetStringWidth(text) + 2
		x := core.pdf.GetX()

		core.pdf.SetX(x - stringWidth)
		core.pdf.Cell(stringWidth, lineHeight, text)
		core.pdf.SetX(x - stringWidth)
	case "C":
	default:
		core.pdf.SetError(errors.New("can't interpret the given text align code"))
	}
}

// PrintLnPdfText
//
//	 prints a line with line break
//
//		text		the text to print
//		styleStr	"" default, "l" light, "i" italic, "b" bold, "m" medium
//		textSize	the text size
//		alignStr	"L" right, "C" center, "R" right
func (core *PDFGenerator) PrintLnPdfText(text string, styleStr string, alignStr string) {
	currentX := core.pdf.GetX()
	core.PrintPdfText(text, styleStr, alignStr)
	core.newLine(currentX)
}

func (core *PDFGenerator) newLine(oldX float64) {
	_, lineHeight := core.pdf.GetFontSize()
	currentY := core.pdf.GetY() + lineHeight + core.data.FontGapY
	core.pdf.SetXY(oldX, currentY)
}

func (core *PDFGenerator) DrawPdfTextRightAligned(posXRight float64, posY float64, text string, styleStr string, textSize float64, elementWith float64, elementHeight float64) {
	core.pdf.SetFont(core.data.FontName, styleStr, textSize)
	stringWidth := core.pdf.GetStringWidth(text) + 2
	core.pdf.SetXY(posXRight-stringWidth, posY)
	core.pdf.WriteAligned(core.pdf.GetStringWidth(text), core.data.LineHeight, text, "R")
	core.pdf.Cell(elementWith, elementHeight, text)
}

func (core *PDFGenerator) DrawLine(x1 float64, y1 float64, x2 float64, y2 float64, color Color) {
	core.pdf.SetDrawColor(int(color.R), int(color.G), int(color.B))
	core.pdf.Line(x1, y1, x2, y2)
}

func (core *PDFGenerator) PlaceImgOnPosXY(logoUrl string, posX int, posY int) (err error) {
	var (
		rsp *http.Response
		tp  string
	)

	rsp, err = http.Get(logoUrl)
	if err != nil {
		core.pdf.SetError(err)
		return core.pdf.Error()
	}

	tp = core.pdf.ImageTypeFromMime(rsp.Header["Content-Type"][0])
	infoPtr := core.pdf.RegisterImageReader(logoUrl, tp, rsp.Body)
	if core.pdf.Ok() {
		imgWd, imgHt := infoPtr.Extent()
		core.pdf.Image(logoUrl, float64(posX), float64(posY), imgWd/2, imgHt/2, false, tp, 0, "")
	}

	return core.pdf.Error()
}

// TODO add to def
func (core *PDFGenerator) PrintTable(header []string, columnWidth []float64, items [][]string) {
	//fillColor := Color{R: 200, G: 200, B: 200}
	//const colNumber = 5
	//header := [colNumber]string{"No", "Description", "Quantity", "Unit Price ($)", "Price ($)"}
	//colWidth := [colNumber]float64{10.0, 50.0, 40.0, 30.0, 30.0}
	//lineHt := 10.0
	//pdfGen.SetCursor(iv.marginLeft, iv.pdfGen.GetY()+iv.lineHeight+10.0)
	//for colJ := 0; colJ < colNumber; colJ++ {
	//	iv.pdfGen.CellFormat(colWidth[colJ], lineHt, header[colJ], "1", 0, "CM", true, 0, "")
	//}

}

func (core *PDFGenerator) printTableHeader() {

}
