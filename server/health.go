package server

import (
	"log"
	"net/http"
	"os"
	"theztd/chuvicka/auth"
	"theztd/chuvicka/metrics"
	"time"

	"github.com/gin-gonic/gin"
)

type Check struct {
	Status        string      `json:"status"`
	Output        string      `json:"output"`
	ComponentId   string      `json:"componentId"`
	Time          string      `json:"time"`
	ComponentType string      `json:"componentType"`
	ObservedValue interface{} `json:"observedValue"`
	ObservedUnit  string      `json:"observedUnit"`
}

type Healthz struct {
	Status      string           `json:"status"`
	Version     string           `json:"version"`
	InstanceId  string           `json:"instanceId"`
	ReleaseId   string           `json:"releaseId"`
	Output      string           `json:"output"`
	Checks      map[string]Check `json:"Checks"`
	Description string           `json:"description"`
	Notes       []string         `json:"notes"`
}

func healthStatus(ctx *gin.Context) {
	// Check postgresDB
	var pgH, influxH Check
	users := []auth.UserRead{}
	dbErr := metrics.DB.Model(&auth.User{}).Find(&users).Limit(1).Error
	if dbErr == nil {
		pgH.Status = "pass"
	} else {
		log.Println("Error [HEALTHZ] - ", dbErr.Error())
		pgH.Status = "fail"
		pgH.Output = dbErr.Error()
		pgH.ComponentType = "db"
		pgH.ComponentId = auth.DBHost + "/" + auth.DBName
		pgH.Time = time.Now().String()
	}

	// Check InfluxDB
	_, influxErr := metrics.List()
	log.Println("DEBUG: ", influxErr)
	if influxErr == nil {
		influxH.Status = "pass"
	} else {
		log.Println("Error [HEALTHZ] - ", influxErr.Error())
		influxH.Status = "fail"
		influxH.Output = influxErr.Error()
		influxH.ComponentType = "db"
		pgH.Time = time.Now().String()
	}

	hostname, _ := os.Hostname()
	healthRet := Healthz{
		Status:     "pass",
		InstanceId: hostname,
		Version:    VERSION,
		Checks: map[string]Check{
			"postgres": pgH,
			"influx":   influxH,
		},
	}

	// ctx.Request.Header.Add("Content-Type", "application/health+json")
	ctx.JSON(http.StatusOK, healthRet)
}
