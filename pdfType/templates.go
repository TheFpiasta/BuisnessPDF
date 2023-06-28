package pdfType

import (
	"SimpleInvoice/generator"
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

func din5008aSenderAdresse(gen *generator.PDFGenerator, senderInfo FullPersonInfo) {

}

func din5008aReceiverAdresse(gen *generator.PDFGenerator) {

}

func din5008atMetaInfo(gen *generator.PDFGenerator) {

}
