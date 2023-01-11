package httpCheck

import (
	"net/http"
	"sync"
)

// vytunena verze by ChatGPT
var client *http.Client
var clientOnce sync.Once

type Response struct {
	Url           string
	StatusCode    uint
	TCPConnection int
	DNSLookup     int
	TTFB          int
	TLSHandshake  int
	ResponseTime  int
}
