package pdfType

import (
	"SimpleInvoice/generator"
	din5008A "SimpleInvoice/norms/letter/DIN-5008-a"
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
	const marginRight = din5008A.Width - din5008A.MetaInfoStopX
	const marginTop = 5.

	const startX = din5008A.HeaderStopX - marginRight
	const startY = din5008A.HeaderStartY + marginTop
	const maxImageHeight = din5008A.HeaderStopY - marginTop

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

	pdfGen.SetCursor(din5008A.MetaInfoStartX, din5008A.MetaInfoStartY)
	for _, datum := range data {
		pdfGen.PrintLnPdfText(datum.name, "", "L")
	}

	const gapNameValue = 2
	pdfGen.SetCursor(din5008A.MetaInfoStartX+maxNameLength+gapNameValue, din5008A.MetaInfoStartY)

	for _, datum := range data {
		pdfGen.PrintLnPdfText(datum.value, "", "L")
	}
}

func din5008aReceiverAdresse(pdfGen *generator.PDFGenerator, receiverAddress FullPersonInfo) {
	pdfGen.SetCursor(din5008A.AddressReceiverTextStartX, din5008A.AddressReceiverTextStartY)

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

	var addressSenderCompanySmall = ""
	var addressSenderRoadSmall = ""

	addressSenderCompanySmall += senderInfo.CompanyName
	if senderInfo.CompanyName != "" && (senderInfo.FullForename != "" || senderInfo.FullSurname != "") {
		addressSenderCompanySmall += ", "
	}

	addressSenderCompanySmall += senderInfo.FullForename
	if senderInfo.FullSurname != "" {
		addressSenderCompanySmall += " "
	}
	addressSenderCompanySmall += senderInfo.FullSurname

	addressSenderRoadSmall += fmt.Sprintf(" - %s %s",
		senderInfo.Address.Road,
		senderInfo.Address.HouseNumber,
	)

	if senderInfo.Address.StreetSupplement != "" {
		addressSenderRoadSmall += ", "
		addressSenderRoadSmall += senderInfo.Address.StreetSupplement
	}

	addressSenderRoadSmall += fmt.Sprintf(", %s %s %s",
		senderInfo.Address.CountryCode,
		senderInfo.Address.ZipCode,
		senderInfo.Address.CityName,
	)

	// not really din conform now -> implement pdfGen new line to top
	pdfGen.SetCursor(din5008A.AddressSenderTextStartX, din5008A.AddressSenderTextStopY-4)
	pdfGen.SetFontSize(din5008A.FontSizeSender8)
	pdfGen.SetFontGapY(din5008A.FontGabSender8)
	pdfGen.PrintLnPdfText(addressSenderRoadSmall, "", "L")
	pdfGen.SetCursor(din5008A.AddressSenderTextStartX, din5008A.AddressSenderTextStopY-7)
	pdfGen.PrintLnPdfText(addressSenderCompanySmall, "", "L")

}
