package generator

import (
	"fmt"
	errorsWithStack "github.com/go-errors/errors"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog"
	"math"
	"net/http"
	"net/url"
	"strings"
)

// NewPDFGenerator construct and return a new PDFGenerator instance.
//
// MetaData is used for all necessary inputs.
//
// Set strictErrorHandling to true, to provide execution of any method, if a pdf internal error is set.
// If strictErrorHandling is set to false, all methods are tried to execute executed, even if a pdf internal error is set.
// This may cause the PDF internal error to be overwritten by a new error.
// Use GetError() to get the current pdf internal error.
func NewPDFGenerator(data MetaData, strictErrorHandling bool, logger *zerolog.Logger, headerFunction func(), footerFunction func(isLastPage bool)) (gen *PDFGenerator, err error) {
	// --> validate inputs
	if data.FontGapY < 0 {
		return nil, errorsWithStack.New(fmt.Sprintf("A negative FontGapY (%f) is not allowed.", data.FontGapY))
	}
	if data.FontSize <= 0 {
		return nil, errorsWithStack.New(fmt.Sprintf("Text size must be grather or equal then 0."))
	}

	validUnits := map[string]bool{"pt": true, "mm": true, "cm": true, "in": true}
	if !validUnits[data.Unit] {
		return nil, errorsWithStack.New(fmt.Sprintf("The Unit must be pt, mm, cm or in."))
	}

	if data.MarginLeft < 0 {
		return nil, errorsWithStack.New(fmt.Sprintf("A negative MarginLeft (%f) is not allowed.", data.MarginLeft))
	}

	if data.MarginTop < 0 {
		return nil, errorsWithStack.New(fmt.Sprintf("A negative MarginTop (%f) is not allowed.", data.MarginTop))
	}

	if data.MarginRight < 0 {
		return nil, errorsWithStack.New(fmt.Sprintf("A negative MarginRight (%f) is not allowed.", data.MarginRight))
	}

	if data.MarginBottom < 0 {
		return nil, errorsWithStack.New(fmt.Sprintf("A negative MarginBottom (%f) is not allowed.", data.MarginBottom))
	}
	// <--

	// create new PDF
	pdf := gofpdf.New("P", data.Unit, "A4", "")
	if data.FontName == "OpenSans" {
		pdf.AddUTF8Font("OpenSans", "", "fonts/OpenSans-Regular.ttf")
		pdf.AddUTF8Font("OpenSans", "l", "fonts/OpenSans-Light.ttf")
		pdf.AddUTF8Font("OpenSans", "i", "fonts/OpenSans-Italic.ttf")
		pdf.AddUTF8Font("OpenSans", "b", "fonts/OpenSans-Bold.ttf")
		pdf.AddUTF8Font("OpenSans", "m", "fonts/OpenSans-Medium.ttf")
	}
	pdf.SetFont(data.FontName, "", data.FontSize)
	pdf.SetMargins(data.MarginLeft, data.MarginTop, data.MarginRight)
	pdf.SetHomeXY()
	pdf.SetAutoPageBreak(true, data.MarginBottom)
	//pdf.AliasNbPages("{entute}")
	//pdf.SetHeaderFuncMode(, true)
	//pdf.SetFooterFunc()

	pdf.SetHeaderFuncMode(headerFunction, true)
	pdf.SetFooterFuncLpi(footerFunction)
	pdf.SetHomeXY()

	if pdf.Err() {
		return nil, pdf.Error()
	}

	// create new PDFGenerator instance
	gen = new(PDFGenerator)
	gen.logger = logger
	pageWidth, pageHeight := pdf.GetPageSize()
	gen.pdf = pdf
	gen.data = data
	gen.strictErrorHandling = strictErrorHandling
	gen.maxSaveX = pageWidth - data.MarginRight
	gen.maxSaveY = pageHeight - data.MarginBottom
	gen.registeredImageTypes = map[string]string{}

	return gen, pdf.Error()
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
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	// --> validate inputs
	valideAlignStrs := map[string]bool{"L": true, "R": true, "C": true}
	if !valideAlignStrs[alignStr] {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("\"%s\" is not a valid alignStr of \"L\", \"R\" or \"C\".", alignStr)))
		return
	}

	if len(text) == 0 {
		core.logger.Warn().Msg("No text to print, return now. Please use NewLine() to print a new line.")
		//core.pdf.SetError(errorsWithStack.New("No text to print, return now."))
		return
	}

	// TODO adjust style string ("B" (bold), "I" (italic), "U" (underscore), "S" (strike-out) or any combination. The default value (specified with an empty string) is regular.)
	// <--

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
		//default:
		//	core.pdf.SetError(errorsWithStack.New("can't interpret the given text align code"))
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
//	"" non-specific font style,
//	"l" light font,
//	"i" italic font,
//	"b" bold font, or
//	"m" medium font.
//
// alignStr set the align mode:
//
//	"L" align the left side of the text to the current cursor position,
//	"R" align the right side of the text to the current cursor position, or
//	"C" align the center of the text to the current cursor position.
func (core *PDFGenerator) PrintLnPdfText(text string, styleStr string, alignStr string) {
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	// --> validate inputs
	valideAlignStrs := map[string]bool{"L": true, "R": true, "C": true}
	if !valideAlignStrs[alignStr] {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("\"%s\" is not a valid alignStr of \"L\", \"R\" or \"C\".", alignStr)))
		return
	}

	if len(text) == 0 {
		core.logger.Warn().Msg("No text to print, return now. Please use NewLine() to print a new line.")
		//core.pdf.SetError(errorsWithStack.New("No text to print, return now. Please use NewLine() to print a new line."))
		return
	}

	// TODO adjust style string ("B" (bold), "I" (italic), "U" (underscore), "S" (strike-out) or any combination. The default value (specified with an empty string) is regular.)
	// <--

	lines := core.extractLinesFromText(text)
	referenceX := core.pdf.GetX()

	for _, line := range lines {
		if line != "" {
			core.PrintPdfText(line, styleStr, alignStr)
		}
		core.NewLine(referenceX)
	}
}

