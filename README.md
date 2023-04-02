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
    "companyName" :  "Musterfraufirma GmbH",
    "supplement" : "",
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
    "companyName" :  "Mustermannfirma GmbH",
    "supplement" : "Platz 1",
    "address" :  {
      "road" : "Burgplatz",
      "houseNumber" : "4",
      "streetSupplement" :  "",
      "zipCode" : "12345",
      "cityName" :  "Musterstadt",
      "country" : "Deutschland",
      "countryCode" :  "DE"
    }
  },
  "senderInfo" : {
    "phone" : "01234567890",
    "web" : "test.test",
    "email" : "hello@test.test",
    "mimeLogoUrl" :  "https://cdn.pictro.de/Test/logoTest.png",
    "mimeLogoScale": 0.25,
    "iban" : "DE02100100100006820101",
    "bic" :  "PBNKDEFF",
    "taxNumber" : "123/456/789",
    "bankName" :  "POSTBANK"
  },
  "invoiceMeta" : {
    "invoiceNumber" : "Re23",
    "invoiceDate" : "02.02.2023",
    "customerNumber" :  "NEKO11"
  },
  "invoiceBody" : {
    "openingText" :  "Sehr geehrte Damen und Herren,\n hiermit stellen wir Ihnen die Rechnung für unsere Leistungenaus.",
    "serviceTimeText" : "Leistungszeitraum: 01.01.2023 - 31.02.2023",
    "headlineText" :  "Rechnung",
    "closingText" : "Wir danken für Ihr Vertrauen und freuen uns auf eine weitere Zusammenarbeit.",
    "ustNotice" :  "",
    "invoicedItems" : [
      {
        "positionNumber": "1",
        "quantity": 40.5,
        "unit": "h",
        "description": "Programmierung",
        "singlePrice": 8900,
        "currency": "€",
        "taxRate": 19
      },
      {
        "positionNumber": "1",
        "quantity": 20.7,
        "unit": "h",
        "description": "Software-Testing, System-Monitoring",
        "singlePrice": 6800,
        "currency": "€",
        "taxRate": 19
      },
      {
        "positionNumber": "1",
        "quantity": 33.1,
        "unit": "h",
        "description": "IT-Beratung",
        "singlePrice": 7500,
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
