package agent

import (
	"fmt"
	"log"
	"theztd/chuvicka/agent/httpCheck"
	"theztd/chuvicka/metrics/influx"
)

func RunChecks(endpoints []string) {
	for _, url := range endpoints {
		log.Println("INFO: [Agent] measure", url)

		// Run httpCheck
		result, _ := httpCheck.Get(url)

		// Report measured metrics
		err := influx.WriteMetric("chuvicka", influx.Metric{
			Name: "http_endpoint",
			Tags: []influx.Tags{
				{Key: "url", Value: url},
				{Key: "StatusCode", Value: fmt.Sprintf("%d", result.StatusCode)},
			},
			Fields: []influx.Fields{
				{Key: "TTFB", Value: float32(result.TTFB)},
				{Key: "TCPConnection", Value: float32(result.TCPConnection)},
				{Key: "TLSHandshake", Value: float32(result.TLSHandshake)},
				{Key: "DNSLookup", Value: float32(result.DNSLookup)},
				{Key: "ResponseTime", Value: float32(result.ResponseTime)},
			},
		})
		if err != nil {
			log.Println("ERROR: [Agent]", url, err)
		}
	}

}
