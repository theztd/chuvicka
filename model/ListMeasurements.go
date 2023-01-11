package model

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func ListMeasurements(bucketName string) ([]string, error) {
	var ret []string

	client := influxdb2.NewClient(Url, Token)
	queryAPI := client.QueryAPI("myorg")
	query := `import "influxdata/influxdb/schema"

	schema.measurementTagValues(
		bucket: "chuvicka",
		measurement: "http_endpoint",
		tag: "url"
	)`
	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	for results.Next() {
		x := strings.ReplaceAll(fmt.Sprintf("%s", results.Record().Value()), "\"", "")
		ret = append(ret, x)
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}

	return ret, nil
}

func ListMeasurementsAPI(bucketName string) ([]string, error) {
	data := url.Values{}

	/*
		FLUX format

		import "influxdata/influxdb/schema"

		schema.measurementTagValues(
		bucket: "chuvicka",
		measurement: "http_endpoint",
		tag: "url"
		)
	*/
	// FLUX format
	fluxQ := `import "influxdata/influxdb/schema"
		schema.measurementTagValues(
		bucket: "%s",
		measurement: "http_endpoint",
		tag: "url")`

	data.Set("q", fmt.Sprintf(fluxQ, bucketName))

	// Old format
	//data.Set("q", fmt.Sprintf("SELECT url FROM %s.http_endpoint", bucketName))
	req, err := http.NewRequest("POST", Url+"/api/v2/query?orgID=c89dcc5170f3030a", strings.NewReader(data.Encode()))
	if err != nil {
		log.Println("ERR: [influx]", err)
	}
	req.Header.Add("Authorization", "Token "+Token)
	req.Header.Add("Content-Type", "application/vnd.flux")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERR: [influx]", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERR:", err)
	}
	log.Println(resp.StatusCode, resp.Status, string(body))

	return []string{}, nil
}
