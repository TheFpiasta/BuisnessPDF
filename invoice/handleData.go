package invoice

import (
	"encoding/json"
	"io"
	"net/http"
)

func (iv *Invoice) parseJsonData(request *http.Request) (err error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			iv.logger.Fatal().Msg(err.Error())
		}
	}(request.Body)

	return json.NewDecoder(request.Body).Decode(&iv.pdfData)
}

// todo übernehmen von https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body ?
//func personCreate(w http.ResponseWriter, r *http.Request) {
//	// If the Content-Type header is present, check that it has the value
//	// application/json. Note that we are using the gddo/httputil/header
//	// package to parse and extract the value here, so the check works
//	// even if the client includes additional charset or boundary
//	// information in the header.
//	if r.Header.Get("Content-Type") != "" {
//		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
//		if value != "application/json" {
//			msg := "Content-Type header is not application/json"
//			http.Error(w, msg, http.StatusUnsupportedMediaType)
//			return
//		}
//	}
//
//	// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
//	// response body. A request body larger than that will now result in
//	// Decode() returning a "http: request body too large" error.
//	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
//
//	// Setup the decoder and call the DisallowUnknownFields() method on it.
//	// This will cause Decode() to return a "json: unknown field ..." error
//	// if it encounters any extra unexpected fields in the JSON. Strictly
//	// speaking, it returns an error for "keys which do not match any
//	// non-ignored, exported fields in the destination".
//	dec := json.NewDecoder(r.Body)
//	dec.DisallowUnknownFields()
//
//	var p Person
//	err := dec.Decode(&p)
//	if err != nil {
//		var syntaxError *json.SyntaxError
//		var unmarshalTypeError *json.UnmarshalTypeError
//
//		switch {
//		// Catch any syntax errors in the JSON and send an error message
//		// which interpolates the location of the problem to make it
//		// easier for the client to fix.
//		case errors.As(err, &syntaxError):
//			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
//			http.Error(w, msg, http.StatusBadRequest)
//
//		// In some circumstances Decode() may also return an
//		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
//		// is an open issue regarding this at
//		// https://github.com/golang/go/issues/25956.
//		case errors.Is(err, io.ErrUnexpectedEOF):
//			msg := fmt.Sprintf("Request body contains badly-formed JSON")
//			http.Error(w, msg, http.StatusBadRequest)
//
//		// Catch any type errors, like trying to assign a string in the
//		// JSON request body to a int field in our Person struct. We can
//		// interpolate the relevant field name and position into the error
//		// message to make it easier for the client to fix.
//		case errors.As(err, &unmarshalTypeError):
//			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
//			http.Error(w, msg, http.StatusBadRequest)
//
//		// Catch the error caused by extra unexpected fields in the request
//		// body. We extract the field name from the error message and
//		// interpolate it in our custom error message. There is an open
//		// issue at https://github.com/golang/go/issues/29035 regarding
//		// turning this into a sentinel error.
//		case strings.HasPrefix(err.Error(), "json: unknown field "):
//			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
//			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
//			http.Error(w, msg, http.StatusBadRequest)
//
//		// An io.EOF error is returned by Decode() if the request body is
//		// empty.
//		case errors.Is(err, io.EOF):
//			msg := "Request body must not be empty"
//			http.Error(w, msg, http.StatusBadRequest)
//
//		// Catch the error caused by the request body being too large. Again
//		// there is an open issue regarding turning this into a sentinel
//		// error at https://github.com/golang/go/issues/30715.
//		case err.Error() == "http: request body too large":
//			msg := "Request body must not be larger than 1MB"
//			http.Error(w, msg, http.StatusRequestEntityTooLarge)
//
//		// Otherwise default to logging the error and sending a 500 Internal
//		// Server Error response.
//		default:
//			log.Print(err.Error())
//			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//		}
//		return
//	}
//
//	// Call decode again, using a pointer to an empty anonymous struct as
//	// the destination. If the request body only contained a single JSON
//	// object this will return an io.EOF error. So if we get anything else,
//	// we know that there is additional data in the request body.
//	err = dec.Decode(&struct{}{})
//	if err != io.EOF {
//		msg := "Request body must only contain a single JSON object"
//		http.Error(w, msg, http.StatusBadRequest)
//		return
//	}
//
//	fmt.Fprintf(w, "Person: %+v", p)
//}

func (iv *Invoice) validateJsonData() (err error) {
	return nil
}
