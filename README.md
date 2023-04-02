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

If you would like to contribute to SimpleInvoice, please fork the repository and submit a pull request. Contributions are always welcome, and we appreciate your help in making SimpleInvoice even better!

## License
SimpleInvoice is released under the MIT license. See 'LICENSE' for more information.