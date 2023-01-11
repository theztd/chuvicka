package model

import (
	"encoding/json"
	"log"
	"net/http"
)

type Buckets struct {
	Buckets []Bucket `json:"buckets"`
}

type Retention struct {
	Interval uint   `json:"everySeconds"`
	Type     string `json:"type"`
}

type Bucket struct {
	Id          uint        `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Labels      []string    `json:"labels"`
	Retentions  []Retention `json:"retentionRules"`
}

func ListBuckets() ([]Bucket, error) {
	req, err := http.NewRequest("GET", Url+"/api/v2/buckets", nil)
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
	var ret Buckets
	log.Println(resp.StatusCode, resp.Status, resp.Body)
	json.NewDecoder(resp.Body).Decode(&ret)

	log.Println("Parsovana odpoved:", ret)

	return ret.Buckets, nil
}
