package generator

import (
	"errors"
	"github.com/jung-kurt/gofpdf"
	"math"
	"net/http"
	"strings"
)

// NewPDFGenerator create a new PDFGenerator instance.
func NewPDFGenerator(data MetaData) (gen *PDFGenerator) {
	gen = new(PDFGenerator)

	pdf := gofpdf.New("P", data.Unit, "A4", "")

	pdf.AddUTF8Font("OpenSans", "", "fonts/OpenSans-Regular.ttf")
	pdf.AddUTF8Font("OpenSans", "l", "fonts/OpenSans-Light.ttf")
	pdf.AddUTF8Font("OpenSans", "i", "fonts/OpenSans-Italic.ttf")
	pdf.AddUTF8Font("OpenSans", "b", "fonts/OpenSans-Bold.ttf")
	pdf.AddUTF8Font("OpenSans", "m", "fonts/OpenSans-Medium.ttf")
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
	_, lineHeight := core.pdf.GetFontSize()
	stringWidth := core.pdf.GetStringWidth(text) + 2

	switch alignStr {
	case "L":
		core.pdf.Cell(stringWidth, lineHeight, text)
	case "R":
		x := core.pdf.GetX()

		core.pdf.SetX(x - stringWidth)
		core.pdf.Cell(stringWidth, lineHeight, text)
	case "C":
		x := core.pdf.GetX()

		core.pdf.SetX(x - stringWidth/2)
		core.pdf.Cell(stringWidth, lineHeight, text)
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
	lines := core.extractLinesFromText(text)
	currentX := core.pdf.GetX()

	for _, line := range lines {
		core.PrintPdfText(line, styleStr, alignStr)
		core.NewLine(currentX)
	}
}

func (core *PDFGenerator) NewLine(oldX float64) {
	_, lineHeight := core.pdf.GetFontSize()
	newY := core.pdf.GetY() + lineHeight + core.data.FontGapY
	core.pdf.SetXY(oldX, newY)
}

func (core *PDFGenerator) extractLinesFromText(text string) (textLines []string) {
	textLines = strings.Split(text, "\n")

	for i, line := range textLines {
		removeStr := 0

		for _, c := range line {
			if c != 32 {
				break
			}
			removeStr++
		}

		if removeStr > 0 {
			textLines[i] = line[removeStr:]
		}
	}

	return textLines
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
	core.pdf.CellFormat(stringWidth, lineHeight, text, borderStr, 0, alignStr, fill, 0, "")
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

func (core *PDFGenerator) PrintTableHeader(cells []string, columnWidth []float64) {
	referenceX := core.pdf.GetX()
	_, lineHeight := core.pdf.GetFontSize()
	newlineHeight := lineHeight + core.data.FontGapY*2

	for i, cell := range cells {
		core.PrintPdfTextFormatted(cell, "b", "LM", "TB", true, Color{R: 239, G: 239, B: 239}, newlineHeight, columnWidth[i])
	}

	core.SetCursor(referenceX, core.pdf.GetY()+newlineHeight)
}

func (core *PDFGenerator) PrintTableBody(cells [][]string, columnWidths []float64, columnAlignStrings []string) {
	referenceX := core.pdf.GetX()
	_, lineHeight := core.pdf.GetFontSize()
	newlineHeight := lineHeight + core.data.FontGapY*2

	for _, row := range cells {
		var extractedLines [][]string
		var maxLines = 0

		for _, cell := range row {
			extractedItem := core.extractLinesFromText(cell)
			maxLines = int(math.Max(float64(maxLines), float64(len(extractedItem))))
			extractedLines = append(extractedLines, extractedItem)
		}

		for i := 0; i < maxLines; i++ {
			core.printTableRow(extractedLines, i, maxLines, columnAlignStrings, newlineHeight, columnWidths, referenceX)
		}
	}
}

func (core *PDFGenerator) printTableRow(extractedLines [][]string, currentLine int, maxItems int, alignStrings []string, newlineHeight float64, columnWidth []float64, referenceX float64) {
	for j, cell := range extractedLines {
		var text = ""
		var borderStr = ""

		if currentLine < len(cell) {
			text = cell[currentLine]
		}

		if currentLine == maxItems-1 {
			borderStr = "B"
		}

		core.PrintPdfTextFormatted(text, "", alignStrings[j], borderStr, false, Color{R: 239, G: 239, B: 239}, newlineHeight, columnWidth[j])
	}
	core.SetCursor(referenceX, core.pdf.GetY()+newlineHeight)
}

func (core *PDFGenerator) PrintTableFooter(cells [][]string, columnWidths []float64, columnAlignStrings []string) {
	referenceX := core.pdf.GetX()
	_, lineHeight := core.pdf.GetFontSize()
	newlineHeight := lineHeight + core.data.FontGapY*2

	for i, row := range cells {
		boarderStr := ""
		fill := false
		styleStr := ""

		if len(cells)-1 == i {
			boarderStr = "BT"
			fill = true
			styleStr = "B"
		}

		for j, cell := range row {
			if cell == "" {
				core.PrintPdfTextFormatted(row[j], "", columnAlignStrings[j], "", false, Color{R: 239, G: 239, B: 239}, newlineHeight, columnWidths[j])
			} else {
				core.PrintPdfTextFormatted(row[j], styleStr, columnAlignStrings[j], boarderStr, fill, Color{R: 239, G: 239, B: 239}, newlineHeight, columnWidths[j])
			}
		}

		core.SetCursor(referenceX, core.pdf.GetY()+newlineHeight)
	}
}
