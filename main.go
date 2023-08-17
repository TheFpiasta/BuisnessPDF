package main

import (
	"SimpleInvoice/pdfType"
	"fmt"
	errorsWithStack "github.com/go-errors/errors"
	"github.com/rs/zerolog"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var logger zerolog.Logger

func invoiceRequest(w http.ResponseWriter, r *http.Request) {
	h := pdfType.NewInvoice(&logger)
	executeHandler(h, w, r)
}

func deliveryNodeRequest(w http.ResponseWriter, r *http.Request) {
	h := pdfType.NewDeliveryNode(&logger)
	executeHandler(h, w, r)
}

func attachmentTableRequest(w http.ResponseWriter, r *http.Request) {
	h := pdfType.NewTableAttachment(&logger)
	executeHandler(h, w, r)
}

func handleRequests() {
	http.HandleFunc("/invoice", invoiceRequest)
	http.HandleFunc("/delivery-node", deliveryNodeRequest)
	http.HandleFunc("/attachment/table", attachmentTableRequest)
	logger.Debug().Msg("start server on localhost:10000")
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {

	const loggingLevel = 0
	const logDir = ""
	const openBrowserOnStartup = false

	err := initLogger(loggingLevel, logDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	if openBrowserOnStartup {
		go openBrowser("http://localhost:10000/")
	}

	handleRequests()
}

func executeHandler(handler pdfType.PdfType, w http.ResponseWriter, r *http.Request) {
	err := handler.SetDataFromRequest(r)
	if err != nil {
		logError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pdf, err := handler.GeneratePDF()
	if err != nil {
		logError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = pdf.Output(w)
	if err != nil {
		logError(err)
	}
}

func openBrowser(url string) {
	var err error

	time.Sleep(1 * time.Second)

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func initLogger(loggingLevel int, logDir string) (err error) {
	var logLevel zerolog.Level

	switch loggingLevel {
	case -1:
		logLevel = zerolog.TraceLevel
	case 0:
		logLevel = zerolog.DebugLevel
	case 1:
		logLevel = zerolog.InfoLevel
	case 2:
		logLevel = zerolog.WarnLevel
	case 3:
		logLevel = zerolog.ErrorLevel
	case 4:
		logLevel = zerolog.FatalLevel
	case 5:
		logLevel = zerolog.PanicLevel
	default:
		logLevel = zerolog.ErrorLevel
	}

	if logDir == "" {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).Level(logLevel).With().Timestamp().Logger()
	} else {
		logName := fmt.Sprintf("%s%s.log", logDir, time.Now().Format("2006-01-02_15-04-05_1111"))
		mainLogFile, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		if err != nil {
			return err
		}

		logger = zerolog.New(mainLogFile).Level(logLevel).With().Timestamp().Logger()
		err = os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func logError(err error) {
	var errStr string
	const printErrStack = true

	if _, ok := err.(*errorsWithStack.Error); ok && printErrStack {
		errStr = err.(*errorsWithStack.Error).ErrorStack()
	} else {
		errStr = err.Error()
	}

	logger.Error().Msgf(errStr)
}
