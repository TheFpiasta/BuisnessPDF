package generator

import (
	"github.com/jung-kurt/gofpdf"
	"net/url"
	"reflect"
	"testing"
)

func TestNewPDFGenerator(t *testing.T) {
	type args struct {
		data                MetaData
		strictErrorHandling bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no error",
			args: args{
				data:                defaultMetaData,
				strictErrorHandling: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPDFGenerator(tt.args.data, tt.args.strictErrorHandling)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewPDFGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestPDFGenerator_DrawLine(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		x1       float64
		y1       float64
		x2       float64
		y2       float64
		color    Color
		lineWith float64
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
			core.DrawLine(tt.args.x1, tt.args.y1, tt.args.x2, tt.args.y2, tt.args.color, tt.args.lineWith)
		})
	}
}

func TestPDFGenerator_NewLine(t *testing.T) {
	type args struct {
		oldX float64
	}
	tests := []struct {
		name    string
		data    MetaData
		args    args
		wantErr bool
	}{
		{
			name:    "default",
			data:    defaultMetaData,
			args:    args{oldX: 15.3},
			wantErr: false,
		},
		{
			name: "font size * 3.14159",
			data: MetaData{
				FontName:     defaultMetaData.FontName,
				FontGapY:     defaultMetaData.FontGapY,
				FontSize:     defaultMetaData.FontSize * 3.14159,
				MarginLeft:   defaultMetaData.MarginLeft,
				MarginTop:    defaultMetaData.MarginTop,
				MarginRight:  defaultMetaData.MarginRight,
				MarginBottom: defaultMetaData.MarginBottom,
				Unit:         defaultMetaData.Unit,
			},
			args:    args{oldX: 15.3},
			wantErr: false,
		},
		{
			name:    "to small old x",
			data:    defaultMetaData,
			args:    args{oldX: -0.1},
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

			_, lineHeight := core.pdf.GetFontSize()
			wantNewY := tt.data.FontGapY + core.pdf.GetY() + lineHeight
			wantX := tt.args.oldX

			core.NewLine(tt.args.oldX)
			if core.pdf.Err() != tt.wantErr {
				t.Errorf("NewLine() error = %v, wantErr %v", core.pdf.Err(), tt.wantErr)
				return
			} else if core.pdf.Err() {
				wantNewY = core.pdf.GetY()
				wantX = core.pdf.GetX()
			}

			if x, y := core.pdf.GetXY(); x != wantX || y != wantNewY {
				t.Errorf("NewLine() got x y = %v %v, want %v %v", x, y, tt.args.oldX, wantNewY)
			}
		})
	}
}

func TestPDFGenerator_PlaceMimeImageFromUrl(t *testing.T) {
	type args struct {
		cdnUrl   *url.URL
		scale    float64
		alignStr string
	}
	tests := []struct {
		name string
		data MetaData
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false)
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			core.PlaceMimeImageFromUrl(tt.args.cdnUrl, tt.args.scale, tt.args.alignStr)
		})
	}
}

func TestPDFGenerator_PrintLnPdfText(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		text     string
		styleStr string
		alignStr string
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
			core.PrintLnPdfText(tt.args.text, tt.args.styleStr, tt.args.alignStr)
		})
	}
}

func TestPDFGenerator_PrintPdfText(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		text     string
		styleStr string
		alignStr string
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
			core.PrintPdfText(tt.args.text, tt.args.styleStr, tt.args.alignStr)
		})
	}
}

func TestPDFGenerator_PrintPdfTextFormatted(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		text            string
		styleStr        string
		alignStr        string
		borderStr       string
		fill            bool
		backgroundColor Color
		cellHeight      float64
		cellWidth       float64
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
			core.PrintPdfTextFormatted(tt.args.text, tt.args.styleStr, tt.args.alignStr, tt.args.borderStr, tt.args.fill, tt.args.backgroundColor, tt.args.cellHeight, tt.args.cellWidth)
		})
	}
}

func TestPDFGenerator_PrintTableBody(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		cells              [][]string
		columnWidths       []float64
		columnAlignStrings []string
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
			core.PrintTableBody(tt.args.cells, tt.args.columnWidths, tt.args.columnAlignStrings)
		})
	}
}

func TestPDFGenerator_PrintTableFooter(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		cells              [][]string
		columnWidths       []float64
		columnAlignStrings []string
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
			core.PrintTableFooter(tt.args.cells, tt.args.columnWidths, tt.args.columnAlignStrings)
		})
	}
}

func TestPDFGenerator_PrintTableHeader(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		cells              []string
		columnWidth        []float64
		columnAlignStrings []string
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
			core.PrintTableHeader(tt.args.cells, tt.args.columnWidth, tt.args.columnAlignStrings)
		})
	}
}

func TestPDFGenerator_addNewPageIfNecessary(t *testing.T) {
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
			core.addNewPageIfNecessary()
		})
	}
}

func TestPDFGenerator_extractLinesFromText(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		text string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantTextLines []string
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
			if gotTextLines := core.extractLinesFromText(tt.args.text); !reflect.DeepEqual(gotTextLines, tt.wantTextLines) {
				t.Errorf("extractLinesFromText() = %v, want %v", gotTextLines, tt.wantTextLines)
			}
		})
	}
}

func TestPDFGenerator_printTableBodyRow(t *testing.T) {
	type fields struct {
		pdf                 *gofpdf.Fpdf
		data                MetaData
		maxSaveX            float64
		maxSaveY            float64
		strictErrorHandling bool
	}
	type args struct {
		extractedLines [][]string
		currentLine    int
		maxItems       int
		alignStrings   []string
		newlineHeight  float64
		columnWidth    []float64
		referenceX     float64
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
			core.printTableBodyRow(tt.args.extractedLines, tt.args.currentLine, tt.args.maxItems, tt.args.alignStrings, tt.args.newlineHeight, tt.args.columnWidth, tt.args.referenceX)
		})
	}
}
