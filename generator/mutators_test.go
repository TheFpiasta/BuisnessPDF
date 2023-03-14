package generator

import (
	"github.com/go-errors/errors"
	"github.com/jung-kurt/gofpdf"
	"reflect"
	"testing"
)

var defaultMetaData = MetaData{
	FontName:     "arial",
	FontGapY:     1,
	FontSize:     1,
	MarginLeft:   1,
	MarginTop:    1,
	MarginRight:  1,
	MarginBottom: 1,
	Unit:         "mm",
}

func TestPDFGenerator_GetCursor(t *testing.T) {
	tests := []struct {
		name      string
		data      MetaData
		setCursor bool
		setX      float64
		setY      float64
		wantX     float64
		wantY     float64
	}{
		{
			name:      "default",
			data:      defaultMetaData,
			setCursor: false,
			setX:      0,
			setY:      0,
			wantX:     1,
			wantY:     1,
		},
		{
			name:      "setCursor",
			data:      defaultMetaData,
			setCursor: true,
			setX:      12.4,
			setY:      32.9,
			wantX:     12.4,
			wantY:     32.9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false)
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}
			if tt.setCursor {
				core.SetUnsafeCursor(tt.setX, tt.setY)
			}
			gotX, gotY := core.GetCursor()
			if gotX != tt.wantX {
				t.Errorf("GetCursor() gotX = %v, want %v", gotX, tt.wantX)
			}
			if gotY != tt.wantY {
				t.Errorf("GetCursor() gotY = %v, want %v", gotY, tt.wantY)
			}
		})
	}
}

func TestPDFGenerator_GetError(t *testing.T) {
	tests := []struct {
		name    string
		data    MetaData
		setErr  string
		wantErr bool
	}{
		{
			name:    "default",
			data:    defaultMetaData,
			setErr:  "",
			wantErr: false,
		},
		{
			name:    "wandErr",
			data:    defaultMetaData,
			setErr:  "test",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false)
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			if tt.setErr != "" {
				core.SetError(errors.New(tt.setErr))
			}

			if err := core.GetError(); (err != nil) != tt.wantErr {
				t.Errorf("GetError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPDFGenerator_GetFontGapY(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			if got := core.GetFontGapY(); got != tt.want {
				t.Errorf("GetFontGapY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetFontName(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			if got := core.GetFontName(); got != tt.want {
				t.Errorf("GetFontName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetFontSize(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			if got := core.GetFontSize(); got != tt.want {
				t.Errorf("GetFontSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetMarginBottom(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			if got := core.GetMarginBottom(); got != tt.want {
				t.Errorf("GetMarginBottom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetMarginLeft(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			if got := core.GetMarginLeft(); got != tt.want {
				t.Errorf("GetMarginLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetMarginRight(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			if got := core.GetMarginRight(); got != tt.want {
				t.Errorf("GetMarginRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetMarginTop(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			if got := core.GetMarginTop(); got != tt.want {
				t.Errorf("GetMarginTop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetPdf(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	tests := []struct {
		name   string
		fields fields
		want   *gofpdf.Fpdf
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			if got := core.GetPdf(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPdf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_SetCursor(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			core.SetCursor(tt.args.x, tt.args.y)
		})
	}
}

func TestPDFGenerator_SetError(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			core.SetError(tt.args.err)
		})
	}
}

func TestPDFGenerator_SetFontGapY(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		fontGapY float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			core.SetFontGapY(tt.args.fontGapY)
		})
	}
}

func TestPDFGenerator_SetFontSize(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		textSize float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			core.SetFontSize(tt.args.textSize)
		})
	}
}

func TestPDFGenerator_SetUnsafeCursor(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := &PDFGenerator{
				pdf:                 tt.fields.pdf,
				data:                tt.fields.data,
				maxSaveX:            tt.fields.maxSaveX,
				maxSaveY:            tt.fields.maxSaveY,
				strictErrorHandling: tt.fields.strictErrorHandling,
			}
			core.SetUnsafeCursor(tt.args.x, tt.args.y)
		})
	}
}
