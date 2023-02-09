package generator

import "github.com/jung-kurt/gofpdf"

// GetPdf returns the full PDF.
//
// Usually used at the end of all manipulations.
func (core *PDFGenerator) GetPdf() *gofpdf.Fpdf {
	return core.pdf
}

// GetError returns the internal PDF error; this will be nil if no error has occurred.
func (core *PDFGenerator) GetError() error {
	return core.pdf.Error()
}

// GetLineHeight returns the specified line height in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) GetLineHeight() float64 {
	return core.data.LineHeight
}

// GetFontName returns the specified font name in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) GetFontName() string {
	return core.data.FontName
}

// GetMarginLeft returns the specified left margin in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) GetMarginLeft() float64 {
	return core.data.MarginLeft
}

// GetMarginTop returns the specified top margin in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) GetMarginTop() float64 {
	return core.data.MarginTop
}

// GetMarginRight returns the specified right margin in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) GetMarginRight() float64 {
	return core.data.MarginRight
}

// GetFontGapY returns the specified gap between two lines in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) GetFontGapY() float64 {
	return core.data.FontGapY
}

// GetFontSize returns the current font size in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) GetFontSize() float64 {
	return core.data.FontSize
}

// SetLineHeight change the line height in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) SetLineHeight(lineHeight float64) {
	core.data.LineHeight = lineHeight
}

// SetFontGapY change the gap between two lines in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) SetFontGapY(fontGapY float64) {
	core.data.FontGapY = fontGapY
}

// SetFontSize change the font size in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) SetFontSize(textSize float64) {
	core.data.FontSize = textSize
}
