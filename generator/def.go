package generator

import "github.com/jung-kurt/gofpdf"

type PDFGenerator struct {
	pdf  *gofpdf.Fpdf
	data MetaData
}

type MetaData struct {
	LineHeight  float64
	FontName    string
	FontGapY    float64
	FontSize    float64
	MarginLeft  float64
	MarginTop   float64
	MarginRight float64
	Unit        string
}

type Generator interface {
	PrintPdfText(text string, styleStr string, alignStr string)
	PrintLnPdfText(text string, styleStr string, alignStr string)
	DrawPdfTextRightAligned(posXRight float64, posY float64, text string, styleStr string, textSize float64, elementWith float64, elementHeight float64)
	DrawLine(x1 float64, y1 float64, x2 float64, y2 float64, color Color, lineWith float64)
	PlaceImgOnPosXY(logoUrl string, posX int, posY int) (err error)
	PrintPdfTextFormatted(text string, styleStr string, alignStr string, borderStr string, fill bool, backgroundColor Color, lineHeight float64, stringWidth float64)
	PrintTable(header []string, columnWidth []float64, items [][]string, summary [][]string, summaryWidths [][]float64)

	GetError() error
	GetPdf() *gofpdf.Fpdf
	SetCursor(x float64, y float64)

	GetLineHeight() float64
	GetTextFont() string
	GetMarginLeft() float64
	GetFontGapY() float64
	GetMarginTop() float64
	GetMarginRight() float64
	GetFontSize() float64

	SetLineHeight(lineHeight float64)
	SetFontGapY(fontGapY float64)
	SetFontSize(textSize float64)
}

type Color struct {
	R uint8
	G uint8
	B uint8
}
