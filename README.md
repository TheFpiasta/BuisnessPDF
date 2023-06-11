# SimpleInvoice
This project is a simple API built in Golang that allows users to create invoice PDFs. It provides a straightforward way to generate professional-looking invoices with customizable features. This project aims to be user-friendly and developer-friendly, with simple and easy-to-understand endpoints.

## How to Run

To run the Invoice PDF Creator API, follow these steps:

1. Install Go Version 1.18.x from the official webpage [https://go.dev/](https://go.dev/)
2. Clone the repository to your local machine. ```git clone https://github.com/TheFpiasta/SimpleInvoice.git```
3. Run ``go run main.go`` in the terminal from the root directory of the project.
4. The server should start running on [http://localhost:10000](http://localhost:10000).

## Endpoints

Once the application is running, you can send a POST request a JSON body containing the necessary invoice details.
The following pdf types are implemented:

| API endpoint   | Description                 | JSON body                                                                                             |
|----------------|-----------------------------|-------------------------------------------------------------------------------------------------------|
| /invoice       | to generate a invoice       | [template](pdfType/pdfInvoiceTemplate.json) <br/> [example](pdfType/pdfInvoiceExample.json)           |
| /delivery-node | to generate a delivery node | [template](pdfType/pdfDeliveryNoteTemplate.json) <br/> [example](pdfType/pdfDeliveryNoteExample.json) |

The API will return a PDF if no error occurred, or the error message in json format.

## Customization

If you would like to contribute to SimpleInvoice, please fork the repository and submit a pull request. Contributions
are always welcome, and we appreciate your help in making SimpleInvoice even better!

## License

SimpleInvoice is released under the MIT license. See 'LICENSE' for more information.
