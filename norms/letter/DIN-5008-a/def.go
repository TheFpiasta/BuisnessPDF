package DIN_5008_a

import "SimpleInvoice/norms/paperSize/DIN-A4"

// in mm
// for DIN A4 paper
const (
	FontSize10 = 10
	FontSize11 = 11
	FontSize12 = 12

	LineSpacing = 1.3 // line spacing factor (130%)
	FontGab10   = 3
	FontGab11   = 3.3
	FontGab12   = 3.6

	HeaderStartX = 0
	HeaderStartY = 0
	HeaderStopX  = 27
	HeaderStopY  = DIN_A4.Width

	AddressSenderTextStartX = 25
	AddressSenderTextStartY = 27
	AddressSenderTextStopX  = 105
	AddressSenderTextStopY  = 44.7

	AddressReceiverTextStartX = 25
	AddressReceiverTextStartY = 44.7
	AddressReceiverTextStopX  = 105
	AddressReceiverTextStopY  = 72

	MetaInfoStartX = 125
	MetaInfoStartY = 32
	MetaInfoStopX  = 200
	MetaInfoSTopY  = 107

	BodyStartX = 25
	BodyStartY = 103.46
	BodyStopX  = 90

	MarginPageNumberY = 4.23
)