// NewLine sets the cursor on the next line dependent on the given X-position.
// (mostly use the start X-point of the current line.)
func (core *PDFGenerator) NewLine(oldX float64) {
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	// --> validate inputs
	if oldX < 0 {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("A negative oldX is not allowed.")))
		return
	}
	// <--

	_, lineHeight := core.pdf.GetFontSize()
	newY := core.pdf.GetY() + lineHeight + core.data.FontGapY
	core.pdf.SetXY(oldX, newY)
}

// extractLinesFromText split a string on newline character (\n) and return the parts as an array.
// Prefixing whitespaces (ONLY " ")! will be automatically removed on each part.
func (core *PDFGenerator) extractLinesFromText(text string) (textLines []string) {
	textLines = strings.Split(text, "\n")

	// remove prefixing whitespaces (" ") from eny line.
	for i, line := range textLines {
		whitespaceCounter := 0

		for _, c := range line {
			if c != 32 {
				break
			}
			whitespaceCounter++
		}

		if whitespaceCounter > 0 {
			textLines[i] = line[whitespaceCounter:]
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
// *alignStr* set the align mode. The default alignment is left middle.
//
//	Horizontal alignment is controlled by including "L", "C" or "R" (left, center, right) in alignStr.
//	Vertical alignment is controlled by including "T", "M", "B" or "A" (top, middle, bottom, baseline) in alignStr.
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
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	// --> validate inputs

	//TODO refactor this
	//valideAlignStrs := map[string]bool{"L": true, "R": true, "C": true}
	//if !valideAlignStrs[alignStr] {
	//	core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("\"%s\" is not a valid alignStr of \"L\", \"R\" or \"C\".", alignStr)))
	//	return
	//}

	if cellHeight <= 0 {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("A negative or zero cellHeight is not allowed.")))
		return
	}
	if cellWidth <= 0 {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("A negative or zero cellHeight is not allowed.")))
		return
	}

	// TODO adjust style string ("B" (bold), "I" (italic), "U" (underscore), "S" (strike-out) or any combination. The default value (specified with an empty string) is regular.)
	// TODO valide borderStr?
	// <--

	core.pdf.SetFont(core.data.FontName, styleStr, core.GetFontSize())
	if fill {
		core.pdf.SetFillColor(int(backgroundColor.R), int(backgroundColor.G), int(backgroundColor.B))
	}
	core.pdf.CellFormat(cellWidth, cellHeight, text, borderStr, 0, alignStr, fill, 0, "")
}

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
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	// --> validate inputs
	if lineWith < 0 {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("A negative lineWith is not allowed.")))
		return
	}

	pageWidth, pageLength := core.pdf.GetPageSize()

	if x1 < 0 || x1 > pageWidth {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("x1 (%f) is out of range [%f, %f].", x1, 0.0, pageWidth)))
		return
	}

	if x2 < 0 || x2 > pageWidth {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("x2 (%f) is out of range [%f, %f].", x2, 0.0, pageWidth)))
		return
	}

	if y1 < 0 || y1 > pageLength {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("y1 (%f) is out of range [%f, %f].", y1, 0.0, pageLength)))
		return
	}

	if y2 < 0 || y2 > pageLength {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("y2 (%f) is out of range [%f, %f].", y2, 0.0, pageLength)))
		return
	}
	// <--

	core.pdf.SetLineWidth(lineWith)
	core.pdf.SetDrawColor(int(color.R), int(color.G), int(color.B))
	core.pdf.Line(x1, y1, x2, y2)
}

