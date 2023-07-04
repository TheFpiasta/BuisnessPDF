package pdfType

import (
	"SimpleInvoice/generator"
	DIN_5008_a "SimpleInvoice/norms/letter/DIN-5008-a"
	"fmt"
	errorsWithStack "github.com/go-errors/errors"
	"net/url"
)

func mimeImg(pdfGen *generator.PDFGenerator, strUrl string, posX float64, posY float64, scale float64) {
	urlStruct, err := url.Parse(strUrl)
	if err != nil {
		pdfGen.SetError(errorsWithStack.New(err.Error()))
		return
	}

	pdfGen.SetUnsafeCursor(posX, posY)

	if !pdfGen.ImageIsRegistered(urlStruct.String()) {
		pdfGen.RegisterMimeImageToPdf(urlStruct)
	}
	pdfGen.PlaceRegisteredImageOnPage(urlStruct.String(), "R", scale)
}

func letterAddressSenderSmall(pdfGen *generator.PDFGenerator, address string, posX float64, posY float64, size float64) {
	pdfGen.SetCursor(posX, posY)
	pdfGen.SetFontSize(size)
	pdfGen.PrintPdfText(address, "", "L")
}

func letterReceiverAddress(pdfGen *generator.PDFGenerator, receiverAddress FullPersonInfo, posX float64, posY float64) {
	pdfGen.SetCursor(posX, posY)
	if receiverAddress.CompanyName != "" {
		pdfGen.PrintLnPdfText(receiverAddress.CompanyName, "", "L")

	}
	if receiverAddress.FullForename != "" || receiverAddress.FullSurname != "" {
		pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", receiverAddress.FullForename, receiverAddress.FullSurname),
			"", "L")
	}
	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", receiverAddress.Address.Road, receiverAddress.Address.HouseNumber),
		"", "L")
	if receiverAddress.Address.StreetSupplement != "" {
		pdfGen.PrintLnPdfText(receiverAddress.Address.StreetSupplement, "", "L")
	}
	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", receiverAddress.Address.ZipCode, receiverAddress.Address.CityName),
		"", "L")
}

func getColumnWithFromPercentage(pdfGen *generator.PDFGenerator, columnPercent []float64) (columnWidth []float64) {
	for _, p := range columnPercent {
		columnWidth = append(columnWidth, getCellWith(pdfGen, p))
	}
	return columnWidth
}

func din5008aMimeImage(pdfGen *generator.PDFGenerator, strUrl string) {
	urlStruct, err := url.Parse(strUrl)
	if err != nil {
		pdfGen.SetError(errorsWithStack.New(err.Error()))
		return
	}
	const marginRight = DIN_5008_a.Width - DIN_5008_a.MetaInfoStopX
	const marginTop = 5.

	const startX = DIN_5008_a.HeaderStopX - marginRight
	const startY = DIN_5008_a.HeaderStartY + marginTop
	const maxImageHeight = DIN_5008_a.HeaderStopY - marginTop

	pdfGen.SetUnsafeCursor(startX, startY)

	if !pdfGen.ImageIsRegistered(urlStruct.String()) {
		pdfGen.RegisterMimeImageToPdf(urlStruct)
	}

	_, imgHeight := pdfGen.GetRegisteredImageExtent(urlStruct.String())

	scale := maxImageHeight / imgHeight
	pdfGen.PlaceRegisteredImageOnPage(urlStruct.String(), "R", scale)
}

func din5008atMetaInfo(pdfGen *generator.PDFGenerator, data []struct {
	name  string
	value string
}) {
	var maxNameLength = 0.
	var maxValueLength = 0.

	for _, datum := range data {
		nameLength := pdfGen.ComputeStringLength(datum.name)
		if nameLength > maxNameLength {
			maxNameLength = nameLength
		}

		valueLength := pdfGen.ComputeStringLength(datum.value)
		if valueLength > maxValueLength {
			maxValueLength = valueLength
		}
	}

	//todo check max width
	//todo check max length

	pdfGen.SetCursor(DIN_5008_a.MetaInfoStartX, DIN_5008_a.MetaInfoStartY)
	for _, datum := range data {
		pdfGen.PrintLnPdfText(datum.name, "", "L")
	}

	const gapNameValue = 2
	pdfGen.SetCursor(DIN_5008_a.MetaInfoStartX+maxNameLength+gapNameValue, DIN_5008_a.MetaInfoStartY)

	for _, datum := range data {
		pdfGen.PrintLnPdfText(datum.value, "", "L")
	}
}

func din5008aReceiverAdresse(pdfGen *generator.PDFGenerator, receiverAddress FullPersonInfo) {

	pdfGen.SetCursor(DIN_5008_a.AddressReceiverTextStartX, DIN_5008_a.AddressReceiverTextStartY)

	if receiverAddress.CompanyName != "" {
		pdfGen.PrintLnPdfText(receiverAddress.CompanyName, "", "L")
	}

	if receiverAddress.FullForename != "" || receiverAddress.FullSurname != "" {
		pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s %s", receiverAddress.Supplement, receiverAddress.FullForename, receiverAddress.FullSurname), "", "L")
	}

	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", receiverAddress.Address.Road, receiverAddress.Address.HouseNumber), "", "L")

	if receiverAddress.Address.StreetSupplement != "" {
		pdfGen.PrintLnPdfText(receiverAddress.Address.StreetSupplement, "", "L")
	}

	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", receiverAddress.Address.ZipCode, receiverAddress.Address.CityName), "", "L")

	if receiverAddress.Address.Country != "" {
		pdfGen.PrintLnPdfText(fmt.Sprintf("%s", receiverAddress.Address.Country), "", "L")
	}
}

func din5008aSenderAdresse(pdfGen *generator.PDFGenerator, senderInfo FullPersonInfo) {

}
