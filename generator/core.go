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
	//pageWidth, _ := core.pdf.GetPageSize()
	//saveWriteArea := pageWidth - core.data.MarginLeft - core.pdf.GetX()
	_, lineHeight := core.pdf.GetFontSize()
	stringWidth := core.pdf.GetStringWidth(text) + 2

	switch alignStr {
	case "L":
		core.pdf.Cell(stringWidth, lineHeight, text)
	case "R":
		x := core.pdf.GetX()

		core.pdf.SetX(x - stringWidth)
		core.pdf.Cell(stringWidth, lineHeight, text)
		//core.pdf.SetX(x - stringWidth)
	case "C":
		x := core.pdf.GetX()

		core.pdf.SetX(x - stringWidth/2)
		core.pdf.Cell(stringWidth, lineHeight, text)
		//core.pdf.SetX(x - stringWidth/2)
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
	newY := core.pdf.GetY() + lineHeight + core.data.FontGapY
	core.pdf.SetXY(oldX, newY)
}

// PrintPdfTextFormatted
//
//	text
//
//	styleStr
//
//	alignStr
//
//	borderStr: specifies how the cell border will be drawn. An empty string indicates no border, "1" indicates a full border, and one or more of "L", "T", "R" and "B" indicate the left, top, right and bottom sides of the border.
//
//	fill: is true to paint the cell background or false to leave it transparent.
//
//	backgroundColor
func (core *PDFGenerator) PrintPdfTextFormatted(text string, styleStr string, alignStr string, borderStr string, fill bool, backgroundColor Color, lineHeight float64, stringWidth float64) {
	core.pdf.SetFont(core.data.FontName, styleStr, core.GetFontSize())
	//stringWidth := core.pdf.GetStringWidth(text) + 2

	core.pdf.SetFillColor(int(backgroundColor.R), int(backgroundColor.G), int(backgroundColor.B))
	switch alignStr {
	case "L":
		core.pdf.CellFormat(stringWidth, lineHeight, text, borderStr, 0, "LM", fill, 0, "")
	case "R":
		x := core.pdf.GetX()
		core.pdf.SetX(x - stringWidth)
		core.pdf.CellFormat(stringWidth, lineHeight, text, borderStr, 0, "LM", fill, 0, "")
	case "C":
		x := core.pdf.GetX()
		core.pdf.SetX(x - stringWidth/2)
		core.pdf.CellFormat(stringWidth, lineHeight, text, borderStr, 0, "LM", fill, 0, "")
	default:
		core.pdf.SetError(errors.New("can't interpret the given text align code"))
	}
}

func (core *PDFGenerator) DrawPdfTextRightAligned(posXRight float64, posY float64, text string, styleStr string, textSize float64, elementWith float64, elementHeight float64) {
	core.pdf.SetFont(core.data.FontName, styleStr, textSize)
	stringWidth := core.pdf.GetStringWidth(text) + 2
	core.pdf.SetXY(posXRight-stringWidth, posY)
	core.pdf.WriteAligned(core.pdf.GetStringWidth(text), core.data.LineHeight, text, "R")
	core.pdf.Cell(elementWith, elementHeight, text)
}

func (core *PDFGenerator) DrawLine(x1 float64, y1 float64, x2 float64, y2 float64, color Color, lineWith float64) {
	core.pdf.SetLineWidth(lineWith)
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

func (core *PDFGenerator) PrintTable(header []string, columnWidth []float64, items [][]string, summary [][2]string, summaryWidths [3]float64) {
	var tableWidth float64
	for _, width := range columnWidth {
		tableWidth += width
	}

	core.printTableHeader(header, columnWidth)
	core.printTableItems(items, columnWidth)
	core.printTableSummary(summary, summaryWidths, tableWidth)
}

func (core *PDFGenerator) printTableHeader(header []string, columnWidth []float64) {
	referenceX := core.pdf.GetX()
	_, lineHeight := core.pdf.GetFontSize()
	newlineHeight := lineHeight + core.data.FontGapY*2

	for i, text := range header {
		core.PrintPdfTextFormatted(text, "b", "L", "TB", true, Color{R: 239, G: 239, B: 239}, newlineHeight, columnWidth[i])
	}

	core.SetCursor(referenceX, core.pdf.GetY()+newlineHeight)
}

func (core *PDFGenerator) printTableItems(items [][]string, columnWidth []float64) {
	referenceX := core.pdf.GetX()
	_, lineHeight := core.pdf.GetFontSize()
	newlineHeight := lineHeight + core.data.FontGapY*2

	for _, item := range items {
		for j, text := range item {
			core.PrintPdfTextFormatted(text, "", "L", "B", false, Color{R: 239, G: 239, B: 239}, newlineHeight, columnWidth[j])
		}
		core.SetCursor(referenceX, core.pdf.GetY()+newlineHeight)
	}
}

func (core *PDFGenerator) printTableSummary(summary [][2]string, summaryWidths [3]float64, tableWidth float64) {
	referenceX := core.pdf.GetX()
	_, lineHeight := core.pdf.GetFontSize()
	newlineHeight := lineHeight + core.data.FontGapY*2

	for i, text := range summary {
		boarderStr := ""
		if len(summary)-1 == i {
			boarderStr = "B"
		}

		core.PrintPdfTextFormatted("", "", "L", "", false, Color{R: 239, G: 239, B: 239}, newlineHeight, summaryWidths[0])
		core.PrintPdfTextFormatted(text[0], "", "L", boarderStr, false, Color{R: 239, G: 239, B: 239}, newlineHeight, summaryWidths[1])
		core.PrintPdfTextFormatted(text[1], "", "L", boarderStr, false, Color{R: 239, G: 239, B: 239}, newlineHeight, summaryWidths[2])
		core.SetCursor(referenceX, core.pdf.GetY()+newlineHeight)
	}

}
