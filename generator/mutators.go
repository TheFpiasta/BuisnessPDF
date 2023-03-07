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

// GetMarginBottom returns the specified bottom margin in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) GetMarginBottom() float64 {
	return core.data.MarginBottom
}

// GetFontGapY returns the specified gap between two lines in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) GetFontGapY() float64 {
	return core.data.FontGapY
}

// SetFontGapY change the gap between two lines in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) SetFontGapY(fontGapY float64) {
	core.data.FontGapY = fontGapY
}

// GetFontSize returns the current font size in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) GetFontSize() float64 {
	return core.data.FontSize
}

// SetFontSize change the font size in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) SetFontSize(textSize float64) {
	core.data.FontSize = textSize
}