// RegisterMimeImageToPdf downloade a JPEG, PNG or GIF image (from mostly a Content Delivery Network (CDN)) URL
// and puts it in the current page.
// The image will be registered in the PDF but not place on a page!
// Use PlaceRegisteredImageOnPage to place the image on a page.
//
// cdnUrl specifies a parsed (CDN) URL.
//
// return imageNameStr, the image identifier for placing the image on a pdf page.
func (core *PDFGenerator) RegisterMimeImageToPdf(cdnUrl *url.URL) (imageNameStr string) {
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	var rsp *http.Response

	rsp, err := http.Get(cdnUrl.String())
	if err != nil {
		core.pdf.SetError(errorsWithStack.New(err))
		return
	}

	imageNameStr = cdnUrl.String()

	imageType := core.pdf.ImageTypeFromMime(rsp.Header["Content-Type"][0])

	switch imageType {
	case "jpg":
	case "png":
	case "gif":
		break
	default:
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("Image type is not supported.")))
		return ""
	}

	core.pdf.RegisterImageReader(cdnUrl.String(), imageType, rsp.Body)
	core.registeredImageTypes[imageNameStr] = imageType

	return imageNameStr
}

// PlaceRegisteredImageOnPage place a registered image (see RegisterMimeImageToPdf) on the current pdf page.
// The top side of the image will be snap to the current cursor position.
//
// imageNameStr specifies the registered image identifier.
//
// scale specifies the scaling factor into which the image is drawn.
// The value must be grater then 0. Use scaling of 1 for no scaling.
// E.g. a value of 0.5 means draw the image in half the size of the original
// and a value of 3 means draw the image in the triple size of the original.
//
// alignStr specifies the horizontal align type. Use:
//
//	"L" for align the left side of the image to the cursor,
//	"R" for align the right side of the image to the cursor, and
//	"C" for align the center of the image to the cursor.
func (core *PDFGenerator) PlaceRegisteredImageOnPage(imageNameStr string, alignStr string, scale float64) {
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	// --> validate inputs
	if scale == 0 {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("Image scale of 0 is not valide.")))
		return
	}

	if t := core.registeredImageTypes[imageNameStr]; t == "" {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("The image is not registerd.")))
		return
	}
	// <--

	imageInfoType := core.pdf.GetImageInfo(imageNameStr)
	posX, posY := core.GetCursor()
	imgWd, imgHt := imageInfoType.Extent()
	imgWd, imgHt = imgWd*scale, imgHt*scale

	switch alignStr {
	case "L":
		break
	case "R":
		posX = posX - imgWd
	case "C":
		posX = posX - (imgWd / 2.)
	default:
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("The image scale must be grater then 0.")))
		return
	}

	if core.pdf.Ok() {
		core.pdf.Image(imageNameStr, posX, posY, imgWd, imgHt, false, core.registeredImageTypes[imageNameStr], 0, "")
	}

	return
}

// PrintTableHeader print a generic and clean styled table header.
//
// cells contain the displayed column names of the table header.
//
// columnWidth defines the width of each column. NOTE: in general use here the same widths as in PrintTableBody()
func (core *PDFGenerator) PrintTableHeader(cells []string, columnWidth []float64, columnAlignStrings []string) {
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	// --> validate inputs
	if len(cells) != len(columnWidth) {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("The length of cells and columnWidth must be equial.")))
		return
	}

	// TODO check all columnWidths
	// <--

	referenceX := core.pdf.GetX()
	_, lineHeight := core.pdf.GetFontSize()
	newlineHeight := lineHeight + core.data.FontGapY*2

	for i, cell := range cells {
		core.PrintPdfTextFormatted(cell, "b", columnAlignStrings[i], "TB", true, Color{R: 239, G: 239, B: 239}, newlineHeight, columnWidth[i])
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
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	// --> validate inputs
	if len(columnWidths) != len(columnAlignStrings) {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("The length of columnWidths and columnAlignStrings must be equial.")))
		return
	}

	// TODO check all cells, that the length is equal to len(columnWidths) or len (columnAlignStrings)

	// <--

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
	// TODO input validation
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
	// --> validate inputs
	if len(columnWidths) != len(columnAlignStrings) {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("The length of columnWidths and columnAlignStrings must be equial.")))
		return
	}

	// TODO check all cells, that the length is equal to len(columnWidths) or len (columnAlignStrings)

	// <--
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

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

func (core *PDFGenerator) NewPage() {
	core.pdf.AddPage()
}

func (core *PDFGenerator) ComputeStringLength(str string) (length float64) {
	return core.pdf.GetStringWidth(str)
}
