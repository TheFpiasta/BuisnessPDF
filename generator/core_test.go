package generator

import (
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

var _logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.DebugLevel).With().Timestamp().Logger()

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
				data:                _defaultMetaData,
				strictErrorHandling: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPDFGenerator(tt.args.data, tt.args.strictErrorHandling, &_logger)

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
			data:    _defaultMetaData,
			args:    args{oldX: 15.3},
			wantErr: false,
		},
		{
			name: "font size * 3.14159",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize * 3.14159,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			args:    args{oldX: 15.3},
			wantErr: false,
		},
		{
			name:    "to small old x",
			data:    _defaultMetaData,
			args:    args{oldX: -0.1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger)
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
			core, err := NewPDFGenerator(tt.data, false, &_logger)
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			core.PlaceMimeImageFromUrl(tt.args.cdnUrl, tt.args.scale, tt.args.alignStr)
		})
	}
}

func TestPDFGenerator_PrintLnPdfText(t *testing.T) {
	type args struct {
		text     string
		styleStr string
		alignStr string
	}

	tests := []struct {
		name    string
		data    MetaData
		args    args
		wantErr bool
	}{
		{
			name: "normal print",
			data: _defaultMetaData,
			args: args{
				text:     "Test abc",
				styleStr: "b",
				alignStr: "L",
			},
			wantErr: false,
		},
		{
			name: "no text",
			data: _defaultMetaData,
			args: args{
				text:     "Test abc",
				styleStr: "b",
				alignStr: "L",
			},
			wantErr: true,
		},
		{
			name: "long text",
			data: _defaultMetaData,
			args: args{
				text: "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, " +
					"sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, " +
					"sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. " +
					"Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
				styleStr: "b",
				alignStr: "L",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger)
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			wantX, wantY := core.pdf.GetXY()
			_, lineHeight := core.pdf.GetFontSize()
			wantY += lineHeight + core.data.FontGapY

			core.PrintLnPdfText(tt.args.text, tt.args.styleStr, tt.args.alignStr)
			if core.GetError() != nil && !tt.wantErr {
				t.Error(err.Error())
				return
			}

			gotX, gotY := core.pdf.GetXY()

			if wantX != gotX || wantY != gotY {
				t.Errorf("PrintLnPdfText() got x y = %v %v, want %v %v", gotX, gotY, wantX, wantY)
			}
		})
	}
}

func TestPDFGenerator_PrintPdfText(t *testing.T) {
	type args struct {
		text     string
		styleStr string
		alignStr string
	}

	tests := []struct {
		name    string
		data    MetaData
		args    args
		wantErr bool
	}{
		{
			name: "normal print",
			data: _defaultMetaData,
			args: args{
				text:     "Test abc",
				styleStr: "b",
				alignStr: "L",
			},
			wantErr: false,
		},
		{
			name: "no text",
			data: _defaultMetaData,
			args: args{
				text:     "Test abc",
				styleStr: "b",
				alignStr: "L",
			},
			wantErr: true,
		},
		{
			name: "long text",
			data: _defaultMetaData,
			args: args{
				text: "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, " +
					"sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, " +
					"sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. " +
					"Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
				styleStr: "b",
				alignStr: "L",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger)
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			currentX, currentY := core.pdf.GetXY()

			core.PrintPdfText(tt.args.text, tt.args.styleStr, tt.args.alignStr)
			if core.GetError() != nil && !tt.wantErr {
				t.Error(err.Error())
				return
			}

			gotX, gotY := core.pdf.GetXY()

			if currentX >= gotX || currentY != gotY {
				t.Errorf("PrintLnPdfText() got x y = %v %v, want x > %v and y = %v", gotX, gotY, currentX, currentY)
			}
		})
	}
}

