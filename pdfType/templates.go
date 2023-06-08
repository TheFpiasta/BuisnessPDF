package pdfType

import (
	"SimpleInvoice/generator"
	"fmt"
	errorsWithStack "github.com/go-errors/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"net/url"
)

func mimeImg(pdfGen *generator.PDFGenerator, strUrl string, posX float64, posY float64, scale float64) {
	urlStruct, err := url.Parse(strUrl)
	if err != nil {
		pdfGen.SetError(errorsWithStack.New(err.Error()))
		return
	}

	pdfGen.SetUnsafeCursor(posX, posY)
	pdfGen.PlaceMimeImageFromUrl(urlStruct, scale, "R")
}

func letterAddressSenderSmall(pdfGen *generator.PDFGenerator, address string, posX float64, posY float64, size float64) {
	pdfGen.SetCursor(posX, posY)
	pdfGen.SetFontSize(size)
	pdfGen.PrintPdfText(address, "", "L")
}

func letterFooter(pdfGen *generator.PDFGenerator, meta PdfMeta, senderInfo SenderInfo, senderAddress FullPersonInfo, lineColor generator.Color) {
	pageWidth, _ := pdfGen.GetPdf().GetPageSize()

	pdfGen.SetFontSize(meta.Font.SizeSmall)
	pdfGen.DrawLine(meta.Margin.Left, 261, pageWidth-meta.Margin.Right, 261, lineColor, 0)
	pdfGen.SetCursor(meta.Margin.Left, 264)
	pdfGen.PrintLnPdfText(senderInfo.Web, "", "L")
	pdfGen.PrintLnPdfText(senderInfo.Phone, "", "L")
	pdfGen.PrintLnPdfText(senderInfo.Email, "", "L")
	pdfGen.SetCursor(105, 264)
	pdfGen.PrintLnPdfText(senderAddress.CompanyName, "", "C")
	pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s", senderAddress.Address.Road, senderAddress.Address.HouseNumber), "", "C")
	pdfGen.PrintLnPdfText(senderAddress.Address.ZipCode+" "+senderAddress.Address.CityName, "", "C")
	pdfGen.PrintLnPdfText(senderInfo.TaxNumber, "", "C")
	pdfGen.SetCursor(190, 264)
	pdfGen.PrintLnPdfText(senderInfo.BankName, "", "R")
	pdfGen.PrintLnPdfText(senderInfo.Iban, "", "R")
	pdfGen.PrintLnPdfText(senderInfo.Bic, "", "R")
	pdfGen.DrawLine(meta.Margin.Left, 282, pageWidth-meta.Margin.Right, 282, lineColor, 0)
	pdfGen.SetFontSize(meta.Font.SizeDefault)
	pdfGen.SetCursor(pageWidth/2, 285)
	pdfGen.SetFontSize(meta.Font.SizeSmall)
	pdfGen.PrintLnPdfText("Seite 1 von 1", "", "C")
	pdfGen.SetFontSize(meta.Font.SizeDefault)
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

func germanNumber(n float64) string {
	p := message.NewPrinter(language.German)
	return p.Sprintf("%.2f", n)
}
