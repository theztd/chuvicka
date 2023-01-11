package main

import (
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
		log.Println("Measure", url)
		log.Println(httpCheck.Get(url))
		//metrics.Save()
	}

}

func runChecks2() {
	for _, url := range getMonitoredEndpoints() {
		log.Println(httpCheck.GetV2(url))
		//metrics.Save()
	}

}
