package generator

import (
	"errors"
	"github.com/jung-kurt/gofpdf"
	"math"
	"net/http"
	"net/url"
	"strings"
)

// NewPDFGenerator create and return a new PDFGenerator instance.
// MetaData is used for all necessary inputs.
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

// SetCursor set the abscissa (x) and ordinate (y) reference point
// in the unit of measure specified in NewPDFGenerator() for the next operation.
// If the passed values are negative, they are relative respectively to the right and bottom of the page.
func (core *PDFGenerator) SetCursor(x float64, y float64) {
	core.pdf.SetXY(x, y)
}

// PrintPdfText prints from the current cursor position a simple text cell in the PDF.
//
// text passed the string to print.
//
// styleStr defines the font style:
//
//	"" non-specific font style
//	"l" light font
//	"i" italic font
//	"b" bold font
//	"m" medium font
//
// alignStr set the align mode:
//
//	"L" align the left side of the text to the current cursor position
//	"R" align the right side of the text to the current cursor position
//	"C" align the center of the text to the current cursor position
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

// PrintLnPdfText prints from the current cursor position a simple text cell in the PDF
// and call NewLine() at the end.
//
// text passed the string to print.
// Use \n escape character to trigger NewLine() inside the text.
//
// styleStr defines the font style:
//
//	"" non-specific font style
//	"l" light font
//	"i" italic font
//	"b" bold font
//	"m" medium font
//
// alignStr set the align mode:
//
//	"L" align the left side of the text to the current cursor position
//	"R" align the right side of the text to the current cursor position
//	"C" align the center of the text to the current cursor position
func (core *PDFGenerator) PrintLnPdfText(text string, styleStr string, alignStr string) {
	lines := core.extractLinesFromText(text)
	currentX := core.pdf.GetX()

	for _, line := range lines {
		core.PrintPdfText(line, styleStr, alignStr)
		core.NewLine(currentX)
	}
}

// NewLine sets the cursor on the next line dependent on the given X-position
// (mostly use the start X-point of the current line)
func (core *PDFGenerator) NewLine(oldX float64) {
	_, lineHeight := core.pdf.GetFontSize()
	newY := core.pdf.GetY() + lineHeight + core.data.FontGapY
	core.pdf.SetXY(oldX, newY)
}

// extractLinesFromText split a string on newline character (\n) and return the parts as an array.
// Prefixing whitespaces (ONLY " ")! will be automatically removed on each part.
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

// PrintPdfTextFormatted prints from the current cursor position a formatted text cell in the PDF
// (e.g. with boarders or background color).
//
// text passed the string to print.
//
// styleStr defines the font style:
//
//	 "" non-specific font style
//		"l" light font
//		"i" italic font
//		"b" bold font
//		"m" medium font
//
// alignStr set the align mode:
//
//	"L" align the left side of the text to the current cursor position
//	"R" align the right side of the text to the current cursor position
//	"C" align the center of the text to the current cursor position
//
// borderStr specifies how the cell border will be drawn:
//
//	An empty string indicates no border,
//	"1" indicates a full border,
//	one or more of "L", "T", "R" and "B" indicate the left, top, right and bottom sides of the border.
//
// fill defines, whether the background is set to the background color or not. If false, use a transparent background.
//
// backgroundColor defines the background color using the Color.
//
// cellHeight specifies the total height of the cell in the unit of measure specified in NewPDFGenerator().
//
// cellWidth specifies the total width of the cell in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) PrintPdfTextFormatted(text string, styleStr string, alignStr string, borderStr string, fill bool, backgroundColor Color, cellHeight float64, cellWidth float64) {
	core.pdf.SetFont(core.data.FontName, styleStr, core.GetFontSize())
	if fill {
		core.pdf.SetFillColor(int(backgroundColor.R), int(backgroundColor.G), int(backgroundColor.B))
	}
	core.pdf.CellFormat(cellWidth, cellHeight, text, borderStr, 0, alignStr, fill, 0, "")
}

//// DrawPdfTextRightAligned is a deprecated method to write text right aligned
////
//func (core *PDFGenerator) DrawPdfTextRightAligned(posXRight float64, posY float64, text string, styleStr string, textSize float64, elementWith float64, elementHeight float64) {
//	core.pdf.SetFont(core.data.FontName, styleStr, textSize)
//	stringWidth := core.pdf.GetStringWidth(text) + 2
//	core.pdf.SetXY(posXRight-stringWidth, posY)
//	core.pdf.WriteAligned(core.pdf.GetStringWidth(text), core.data.LineHeight, text, "R")
//	core.pdf.Cell(elementWith, elementHeight, text)
//}

// DrawLine draw a user defines line between two points.
//
// x1 and y1 defines the abscissa (x) and ordinate (y) cursor start point.
//
// x2 and y2 defines the abscissa (x) and ordinate (y) cursor end point.
//
// color specifies the color of the line.
//
// lineWith specifies the thinness of the line in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) DrawLine(x1 float64, y1 float64, x2 float64, y2 float64, color Color, lineWith float64) {
	core.pdf.SetLineWidth(lineWith)
	core.pdf.SetDrawColor(int(color.R), int(color.G), int(color.B))
	core.pdf.Line(x1, y1, x2, y2)
}

