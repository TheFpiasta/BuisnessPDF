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
    "company_name": "ACME Corp",
    "invoice_number": "123",
    "date": "2022-01-01",
    "items": [
        {
            "description": "Item 1",
            "quantity": 1,
            "price": 10
        },
        {
            "description": "Item 2",
            "quantity": 2,
            "price": 5
        }
    ]
}
```
The API will return a PDF of the invoice.

## Customization

The API will allow some customization of the invoice through query parameters.
Right now, that is not implemented.
