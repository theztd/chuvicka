package server

import (
	"log"
	"net/http"
	"theztd/chuvicka/metrics"
	"theztd/chuvicka/metrics/influx"

	"github.com/gin-gonic/gin"
)

// views/
func index(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "ui.tmpl", gin.H{
		// "appList": apps,
	})
}

func metricList(ctx *gin.Context) {
	eps, err := metrics.List()

	graphData := map[string][]influx.MetricResult{}
	for _, ep := range eps {
		graphData[ep.Url], _ = influx.GetMetrics(BucketName, ep.Url)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"buckets":   []string{},
			"status":    err,
			"graphData": graphData,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"buckets":   eps,
		"status":    "Ok",
		"graphData": graphData,
	})
	return
}

func metricCreate(ctx *gin.Context) {
	type newMetric struct {
		Url string `json:"url"`
	}

	input := newMetric{}

	if ctx.ShouldBindJSON(&input) != nil {
		log.Println("ERR: [server] Unable to parse given data.")
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Bad request"})
		return
	}
	log.Println("DEBUG: [server] received json data", input.Url)

	// get retention in day, convert it to seconds
	newEp := metrics.Endpoint{}
	newEp.Url = input.Url
	err := newEp.Add()
	if err != nil {
		log.Println("ERR: [server]", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "Unable to process your request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Metric has been created"})
}

func metricDelete(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "Have to be implemented soon",
		"status": "Ok",
	})
}

// func metricGet(ctx *gin.Context) {
// 	type inData struct {
// 		Url string `json:"url"`
// 	}

// 	input := inData{}

// 	influx.GetMetrics(BucketName, input.Url)
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"msg":    "Have to be implemented soon",
// 		"status": "Ok",
// 	})
// }
