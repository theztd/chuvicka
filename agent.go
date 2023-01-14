package main

import (
	"fmt"
	"log"
	"theztd/chuvicka/httpCheck"
	"theztd/chuvicka/model"
)

func getMonitoredEndpoints() []string {
	urls, err := model.ListMeasurements("chuvicka")
	if err != nil {
		log.Println("ERR: [agent]", err)
	}
	return urls
}

func runChecks() {
	for _, url := range getMonitoredEndpoints() {
		log.Println("INFO: [Agent] measure", url)

		// Run check
		result, _ := httpCheck.Get(url)

		// Report measured metrics
		err := model.WriteMetric("chuvicka", model.Metric{
			Name: "http_endpoint",
			Tags: []model.Tags{
				{Key: "url", Value: url},
				{Key: "StatusCode", Value: fmt.Sprintf("%d", result.StatusCode)},
			},
			Fields: []model.Fields{
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