func TestPDFGenerator_PrintPdfTextFormatted(t *testing.T) {
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
		name    string
		data    MetaData
		args    args
		wantErr bool
	}{
		{
			name: "normal print",
			data: _defaultMetaData,
			args: args{
				text:            "Test abc",
				styleStr:        "b",
				alignStr:        "L",
				borderStr:       "1",
				fill:            false,
				backgroundColor: Color{0, 0, 58},
				cellHeight:      60,
				cellWidth:       100,
			},
			wantErr: false,
		},
		{
			name: "no text",
			data: _defaultMetaData,
			args: args{
				text:            "Test abc",
				styleStr:        "b",
				alignStr:        "L",
				borderStr:       "",
				fill:            false,
				backgroundColor: Color{100, 100, 100},
				cellHeight:      10,
				cellWidth:       10,
			},
			wantErr: false,
		},
		{
			name: "long text",
			data: _defaultMetaData,
			args: args{
				text: "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, " +
					"sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, " +
					"sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. " +
					"Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
				styleStr:        "b",
				alignStr:        "L",
				borderStr:       "1L",
				fill:            true,
				backgroundColor: Color{30, 48, 64},
				cellHeight:      80,
				cellWidth:       90,
			},
			wantErr: false,
		},
		{
			name: "wrong cellHeight",
			data: _defaultMetaData,
			args: args{
				text:            "Test",
				styleStr:        "",
				alignStr:        "",
				borderStr:       "",
				fill:            false,
				backgroundColor: Color{30, 48, 64},
				cellHeight:      0,
				cellWidth:       90,
			},
			wantErr: true,
		},
		{
			name: "wrong cellWidth",
			data: _defaultMetaData,
			args: args{
				text:            "Test",
				styleStr:        "",
				alignStr:        "",
				borderStr:       "",
				fill:            false,
				backgroundColor: Color{30, 48, 64},
				cellHeight:      10,
				cellWidth:       0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger)
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			wantX, wantY := core.pdf.GetXY()
			wantX += tt.args.cellWidth

			core.PrintPdfTextFormatted(tt.args.text, tt.args.styleStr, tt.args.alignStr, tt.args.borderStr, tt.args.fill, tt.args.backgroundColor, tt.args.cellHeight, tt.args.cellWidth)
			if core.GetError() != nil && !tt.wantErr {
				t.Error(err.Error())
				return
			}

			if core.GetError() != nil && tt.wantErr {
				return
			}

			gotX, gotY := core.pdf.GetXY()

			if wantX != gotX || wantY != gotY {
				t.Errorf("PrintLnPdfText() got x y = %v %v, want %v %v", gotX, gotY, wantX, wantY)
			}
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

//addNewPageIfNecessary currently not implemented
//func TestPDFGenerator_addNewPageIfNecessary(t *testing.T) {
//	type fields struct {
//		pdf                 *gofpdf.Fpdf
//		data                MetaData
//		maxSaveX            float64
//		maxSaveY            float64
//		strictErrorHandling bool
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			core := &PDFGenerator{
//				pdf:                 tt.fields.pdf,
//				data:                tt.fields.data,
//				maxSaveX:            tt.fields.maxSaveX,
//				maxSaveY:            tt.fields.maxSaveY,
//				strictErrorHandling: tt.fields.strictErrorHandling,
//			}
//			core.addNewPageIfNecessary()
//		})
//	}
//}

func TestPDFGenerator_extractLinesFromText(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name          string
		data          MetaData
		args          args
		wantTextLines []string
	}{
		{
			name:          "default",
			data:          _defaultMetaData,
			args:          args{text: "Hi\nFrom \nThe\n Test \n !!!1!11\n"},
			wantTextLines: []string{"Hi", "From ", "The", "Test ", "!!!1!11", ""},
		},
		{
			name:          "nothing to do",
			data:          _defaultMetaData,
			args:          args{text: ""},
			wantTextLines: []string{""},
		},
		{
			name:          "only \\n",
			data:          _defaultMetaData,
			args:          args{text: "\n"},
			wantTextLines: []string{"", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger)
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			if gotTextLines := core.extractLinesFromText(tt.args.text); !reflect.DeepEqual(gotTextLines, tt.wantTextLines) {
				var gotTmp []string
				var wantTmp []string
				for _, line := range gotTextLines {
					gotTmp = append(gotTmp, "\""+line+"\"")
				}
				for _, line := range tt.wantTextLines {
					wantTmp = append(wantTmp, "\""+line+"\"")
				}
				t.Errorf("extractLinesFromText() = \n%v, want \n%v", "["+strings.Join(gotTmp, ", ")+"]", "["+strings.Join(wantTmp, ", ")+"]")
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
