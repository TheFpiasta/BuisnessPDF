package generator

import (
	"github.com/jung-kurt/gofpdf"
	"net/url"
)

// PDFGenerator is a light-way PDF document generator witch simplify and enhanced [github.com/jung-kurt/gofpdf].
// The implemented methods focused primary on creating easy clean invoices or B2B letters.
type PDFGenerator struct {
	pdf                 *gofpdf.Fpdf
	data                MetaData
	maxSaveX            float64
	maxSaveY            float64
	strictErrorHandling bool
}

// MetaData sums all necessary inputs for NewPDFGenerator().
//
// LineHeight defines the total height of a text line in the Unit of measure.
//
// FontName define font familie used to print character strings. Standard families (case insensitive):
//
//	"Courier" for fixed-width,
//	"Helvetica" or "Arial" for sans serif,
//	"Times" for serif,
//	"Symbol" or "ZapfDingbats" for symbolic.
//	"OpenSans" for TrueType support with utf-8 symbols.
//
// FontGapY defines the gap between two text lines in the Unit of measure.
//
// FontSize defines the font size measured in points.
//
// MarginLeft defines the left page margin in the Unit of measure.
//
// MarginTop defines the top page margin in the Unit of measure.
//
// MarginRight defines the right page margin in the Unit of measure.
// If the value is less than zero, it is set to the same as the left margin.
//
// MarginBottom defines the bottom page margin in the Unit of measure.
// On top of the bottom margin is the footer section.
//
// Unit specifies the unit of length used in size parameters for elements other than fonts,
// which are always measured in points. An empty string will be replaced with "mm". Specify
//
//	"pt" for point,
//	"mm" for millimeter,
//	"cm" for centimeter, or
//	"in" for inch.
type MetaData struct {
	LineHeight   float64
	FontName     string
	FontGapY     float64
	FontSize     float64
	MarginLeft   float64
	MarginTop    float64
	MarginRight  float64
	MarginBottom float64
	Unit         string
}

// Generator specify all public methods.
type Generator interface {
	PrintPdfText(text string, styleStr string, alignStr string)
	PrintLnPdfText(text string, styleStr string, alignStr string)
	DrawLine(x1 float64, y1 float64, x2 float64, y2 float64, color Color, lineWith float64)
	PlaceMimeImageFromUrl(cdnUrl *url.URL, posX float64, posY float64, scale float64) (err error)
	PrintPdfTextFormatted(text string, styleStr string, alignStr string, borderStr string, fill bool, backgroundColor Color, cellHeight float64, cellWidth float64)

	PrintTableHeader(cells []string, columnWidth []float64)
	PrintTableBody(cells [][]string, columnWidths []float64, columnAlignStrings []string)
	PrintTableFooter(cells [][]string, columnWidths []float64, columnAlignStrings []string)

	GetError() error
	GetPdf() *gofpdf.Fpdf
	SetCursor(x float64, y float64)
	SetUnsafeCursor(x float64, y float64)
	NewLine(oldX float64)

	GetLineHeight() float64
	SetLineHeight(lineHeight float64)
	GetFontName() string
	GetFontGapY() float64
	SetFontGapY(fontGapY float64)
	GetMarginLeft() float64
	GetMarginTop() float64
	GetMarginRight() float64
	GetFontSize() float64
	SetFontSize(textSize float64)
	GetMarginBottom() float64
}

// Color represents a specific color in red, green and blue values, each from 0 to 255
type Color struct {
	R uint8
	G uint8
	B uint8
}
