package bigip

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type httpLogger struct {
	log *log.Logger
}

func newLogger() *httpLogger {
	return &httpLogger{
		log: log.New(os.Stdout, "log - ", log.LstdFlags),
	}
}

func (l *httpLogger) LogRequest(req *http.Request) {
	body, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}

	l.log.Printf("HTTP Request method=%s url=%s, body=%s", req.Method, req.URL.String(), body)
}

func (l *httpLogger) LogResponse(req *http.Request, res *http.Response, httperr error, duration time.Duration) {
	duration /= time.Millisecond
	defer res.Body.Close()

	// Dump response body
	body, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		l.log.Printf("HTTP Response status=error durationMs=%d error=%q", duration, httperr.Error())
	} else {
		l.log.Printf("HTTP Response status=%d durationMs=%d, body=%s", res.StatusCode, duration, body)
	}
}
