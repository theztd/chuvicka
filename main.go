package main

import (
	"fmt"
	"os"
	"theztd/chuvicka/model"

	"github.com/joho/godotenv"
)

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

	case "agent2":
		runChecks2()

	case "check":
		fmt.Println("Run status check of the configuration and resources")

	default:
		fmt.Println("Read documentation if you needs help")
	}

}
