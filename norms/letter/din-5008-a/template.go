package din5008a

import (
	"SimpleInvoice/generator"
	"fmt"
	errorsWithStack "github.com/go-errors/errors"
	"github.com/rs/zerolog"
	"net/url"
)

type InfoData struct {
	Name  string
	Value string
}

func MimeImage(pdfGen *generator.PDFGenerator, strUrl string) {
	urlStruct, err := url.Parse(strUrl)
	if err != nil {
		pdfGen.SetError(errorsWithStack.New(err.Error()))
		return
	}
	const marginRight = Width - MetaInfoStopX
	const marginTop = 5.

	const startX = HeaderStopX - marginRight
	const startY = HeaderStartY + marginTop
	const maxImageHeight = HeaderStopY - marginTop

	pdfGen.SetUnsafeCursor(startX, startY)

	if !pdfGen.ImageIsRegistered(urlStruct.String()) {
		pdfGen.RegisterMimeImageToPdf(urlStruct)
	}

	_, imgHeight := pdfGen.GetRegisteredImageExtent(urlStruct.String())

	scale := maxImageHeight / imgHeight
	pdfGen.PlaceRegisteredImageOnPage(urlStruct.String(), "R", scale)
}

func MetaInfo(pdfGen *generator.PDFGenerator, data []InfoData) {
	var maxNameLength = 0.
	var maxValueLength = 0.

	pdfGen.SetFontSize(FontSize10)
	pdfGen.SetFontGapY(FontGab10)

	for _, datum := range data {
		nameLength := pdfGen.ComputeStringLength(datum.Name)
		if nameLength > maxNameLength {
			maxNameLength = nameLength
		}

		valueLength := pdfGen.ComputeStringLength(datum.Value)
		if valueLength > maxValueLength {
			maxValueLength = valueLength
		}
	}

	//todo check max width
	//todo check max length

	pdfGen.SetCursor(MetaInfoStartX, MetaInfoStartY)
	for _, datum := range data {
		pdfGen.PrintLnPdfText(datum.Name, "", "L")
	}

	const gapNameValue = 2
	pdfGen.SetCursor(MetaInfoStartX+maxNameLength+gapNameValue, MetaInfoStartY)

	for _, datum := range data {
		pdfGen.PrintLnPdfText(datum.Value, "", "L")
	}

	_, y := pdfGen.GetCursor()

	pdfGen.DrawLine(MetaInfoStartX, MetaInfoStartY, MetaInfoStartX, y-FontGab10)
}

