package pdfType

import (
	"SimpleInvoice/generator"
	"fmt"
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

func getCellWith(pdfGen *generator.PDFGenerator, percent float64) float64 {
	maxSavePrintingWidth, _ := pdfGen.GetPdf().GetPageSize()
	maxSavePrintingWidth = maxSavePrintingWidth - pdfGen.GetMarginLeft() - pdfGen.GetMarginRight()

	return (percent * maxSavePrintingWidth) / 100.0
}
