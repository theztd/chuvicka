/*
Muj kod, ktery je rychly a plne funkcni
*/
package httpCheck

import (
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"time"
)

func Request(e Endpoint) (ret Response, err error) {
	// returning time in miliseconds
	err = nil
	ret = Response{
		Url: e.Url,
	}
	start := time.Now()

	req, _ := http.NewRequest(e.Method, e.Url, nil)
	req.Header.Add("User-Agent", "agent-chuvicka")
	if len(e.Headers) > 0 {
		for k, v := range e.Headers {
			req.Header.Add(k, v)
		}
	}

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

	//defer req.Body.Close()
	// ret.StatusCode = uint(req.Response.StatusCode)
	// d, err := io.ReadAll(req.Body)
	// if err == nil {
	// 	ret.Result = string(d)
	// }
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err == nil {
		ret.StatusCode = uint(resp.StatusCode)
	}

	ret.ResponseTime = int(time.Since(start))

	return ret, err
}
