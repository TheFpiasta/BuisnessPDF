package generator

import "github.com/jung-kurt/gofpdf"

func (core *PDFGenerator) GetPdf() *gofpdf.Fpdf {
	return core.pdf
}

func (core *PDFGenerator) GetError() error {
	return core.pdf.Error()
}

func (core *PDFGenerator) GetLineHeight() float64 {
	return core.data.LineHeight
}

func (core *PDFGenerator) GetTextFont() string {
	return core.data.FontName
}

func (core *PDFGenerator) GetMarginLeft() float64 {
	return core.data.MarginLeft
}

func (core *PDFGenerator) GetFontGapY() float64 {
	return core.data.FontGapY
}

func (core *PDFGenerator) GetMarginTop() float64 {
	return core.data.MarginTop
}

func (core *PDFGenerator) GetMarginRight() float64 {
	return core.data.MarginRight
}

func (core *PDFGenerator) GetFontSize() float64 {
	return core.data.FontSize
}

func (core *PDFGenerator) SetLineHeight(lineHeight float64) {
	core.data.LineHeight = lineHeight
}

func (core *PDFGenerator) SetFontGapY(fontGapY float64) {
	core.data.FontGapY = fontGapY
}

func (core *PDFGenerator) SetFontSize(textSize float64) {
	core.data.FontSize = textSize
}
