/*
	Vytunena verze by chatGPT

	Je pomalejsi, ale uspornejsi na zdroje. Zpomaleni je v radu jednotek procent, setreni zdroju cca o 20%

	Uspory je dosazeno poolovanim requestu, uvidime co to udela pri stovkach pozadavku
*/
package httpCheck

import (
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"time"
)

func getClient() *http.Client {
	clientOnce.Do(func() {
		client = &http.Client{}
	})
	return client
}

func GetV2(url string) (Response, error) {
	var err error

	client := getClient()

	ret := Response{
		Url: url,
	}

	start := time.Now()
	req, _ := http.NewRequest("GET", url, nil)
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			ret.TCPConnection = int(time.Since(start))
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			ret.DNSLookup = int(time.Since(start))
		},
		GotFirstResponseByte: func() {
			ret.TTFB = int(time.Since(start))
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			if err == nil {
				ret.TLSHandshake = int(time.Since(start))
			}
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := client.Do(req)
	if err == nil {
		ret.StatusCode = uint(resp.StatusCode)
	}

	ret.ResponseTime = int(time.Since(start))

	return ret, err
}
