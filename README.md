# SimpleInvoice
This project is a simple API built in Golang that allows users to create invoice PDFs. It provides a straightforward way to generate professional-looking invoices with customizable features. This project aims to be user-friendly and developer-friendly, with simple and easy-to-understand endpoints.

## How to Run

To run the Invoice PDF Creator API, follow these steps:

1. Clone the repository to your local machine.
2. Make sure you have Golang installed on your machine.
3. Run go run main.go in the terminal from the root directory of the project.
4. The server should start running on http://localhost:8080.

## Endpoints

The API has only one endpoint /pdf that takes a JSON payload with the necessary data to create an invoice. The JSON payload should have the following structure:

```json
{
  "senderAddress": {
    "fullForename" : "Paul",
    "fullSurname" : "Musterfrau",
    "companyName" :  "Musterfirma",
    "supplement" : "Raum 543",
    "address" :  {
      "road" : "Musterstraße",
      "houseNumber" : "1a",
      "streetSupplement" :  "Hinterhaus",
      "zipCode" : "12345",
      "cityName" :  "Musterstadt",
      "country" : "Deutschland",
      "countryCode" :  "DE"
    }
  },
  "receiverAddress": {
    "fullForename" : "Maria",
    "fullSurname" : "Mustermann",
    "companyName" :  "Mustermann GmbH",
    "supplement" : "Platz 1",
    "address" :  {
      "road" : "Burgplatz",
      "houseNumber" : "4",
      "streetSupplement" :  "",
      "zipCode" : "00000",
      "cityName" :  "Beste Stadt",
      "country" : "Spanien",
      "countryCode" :  "ES"
    }
  },
  "senderInfo" : {
    "phone" : "01745412112",
    "email" : "paul@musterfirma.de",
    "logoSvg" :  "https://cdn.pixabay.com/photo/2017/03/16/21/18/logo-2150297__340.png",
    "iban" : "DE12345678901234567890",
    "bic" :  "DEUTDEDB123",
    "taxNumber" : "123/456/789",
    "bankName" :  "Musterbank"
  },
  "invoiceMeta" : {
    "invoiceNumber" : "2",
    "invoiceDate" : "29.03.2023",
    "customerNumber" :  "KD222"
  },
  "invoiceBody" : {
    "openingText" :  "Sehr geehrte Damen und Herren,\n hiermit stellen wir Ihnen die Rechnung für unsere Leistungenaus.",
    "serviceTimeText" : "Leistungszeitraum: 01.01.2020 - 31.12.2020",
    "headlineText" :  "Dauerrechnung",
    "closingText" : "Wir danken für Ihr Vertrauen und freuen uns auf eine weitere Zusammenarbeit.",
    "ustNotice" :  "Nach § 19 UStG wird MEGA viel Umsatzsteuer berechnet.",
    "invoicedItems" : [
      {
        "positionNumber": "1",
        "quantity": 40.5,
        "unit": "h",
        "description": "agiles Software-Testing, System-Monitoring, \n Programmierung",
        "singlePrice": 4500,
        "currency": "€",
        "taxRate": 19
      },
      {
        "positionNumber": "1",
        "quantity": 40.5,
        "unit": "h",
        "description": "agiles Software-Testing, System-Monitoring, \n Programmierung",
        "singlePrice": 4500,
        "currency": "€",
        "taxRate": 19
      }
    ]
  }
}
```
The API will return a PDF of the invoice.

## Customization

The API will allow some customization of the invoice through query parameters.
Right now, that is not implemented.
