package pdfType

import (
	"SimpleInvoice/generator"
	din5008A "SimpleInvoice/norms/letter/DIN-5008-a"
	"fmt"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func getAddressLine(address FullPersonInfo) string {
	var addressSenderSmallText = ""

	addressSenderSmallText += address.CompanyName
	if address.CompanyName != "" && (address.FullForename != "" || address.FullSurname != "") {
		addressSenderSmallText += ", "
	}

	addressSenderSmallText += address.FullForename
	if address.FullSurname != "" {
		addressSenderSmallText += " "
	}
	addressSenderSmallText += address.FullSurname

	addressSenderSmallText += fmt.Sprintf(" - %s %s",
		address.Address.Road,
		address.Address.HouseNumber,
	)

	if address.Address.StreetSupplement != "" {
		addressSenderSmallText += ", "
		addressSenderSmallText += address.Address.StreetSupplement
	}

	addressSenderSmallText += fmt.Sprintf(", %s %s %s",
		address.Address.CountryCode,
		address.Address.ZipCode,
		address.Address.CityName,
	)

	return addressSenderSmallText
}

func germanNumber[T float64 | int](n T) string {
	p := message.NewPrinter(language.German)

	switch fmt.Sprintf("%T", *new(T)) {
	case "float64":
		return p.Sprintf("%.2f", n)
	case "int":
		return p.Sprintf("%d", n)
	default:
		return "GERMAN NUMBER FAILED"
	}
}

func getCellWith(pdfGen *generator.PDFGenerator, percent float64) float64 {
	maxSavePrintingWidth, _ := pdfGen.GetPdf().GetPageSize()
	maxSavePrintingWidth = maxSavePrintingWidth - pdfGen.GetMarginLeft() - pdfGen.GetMarginRight()

	return (percent * maxSavePrintingWidth) / 100.0
}

func showDebugFramesDin5008A(pdfGen *generator.PDFGenerator, logger *zerolog.Logger) {
	logger.Warn().Msg("show debug frames")
	pdfGen.DrawLine(din5008A.HeaderStartX, din5008A.HeaderStopY, din5008A.HeaderStopX, din5008A.HeaderStopY, generator.Color{R: 255, G: 64, B: 64}, 0)

	pdfGen.DrawLine(din5008A.AddressSenderTextStartX, din5008A.AddressSenderTextStartY, din5008A.AddressSenderTextStartX, din5008A.AddressSenderTextStopY, generator.Color{R: 255, G: 64, B: 64}, 0)
	pdfGen.DrawLine(din5008A.AddressSenderTextStartX, din5008A.AddressSenderTextStopY, din5008A.AddressSenderTextStopX, din5008A.AddressSenderTextStopY, generator.Color{R: 255, G: 64, B: 64}, 0)
	pdfGen.DrawLine(din5008A.AddressSenderTextStopX, din5008A.AddressSenderTextStartY, din5008A.AddressSenderTextStopX, din5008A.AddressSenderTextStopY, generator.Color{R: 255, G: 64, B: 64}, 0)

	pdfGen.DrawLine(din5008A.AddressReceiverTextStartX, din5008A.AddressReceiverTextStartY, din5008A.AddressReceiverTextStartX, din5008A.AddressReceiverTextStopY, generator.Color{R: 255, G: 64, B: 64}, 0)
	pdfGen.DrawLine(din5008A.AddressReceiverTextStartX, din5008A.AddressReceiverTextStopY, din5008A.AddressReceiverTextStopX, din5008A.AddressReceiverTextStopY, generator.Color{R: 255, G: 64, B: 64}, 0)
	pdfGen.DrawLine(din5008A.AddressReceiverTextStopX, din5008A.AddressReceiverTextStartY, din5008A.AddressReceiverTextStopX, din5008A.AddressReceiverTextStopY, generator.Color{R: 255, G: 64, B: 64}, 0)

	pdfGen.DrawLine(din5008A.MetaInfoStartX, din5008A.MetaInfoStartY, din5008A.MetaInfoStopX, din5008A.MetaInfoStartY, generator.Color{R: 255, G: 64, B: 64}, 0)
	pdfGen.DrawLine(din5008A.MetaInfoStartX, din5008A.MetaInfoStopY, din5008A.MetaInfoStopX, din5008A.MetaInfoStopY, generator.Color{R: 255, G: 64, B: 64}, 0)
	pdfGen.DrawLine(din5008A.MetaInfoStartX, din5008A.MetaInfoStartY, din5008A.MetaInfoStartX, din5008A.MetaInfoStopY, generator.Color{R: 255, G: 64, B: 64}, 0)
	pdfGen.DrawLine(din5008A.MetaInfoStopX, din5008A.MetaInfoStartY, din5008A.MetaInfoStopX, din5008A.MetaInfoStopY, generator.Color{R: 255, G: 64, B: 64}, 0)

	pdfGen.DrawLine(din5008A.BodyStartX, din5008A.BodyStartY, din5008A.BodyStopX, din5008A.BodyStartY, generator.Color{R: 255, G: 64, B: 64}, 0)
	pdfGen.DrawLine(din5008A.BodyStartX, din5008A.BodyStartY, din5008A.BodyStartX, din5008A.Height-10, generator.Color{R: 255, G: 64, B: 64}, 0)
	pdfGen.DrawLine(din5008A.BodyStopX, din5008A.BodyStartY, din5008A.BodyStopX, din5008A.Height-10, generator.Color{R: 255, G: 64, B: 64}, 0)
}
