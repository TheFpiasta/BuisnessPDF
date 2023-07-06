package generator

import (
	"fmt"
	errorsWithStack "github.com/go-errors/errors"
	"github.com/jung-kurt/gofpdf"
)

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

// SetError set an internal PDF error.
func (core *PDFGenerator) SetError(err error) {
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	core.pdf.SetError(err)
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
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	// --> validate inputs
	if fontGapY < 0 {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("Text size must be grather or equal then 0.")))
		return
	}
	// <--

	core.data.FontGapY = fontGapY
}

// GetFontSize returns the current font size in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) GetFontSize() float64 {
	return core.data.FontSize
}

// SetFontSize change the font size in the unit of measure specified in NewPDFGenerator().
func (core *PDFGenerator) SetFontSize(textSize float64) {
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	// --> validate inputs
	if textSize <= 0 {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("Text size must be grather then 0.")))
		return
	}
	// <--

	core.data.FontSize = textSize
}

// GetCursor returns the abscissa (x) and ordinate (y) cursor point
func (core *PDFGenerator) GetCursor() (x float64, y float64) {
	return core.pdf.GetXY()
}

// SetCursor set manual the abscissa (x) and ordinate (y) reference point
// in the unit of measure specified in NewPDFGenerator() for the next operation.
// The position must be inside the writing area, restricted by the defined margins in NewPDFGenerator()
func (core *PDFGenerator) SetCursor(x float64, y float64) {
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	// --> validate inputs
	if x < core.data.MarginLeft || x > core.maxSaveX {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("New cursor position x = %f is out of range [%f, %f].", x, core.data.MarginLeft, core.maxSaveX)))
		return
	}

	if y < core.data.MarginTop || y > core.maxSaveY {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("New cursor position y = %f is out of range [%f, %f].", y, core.data.MarginTop, core.maxSaveY)))
		return
	}
	// <--

	core.pdf.SetXY(x, y)
}

// SetUnsafeCursor set manual the abscissa (x) and ordinate (y) reference point
// in the unit of measure specified in NewPDFGenerator() for the next operation.
// The position must be inside the page area, restricted by the page size.
func (core *PDFGenerator) SetUnsafeCursor(x float64, y float64) {
	if core.strictErrorHandling == true && core.pdf.Err() {
		return
	}

	// --> validate inputs
	pageWidth, pageHeight := core.pdf.GetPageSize()
	if x < 0 || x > pageWidth {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("New cursor position x = %f is out of range [%f, %f].", x, 0.0, pageWidth)))
		return
	}

	if y < 0 || y > pageHeight {
		core.pdf.SetError(errorsWithStack.New(fmt.Sprintf("New cursor position y = %f is out of range [%f, %f].", y, 0.0, pageHeight)))
		return
	}
	// <--

	core.pdf.SetXY(x, y)
}

func (core *PDFGenerator) GetRegisteredImageExtent(imageNameStr string) (w float64, h float64) {
	return core.pdf.GetImageInfo(imageNameStr).Extent()
}

func (core *PDFGenerator) ImageIsRegistered(imageNameStr string) bool {
	return core.registeredImageTypes[imageNameStr] != ""
}

func (core *PDFGenerator) GetCurrentPageNumber() int {
	return core.pdf.PageNo()
}

func (core *PDFGenerator) GetTotalNumber() int {
	return core.pdf.PageNo()
}

func (core *PDFGenerator) GoToPage(pageNumber int) {
	core.pdf.SetPage(pageNumber)
}