// PlaceMimeImageFromUrl downloade a JPEG, PNG or GIF image (from mostly a Content Delivery Network (CDN)) URL and puts it in the current page.
//
// cdnUrl specifies a parsed (CDN) URL.
//
// posX and posY define the top left cursor abscissa (x) and ordinate (y) position in the unit of measure specified in NewPDFGenerator(), where the image should be pleased.
//
// scale specifies the scaling factor into which the image is drawn.
// The value must be grater then 0. Use scaling of 1 for no scaling.
// E.g. a value of 0.5 means draw the image in half the size of the original
// and a value of 3 means draw the image in the triple size of the original.
func (core *PDFGenerator) PlaceMimeImageFromUrl(cdnUrl *url.URL, posX float64, posY float64, scale float64) (err error) {
	//TODO scale of (0, ...] abfangen

	var rsp *http.Response

	rsp, err = http.Get(cdnUrl.String())
	if err != nil {
		core.pdf.SetError(err)
		return core.pdf.Error()
	}

	imageMimeType := core.pdf.ImageTypeFromMime(rsp.Header["Content-Type"][0])
	imageInfoType := core.pdf.RegisterImageReader(cdnUrl.String(), imageMimeType, rsp.Body)
	if core.pdf.Ok() {
		imgWd, imgHt := imageInfoType.Extent()
		core.pdf.Image(cdnUrl.String(), posX, posY, imgWd*scale, imgHt*scale, false, imageMimeType, 0, "")
	}

	return core.pdf.Error()
}

// PrintTableHeader print a generic and clean styled table header.
//
// cells contains the displayed column names of the table header.
//
// columnWidth defines the width of each column. NOTE: in general use here the same widths as in PrintTableBody()
func (core *PDFGenerator) PrintTableHeader(cells []string, columnWidth []float64) {
	referenceX := core.pdf.GetX()
	_, lineHeight := core.pdf.GetFontSize()
	newlineHeight := lineHeight + core.data.FontGapY*2

	for i, cell := range cells {
		core.PrintPdfTextFormatted(cell, "b", "LM", "TB", true, Color{R: 239, G: 239, B: 239}, newlineHeight, columnWidth[i])
	}

	core.SetCursor(referenceX, core.pdf.GetY()+newlineHeight)
}

// PrintTableBody prints a generic and clean styled table content rows.
//
// cells contains an array with includes all rows.
// Each row is an array by its self includes the information of each cell.
// In fact, cells can be described as [rowNumber][columnNumber]contentString.
//
// columnWidths defines the width of each column. NOTE: in general use here the same widths as in PrintTableHeader().
//
// columnAlignStrings specifies the align type of each column. Use:
//
//	"L" for align the left side of the text to the left side of the table cell,
//	"R" for align the right side of the text to the right side of the table cell, and
//	"C" for align the text to the center of the table cell.
//
// E.g. in an invoice table, typically use "L" for all strings and "R" for salary.
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
			core.printTableBodyRow(extractedLines, i, maxLines, columnAlignStrings, newlineHeight, columnWidths, referenceX)
		}
	}
}

// printTableBodyRow prints one row of a table body.
//
// extractedLines includes all cells with the braked text.
//
// currentLine is used for detecting, if this row is the last row in the body.
//
// maxItems is used deine the death of the row.
// If one column includes less the maxItems content strings,
// the rest will be filled with empty content items to print the bottom boarder correctly.
//
// alignStrings specifies the align type of each column. Use:
//
//	"L" for align the left side of the text to the left side of the table cell,
//	"R" for align the right side of the text to the right side of the table cell, and
//	"C" for align the text to the center of the table cell.
//
// E.g. in an invoice table, typically use "L" for all strings and "R" for salary.
//
// newlineHeight specifies the overall row Height, for the amount of the maximum items in one column.
//
// columnWidth defines the width of each column. NOTE: use here the same widths as in PrintTableBody().
//
// referenceX defines the left row position to set the cursor at the end to a new line.
func (core *PDFGenerator) printTableBodyRow(extractedLines [][]string, currentLine int, maxItems int, alignStrings []string, newlineHeight float64, columnWidth []float64, referenceX float64) {
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

// PrintTableFooter prints a generic and clean styled table footer.
// The last row of the footer will be print in the same style as the table header.
//
// cells contains an array with includes all rows.
// Each row is an array by its self includes the information of each cell.
// In fact, cells can be described as [rowNumber][columnNumber]contentString.
//
// columnWidths defines the width of each column.
//
// columnAlignStrings specifies the align type of each column. Use:
//
//	"L" for align the left side of the text to the left side of the table cell,
//	"R" for align the right side of the text to the right side of the table cell, and
//	"C" for align the text to the center of the table cell.
//
// E.g. in an invoice table, typically use "L" for all strings and "R" for salary.
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
