package main

import (
	"log"
	"net/http"
	"theztd/chuvicka/model"

	"github.com/gin-gonic/gin"
)

// views/
func index(ctx *gin.Context) {
	// apps, err := model.ListBuckets()
	// if err != nil {
	// 	log.Println("ERR: [ListTables]", err)
	// }

	ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
		// "appList": apps,
	})
}

func metricList(ctx *gin.Context) {
	metricUrls, err := model.ListMeasurements("chuvicka")

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"buckets": []string{},
			"status":  err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"buckets": metricUrls,
		"status":  "Ok",
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
	err := model.WriteMetric("chuvicka", model.Metric{
		Name: "http_endpoint",
		Tags: []model.Tags{
			{Key: "url", Value: input.Url},
			{Key: "StatusCode", Value: "999"},
		},
		Fields: []model.Fields{
			{Key: "ResponseTime", Value: 999},
		},
	})
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

func webUI() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*.tmpl")
	r.GET("/", index)
	r.GET("/api/metrics", metricList)
	r.POST("/api/metrics", metricCreate)
	r.DELETE("/api/metrics/", metricDelete)

	// Admin part
	r.GET("/admin", admin)
	r.GET("/api/tables", bucketList)
	r.POST("/api/tables", bucketCreate)

	r.Run()

}
