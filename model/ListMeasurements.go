package model

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func WriteMetric(bucket string, metric Metric) error {
	client := influxdb2.NewClient(Url, Token)
	writeApi := client.WriteAPIBlocking("myorg", bucket)

	p := influxdb2.NewPointWithMeasurement(metric.Name)

	for _, tag := range metric.Tags {
		p.AddTag(tag.Key, tag.Value)
	}

	for _, item := range metric.Fields {
		p.AddField(item.Key, item.Value)
	}

	p.SetTime(time.Now())

	err := writeApi.WritePoint(context.Background(), p)
	if err != nil {
		return err
	}

	return nil
}

func GetMetrics(bucketName string, url string) ([]MetricResult, error) {
	log.Println("ListMetrics: ")
	var ret []MetricResult

	client := influxdb2.NewClient(Url, Token)
	queryAPI := client.QueryAPI("myorg")
	query := fmt.Sprintf(`
	from(bucket: "%s")
		|> range(start: -1h)
		|> filter(fn: (r) => r["_measurement"] == "http_endpoint")
		|> filter(fn: (r) => r["_field"] == "ResponseTime")
		|> filter(fn: (r) => r["url"] == "%s")
		|> yield(name: "mean")
	`, bucketName, url)

	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	for results.Next() {
		out := MetricResult{}

		x := results.Record().String()
		for _, i := range strings.Split(x, ",") {
			iList := strings.Split(i, ":")
			switch iList[0] {
			case "_field":
				out.Name = iList[1]
			case "StatusCode":
				out.StatusCode, _ = strconv.Atoi(string(iList[1]))
			case "url":
				out.Url = strings.Join(iList[1:], ":")
			case "result":
				out.Result = iList[1]
			case "_time":
				out.Time = strings.Join(iList[1:], ":")
			case "_value":
				out.Value, _ = strconv.ParseFloat(string(iList[1]), 32)
			}
		}
		ret = append(ret, out)
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}

	return ret, nil
}

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