func ReceiverAdresse(pdfGen *generator.PDFGenerator, receiverAddress FullAdresse) {
	pdfGen.SetCursor(AddressReceiverTextStartX, AddressReceiverTextStartY)

	pdfGen.SetFontSize(FontSize10)
	pdfGen.SetFontGapY(FontGabSender8)

	if receiverAddress.CompanyName != "" {
		pdfGen.PrintLnPdfText(receiverAddress.CompanyName, "", "L")
	}

	if receiverAddress.FullForename != "" || receiverAddress.FullSurname != "" {
		pdfGen.PrintLnPdfText(fmt.Sprintf("%s %s %s", receiverAddress.NameTitle, receiverAddress.FullForename, receiverAddress.FullSurname), "", "L")
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

func SenderAdresse(pdfGen *generator.PDFGenerator, senderInfo FullAdresse) {

	var addressSenderCompanySmall = ""

	if senderInfo.CompanyName != "" {
		addressSenderCompanySmall += fmt.Sprintf("%s", senderInfo.CompanyName)

		if senderInfo.FullForename != "" || senderInfo.FullSurname != "" {
			addressSenderCompanySmall += ", "
		}
	}

	if senderInfo.NameTitle != "" && (senderInfo.FullForename != "" || senderInfo.FullSurname != "") {
		addressSenderCompanySmall += fmt.Sprintf("%s ", senderInfo.NameTitle)
	}

	if senderInfo.FullForename != "" {
		addressSenderCompanySmall += fmt.Sprintf("%s ", senderInfo.FullForename)
	}
	if senderInfo.FullSurname != "" {
		addressSenderCompanySmall += fmt.Sprintf("%s ", senderInfo.FullSurname)
	}

	var addressSenderRoadSmall = ""

	addressSenderRoadSmall += fmt.Sprintf("%s %s",
		senderInfo.Address.Road,
		senderInfo.Address.HouseNumber,
	)

	if senderInfo.Address.StreetSupplement != "" {
		addressSenderRoadSmall += fmt.Sprintf(", %s", senderInfo.Address.StreetSupplement)
	}

	addressSenderRoadSmall += fmt.Sprintf(", %s %s",
		senderInfo.Address.ZipCode,
		senderInfo.Address.CityName,
	)

	if senderInfo.Address.CountryCode != "" {
		addressSenderRoadSmall += fmt.Sprintf(", %s", senderInfo.Address.CountryCode)
	}

	pdfGen.SetCursor(AddressSenderTextStartX, AddressSenderTextStopY)
	pdfGen.PreviousLine(AddressSenderTextStartX)

	pdfGen.SetFontSize(FontSizeSender8)
	pdfGen.SetFontGapY(FontGabSender8)
	pdfGen.PrintPdfText(addressSenderCompanySmall, "", "L")
	pdfGen.PreviousLine(AddressSenderTextStartX)
	pdfGen.PrintPdfText(addressSenderRoadSmall, "", "L")

	pdfGen.SetFontSize(FontSize10)
	pdfGen.SetFontGapY(FontGab10)
}

func Footer(content func(maxFooterHeight float64) (footerStartY float64), pdfGen *generator.PDFGenerator) (footerStartY float64, err error) {
	const startAtY = Height - MarginPageNumberY

	footerStartY = Height

	pdfGen.SetFontSize(FontSize10)
	pdfGen.SetFontGapY(FontGabReceiver8)
	pdfGen.DrawLine(BodyStartX, startAtY, BodyStopX, startAtY)
	pdfGen.SetUnsafeCursor(BodyStartX, startAtY)

	footerStartY = content(startAtY)
	if footerStartY <= BodyStartY || footerStartY > Height {
		return -1, errorsWithStack.New(fmt.Sprintf("footerStartY %.4f out of range [%.4f, %.4f)", footerStartY, BodyStartY, Width))
	}

	pdfGen.DrawLine(BodyStartX, footerStartY-1, BodyStopX, footerStartY-1)

	return footerStartY, nil
}

func PageNumbering(pdfGen *generator.PDFGenerator, footerStartY float64) {
	if pdfGen.GetTotalNumber() == 1 {
		// if pdf has only one page, no page number is required by DIN 5008 A
		return
	}

	pdfGen.SetFontSize(FontSize10)
	pdfGen.SetFontGapY(0)

	pages := pdfGen.GetTotalNumber()

	for i := 1; i <= pages; i++ {
		pdfGen.GoToPage(i)
		pdfGen.SetUnsafeCursor(BodyStopX, footerStartY-MarginPageNumberY)
		pdfGen.PreviousLine(BodyStopX)
		text := fmt.Sprintf("Seite %d von %d", i, pages)
		pdfGen.PrintPdfText(text, "", "R")
	}

	pdfGen.SetFontSize(FontSize10)
	pdfGen.SetFontGapY(FontGab10)
}

func Body(pdfGen *generator.PDFGenerator, bodyGenerationFunc func()) {
	pdfGen.SetCursor(BodyStartX, BodyStartY)
	pdfGen.SetFontSize(FontSize10)
	pdfGen.SetFontGapY(FontGab10)

	bodyGenerationFunc()
}

func ShowDebugFrames(pdfGen *generator.PDFGenerator, logger *zerolog.Logger) {
	logger.Warn().Msg("show debug frames")
	pdfGen.DrawLine(HeaderStartX, HeaderStopY, HeaderStopX, HeaderStopY)

	pdfGen.DrawLine(AddressSenderTextStartX, AddressSenderTextStartY, AddressSenderTextStartX, AddressSenderTextStopY)
	pdfGen.DrawLine(AddressSenderTextStartX, AddressSenderTextStopY, AddressSenderTextStopX, AddressSenderTextStopY)
	pdfGen.DrawLine(AddressSenderTextStopX, AddressSenderTextStartY, AddressSenderTextStopX, AddressSenderTextStopY)

	pdfGen.DrawLine(AddressReceiverTextStartX, AddressReceiverTextStartY, AddressReceiverTextStartX, AddressReceiverTextStopY)
	pdfGen.DrawLine(AddressReceiverTextStartX, AddressReceiverTextStopY, AddressReceiverTextStopX, AddressReceiverTextStopY)
	pdfGen.DrawLine(AddressReceiverTextStopX, AddressReceiverTextStartY, AddressReceiverTextStopX, AddressReceiverTextStopY)

	pdfGen.DrawLine(MetaInfoStartX, MetaInfoStartY, MetaInfoStopX, MetaInfoStartY)
	pdfGen.DrawLine(MetaInfoStartX, MetaInfoStopY, MetaInfoStopX, MetaInfoStopY)
	pdfGen.DrawLine(MetaInfoStartX, MetaInfoStartY, MetaInfoStartX, MetaInfoStopY)
	pdfGen.DrawLine(MetaInfoStopX, MetaInfoStartY, MetaInfoStopX, MetaInfoStopY)

	pdfGen.DrawLine(BodyStartX, BodyStartY, BodyStopX, BodyStartY)
	pdfGen.DrawLine(BodyStartX, BodyStartY, BodyStartX, Height-10)
	pdfGen.DrawLine(BodyStopX, BodyStartY, BodyStopX, Height-10)
}
