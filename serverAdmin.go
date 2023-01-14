package main

import (
	"log"
	"net/http"
	"theztd/chuvicka/model"

	"github.com/gin-gonic/gin"
)

// views/
func admin(ctx *gin.Context) {
	apps, err := model.ListBuckets()
	if err != nil {
		log.Println("ERR: [ListTables]", err)
	}

	ctx.HTML(http.StatusOK, "admin.tmpl", gin.H{
		"appList": apps,
	})
}

func bucketList(ctx *gin.Context) {

	log.Println("---------------------------")
	d, _ := model.ListMeasurements("chuvicka")
	log.Println("DEBUG: [ListMeasurements] result:", d)
	log.Println("---------------------------")

	buckets, err := model.ListBuckets()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"buckets": []string{},
			"status":  err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"buckets": buckets,
		"status":  "Ok",
	})
	return
}

func bucketCreate(ctx *gin.Context) {
	type newBucket struct {
		Name          string `json:"tableName"`
		RetentionDays uint   `json:"retentionDays"`
	}

	input := newBucket{}

	if ctx.ShouldBindJSON(&input) != nil {
		log.Println("ERR: [server] Unable to parse given data.")
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Bad request"})
		return
	}
	log.Println("DEBUG: [server] received json data", input)

	// get retention in day, convert it to seconds
	err := model.CreateBucket("marek", input.Name, 60*60*24*input.RetentionDays)
	if err != nil {
		log.Println("ERR: [server]", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "Unable to process your request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Table has been created"})
}
