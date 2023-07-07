package din5008a

import dinA4 "SimpleInvoice/norms/paperSize/din-a4"

// in mm
// for DIN A4 paper
const (
	Width  = dinA4.Width
	Height = dinA4.Height

	FontSizeSender8   = 8
	FontGabSender8    = 0.5
	FontSizeReceiver8 = 8
	FontGabReceiver8  = 1

	FontSize10 = 10.
	FontSize11 = 11.
	FontSize12 = 12.

	LineSpacing = 1.3 // line spacing factor (130%)
	FontGab10   = 3.
	FontGab11   = 3.3
	FontGab12   = 3.6

	HeaderStartX = 0.
	HeaderStartY = 0.
	HeaderStopX  = dinA4.Width
	HeaderStopY  = 27.

	AddressSenderTextStartX = 25.
	AddressSenderTextStartY = 27.
	AddressSenderTextStopX  = 105.
	AddressSenderTextStopY  = 44.7

	AddressReceiverTextStartX = 25.
	AddressReceiverTextStartY = 44.7
	AddressReceiverTextStopX  = 105
	AddressReceiverTextStopY  = 72.

	MetaInfoStartX = 125.
	MetaInfoStartY = 32.
	MetaInfoStopX  = 200.
	MetaInfoStopY  = 95

	BodyStartX = 25.
	BodyStartY = 103.46
	BodyStopX  = 190

	MarginPageNumberY = 4.23
)

type FullAdresse struct {
	FullForename string `json:"fullForename"`
	FullSurname  string `json:"fullSurname"`
	CompanyName  string `json:"companyName"`
	NameTitle    string `json:"nameTitle"`
	Address      struct {
		Road             string `json:"road"`
		HouseNumber      string `json:"houseNumber"`
		StreetSupplement string `json:"streetSupplement"`
		ZipCode          string `json:"zipCode"`
		CityName         string `json:"cityName"`
		Country          string `json:"country"`
		CountryCode      string `json:"countryCode"`
	} `json:"address"`
}
