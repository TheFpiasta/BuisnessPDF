# SimpleInvoice
This project is a simple API built in Golang that allows users to create invoice PDFs. It provides a straightforward way to generate professional-looking invoices with customizable features. This project aims to be user-friendly and developer-friendly, with simple and easy-to-understand endpoints.

## How to Run

To run the Invoice PDF Creator API, follow these steps:

1. Install Go Version 1.18.x from the official webpage [https://go.dev/](https://go.dev/)
2. Clone the repository to your local machine. ```git clone https://github.com/TheFpiasta/SimpleInvoice.git```
3. Run ``go run main.go`` in the terminal from the root directory of the project.
4. The server should start running on [http://localhost:10000](http://localhost:10000).

## Endpoints

Once the application is running, you can send a POST request to the /pdf/ endpoint with a JSON body containing the necessary invoice details. The following is an example of a valid JSON request body:

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
    "mimeLogoUrl" :  "https://cdn.pixabay.com/photo/2017/03/16/21/18/logo-2150297__340.png",
    "mimeLogoScale": 0.5,
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

If you would like to contribute to SimpleInvoice, please fork the repository and submit a pull request. Contributions are always welcome, and we appreciate your help in making SimpleInvoice even better!

## License
SimpleInvoice is released under the MIT license. See 'LICENSE' for more information.