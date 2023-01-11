package model

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func CreateBucket(organization string, bucketName string, expireInSec uint) error {
	postData := []byte(fmt.Sprintf(`{
		"orgID": "'"%s"'",
		"name": "%s",
		"retentionRules": [
		  {
			"type": "expire",
			"everySeconds": %d,
			"shardGroupDurationSeconds": 0
		  }
		]
	  }`, organization, bucketName, expireInSec))

	req, err := http.NewRequest("POST", Url+"/api/v2/buckets", bytes.NewBuffer(postData))
	if err != nil {
		log.Println("ERR: [influx]", err)
	}
	req.Header.Add("Authorization", "Token "+Token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERR: [influx]", err)
	}
	defer resp.Body.Close()
	//var ret Buckets
	log.Println(resp.StatusCode, resp.Status, resp.Body)
	//json.NewDecoder(resp.Body).Decode(&ret)

	//log.Println("Parsovana odpoved:", ret)

	return nil
}
