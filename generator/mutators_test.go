package generator

import (
	"github.com/go-errors/errors"
	errorsWithStack "github.com/go-errors/errors"
	"reflect"
	"testing"
)

var _defaultMetaData = MetaData{
	FontName:     "Arial",
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
			name:      "integer value",
			data:      _defaultMetaData,
			setCursor: false,
			setX:      0,
			setY:      0,
			wantX:     1,
			wantY:     1,
		},
		{
			name:      "float value",
			data:      _defaultMetaData,
			setCursor: true,
			setX:      12.4,
			setY:      32.9,
			wantX:     12.4,
			wantY:     32.9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
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
			name:    "set no pdf error",
			data:    _defaultMetaData,
			setErr:  "",
			wantErr: false,
		},
		{
			name:    "set pdf error",
			data:    _defaultMetaData,
			setErr:  "test",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
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
	tests := []struct {
		name string
		data MetaData
		want float64
	}{
		{
			name: "integer value",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: _defaultMetaData.FontGapY,
		},
		{
			name: "float value",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     9.6,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: 9.6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			if got := core.GetFontGapY(); got != tt.want {
				t.Errorf("GetFontGapY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetFontName(t *testing.T) {
	tests := []struct {
		name string
		data MetaData
		want string
	}{
		{
			name: "Arial",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: _defaultMetaData.FontName,
		},
		{
			name: "Times",
			data: MetaData{
				FontName:     "Times",
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: "Times",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			if got := core.GetFontName(); got != tt.want {
				t.Errorf("GetFontName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetFontSize(t *testing.T) {
	tests := []struct {
		name string
		data MetaData
		want float64
	}{
		{
			name: "integer value",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: _defaultMetaData.FontSize,
		},
		{
			name: "float size",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     0.6,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: 0.6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			if got := core.GetFontSize(); got != tt.want {
				t.Errorf("GetFontSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetMarginBottom(t *testing.T) {
	tests := []struct {
		name string
		data MetaData
		want float64
	}{
		{
			name: "integer value",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: _defaultMetaData.MarginBottom,
		},
		{
			name: "float value",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: 8.7,
				Unit:         _defaultMetaData.Unit,
			},
			want: 8.7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			if got := core.GetMarginBottom(); got != tt.want {
				t.Errorf("GetMarginBottom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetMarginLeft(t *testing.T) {
	tests := []struct {
		name string
		data MetaData
		want float64
	}{
		{
			name: "integer value",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: _defaultMetaData.MarginLeft,
		},
		{
			name: "float value",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   6.7,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: 6.7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			if got := core.GetMarginLeft(); got != tt.want {
				t.Errorf("GetMarginLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetMarginRight(t *testing.T) {
	tests := []struct {
		name string
		data MetaData
		want float64
	}{
		{
			name: "integer value",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: _defaultMetaData.MarginRight,
		},
		{
			name: "float value",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  0.9,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: 0.9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			if got := core.GetMarginRight(); got != tt.want {
				t.Errorf("GetMarginRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetMarginTop(t *testing.T) {
	tests := []struct {
		name string
		data MetaData
		want float64
	}{
		{
			name: "integer value",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    _defaultMetaData.MarginTop,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: _defaultMetaData.MarginTop,
		},
		{
			name: "float value",
			data: MetaData{
				FontName:     _defaultMetaData.FontName,
				FontGapY:     _defaultMetaData.FontGapY,
				FontSize:     _defaultMetaData.FontSize,
				MarginLeft:   _defaultMetaData.MarginLeft,
				MarginTop:    4.4,
				MarginRight:  _defaultMetaData.MarginRight,
				MarginBottom: _defaultMetaData.MarginBottom,
				Unit:         _defaultMetaData.Unit,
			},
			want: 4.4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			if got := core.GetMarginTop(); got != tt.want {
				t.Errorf("GetMarginTop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPDFGenerator_GetPdf(t *testing.T) {
	tests := []struct {
		name string
		data MetaData
	}{
		{
			name: "default",
			data: _defaultMetaData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			want := core.pdf
			if got := core.GetPdf(); !reflect.DeepEqual(got, want) {
				t.Errorf("GetPdf() = %v, want %v", got, want)
			}
		})
	}
}

func TestPDFGenerator_SetCursor(t *testing.T) {
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name    string
		data    MetaData
		args    args
		want    args
		wantErr bool
	}{
		{
			name: "correct input",
			data: _defaultMetaData,
			args: args{
				x: 10,
				y: 8.5,
			},
			want: args{
				x: 10,
				y: 8.5,
			},
			wantErr: false,
		},
		{
			name: "to small x",
			data: _defaultMetaData,
			args: args{
				x: 0.9,
				y: 8.5,
			},
			want: args{
				x: _defaultMetaData.MarginLeft,
				y: _defaultMetaData.MarginTop,
			},
			wantErr: true,
		},
		{
			name: "to small y",
			data: _defaultMetaData,
			args: args{
				x: 10,
				y: 0.9,
			},
			want: args{
				x: _defaultMetaData.MarginLeft,
				y: _defaultMetaData.MarginTop,
			},
			wantErr: true,
		},
		{
			name: "to big x",
			data: _defaultMetaData,
			args: args{
				x: 1000,
				y: 8.5,
			},
			want: args{
				x: _defaultMetaData.MarginLeft,
				y: _defaultMetaData.MarginTop,
			},
			wantErr: true,
		},
		{
			name: "to big y",
			data: _defaultMetaData,
			args: args{
				x: 10,
				y: 850,
			},
			want: args{
				x: _defaultMetaData.MarginLeft,
				y: _defaultMetaData.MarginTop,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			core.SetCursor(tt.args.x, tt.args.y)

			if gotErr := core.pdf.Err(); gotErr != tt.wantErr {
				t.Errorf("SetCursor() set a error = %v, want %v", gotErr, tt.wantErr)
			}
			if x, y := core.pdf.GetXY(); x != tt.want.x || y != tt.want.y {
				t.Errorf("SetCursor() = %v %v, want %v %v", x, y, tt.want.x, tt.want.y)
			}
		})
	}
}

func TestPDFGenerator_SetError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		data MetaData
		args args
	}{
		{
			name: "nil err",
			data: _defaultMetaData,
			args: args{err: nil},
		},
		{
			name: "errors",
			data: _defaultMetaData,
			args: args{err: errors.New("test error")},
		},
		{
			name: "errors with stack",
			data: _defaultMetaData,
			args: args{err: errorsWithStack.New("test error")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}
			core.SetError(tt.args.err)

			if got := core.pdf.Error(); !reflect.DeepEqual(got, tt.args.err) {
				t.Errorf("GetPdf() = %v, want %v", got, tt.args.err)
			}

		})
	}
}

func TestPDFGenerator_SetFontGapY(t *testing.T) {
	type args struct {
		fontGapY float64
	}
	tests := []struct {
		name    string
		data    MetaData
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "integer value",
			data:    _defaultMetaData,
			args:    args{fontGapY: 6},
			want:    6,
			wantErr: false,
		},
		{
			name:    "float value",
			data:    _defaultMetaData,
			args:    args{fontGapY: 3.14},
			want:    3.14,
			wantErr: false,
		},
		{
			name:    "to small value",
			data:    _defaultMetaData,
			args:    args{fontGapY: -3.14},
			want:    _defaultMetaData.FontGapY,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			core.SetFontGapY(tt.args.fontGapY)
			if gotErr := core.pdf.Err(); gotErr != tt.wantErr {
				t.Errorf("SetFontGapY() get error = %v, want %v", gotErr, tt.wantErr)
			}

			if core.data.FontGapY != tt.want {
				t.Errorf("SetFontGapY() font gab y = %v, want %v", core.data.FontGapY, tt.wantErr)
			}
		})
	}
}

func TestPDFGenerator_SetFontSize(t *testing.T) {
	type args struct {
		textSize float64
	}
	tests := []struct {
		name    string
		data    MetaData
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "integer value",
			data:    _defaultMetaData,
			args:    args{textSize: 18},
			want:    18,
			wantErr: false,
		},
		{
			name:    "float value",
			data:    _defaultMetaData,
			args:    args{textSize: 2.3},
			want:    2.3,
			wantErr: false,
		},
		{
			name:    "to small value",
			data:    _defaultMetaData,
			args:    args{textSize: 0},
			want:    _defaultMetaData.FontSize,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			core.SetFontSize(tt.args.textSize)
			if gotErr := core.pdf.Err(); gotErr != tt.wantErr {
				t.Errorf("SetFontSize() get error = %v, want %v", gotErr, tt.wantErr)
			}

			if core.data.FontSize != tt.want {
				t.Errorf("SetFontSize() font size = %v, want %v", core.data.FontSize, tt.wantErr)
			}

		})
	}
}

func TestPDFGenerator_SetUnsafeCursor(t *testing.T) {
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name    string
		data    MetaData
		args    args
		want    args
		wantErr bool
	}{
		{
			name: "correct input",
			data: _defaultMetaData,
			args: args{
				x: 0.9,
				y: 12,
			},
			want: args{
				x: 0.9,
				y: 12,
			},
			wantErr: false,
		},
		{
			name: "to small x",
			data: _defaultMetaData,
			args: args{
				x: -1,
				y: 12,
			},
			want: args{
				x: _defaultMetaData.MarginLeft,
				y: _defaultMetaData.MarginTop,
			},
			wantErr: true,
		},
		{
			name: "to small y",
			data: _defaultMetaData,
			args: args{
				x: 0.9,
				y: -6.6,
			},
			want: args{
				x: _defaultMetaData.MarginLeft,
				y: _defaultMetaData.MarginTop,
			},
			wantErr: true,
		},
		{
			name: "to big x",
			data: _defaultMetaData,
			args: args{
				x: 3698,
				y: 12,
			},
			want: args{
				x: _defaultMetaData.MarginLeft,
				y: _defaultMetaData.MarginTop,
			},
			wantErr: true,
		},
		{
			name: "to big y",
			data: _defaultMetaData,
			args: args{
				x: 0.9,
				y: 5813,
			},
			want: args{
				x: _defaultMetaData.MarginLeft,
				y: _defaultMetaData.MarginTop,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, err := NewPDFGenerator(tt.data, false, &_logger, func() {}, func(isLastPage bool) {})
			if err != nil {
				t.Errorf("init core error\n%s", err.Error())
				return
			}

			core.SetUnsafeCursor(tt.args.x, tt.args.y)

			if gotErr := core.pdf.Err(); gotErr != tt.wantErr {
				t.Errorf("SetUnsafeCursor() set a error = %v, want %v", gotErr, tt.wantErr)
			}
			if x, y := core.pdf.GetXY(); x != tt.want.x || y != tt.want.y {
				t.Errorf("SetUnsafeCursor() = %v %v, want %v %v", x, y, tt.want.x, tt.want.y)
			}
		})
	}
}
