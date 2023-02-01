package main

import (
	"fmt"
	"log"
	"os"
	"theztd/chuvicka/model"

	"github.com/joho/godotenv"
)

const bucketName string = "chuvicka"

func main() {
	godotenv.Load(".env")

	model.Url = os.Getenv("INFLUXDB_URL")
	model.Token = os.Getenv("INFLUXDB_TOKEN")

	if len(os.Args) < 2 {
		fmt.Println("Required argument! (server / agent / check)")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		webUI()

	case "agent":
		runChecks()

	case "check":
		fmt.Println("Run status check of the configuration and resources")
		data, err := model.GetMetrics("chuvicka", "http://troll.fejk.net/v1/v1/pomala_url")
		if err != nil {
			log.Panicln("ERR:", err)
		}
		for i, d := range data {
			log.Println(i, d)
		}

	default:
		fmt.Println("Read documentation if you needs help")
	}

}
