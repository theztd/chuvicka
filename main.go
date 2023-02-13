package main

import (
	"fmt"
	"log"
	"os"
	"theztd/chuvicka/agent"
	"theztd/chuvicka/auth"
	"theztd/chuvicka/metrics"
	"theztd/chuvicka/server"

	"theztd/chuvicka/metrics/influx"

	"github.com/joho/godotenv"
)

func loadConf() {
	godotenv.Load(".env")

	influx.Url = os.Getenv("INFLUXDB_URL")
	influx.Token = os.Getenv("INFLUXDB_TOKEN")

	// docasne, jen pro testovani, pak se bude brat z account tabulky
	server.BucketName = "chuvicka"

	auth.DBHost = os.Getenv("AUTH_DB_HOST")
	auth.DBPort = os.Getenv("AUTH_DB_PORT")
	auth.DBUser = os.Getenv("AUTH_DB_USER")
	auth.DBPassword = os.Getenv("AUTH_DB_PASSWORD")
	auth.DBName = os.Getenv("AUTH_DB_NAME")

	auth.JWTHash = os.Getenv("JWT_HASH")

	metrics.DBHost = os.Getenv("AUTH_DB_HOST")
	metrics.DBPort = os.Getenv("AUTH_DB_PORT")
	metrics.DBUser = os.Getenv("AUTH_DB_USER")
	metrics.DBPassword = os.Getenv("AUTH_DB_PASSWORD")
	metrics.DBName = os.Getenv("AUTH_DB_NAME")
}

func main() {
	loadConf()
	if auth.Connect() != nil {
		log.Panicln("Exit!")
	}

	if metrics.Connect() != nil {
		log.Panicln("Exit!")
	}
	metrics.Migrate()
	log.Println("DEBUG: Metrics DB migration...")

	auth.Migrate()
	log.Println("DEBUG: Auth DB migration...")

	if len(os.Args) < 2 {
		fmt.Println("Required argument! (server / agent / check)")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "server":
		server.Run()

	case "list":
		/*
			List all monitored endpoints
		*/
		eps, err := metrics.List()
		if err != nil {
			log.Println("ERR: List metrics error", err)
		}
		for i, ep := range eps {
			fmt.Println(i, " - ", ep)
		}

	case "init-endpoints":
		/*
			Add endpoint to monitoring
		*/
		for _, url := range []string{"https://www.google.com", "https://www.root.cz", "https://troll.fejk.net/v1/sloooowww"} {
			newEp := metrics.Endpoint{}
			newEp.Url = url
			newEp.Add()
		}

	case "useradd":
		newUser := auth.User{}
		newUser.Email = "karel@pokusnak.com"
		newUser.Login = "kaja"
		newUser.HashPassword("tajne")
		newUser.Register()

	case "auth":
		login, password := "kaja", "tajne"
		user, err := auth.Auth(login, password)
		if err != nil {
			log.Println("DEBUG: Non Authorized user")
		} else {
			log.Println("DEBUG: Authorized ", user)
		}

	/*
		Agent gathering metrics
	*/
	case "agent":
		metrics, err := metrics.List()
		if err != nil {
			log.Println("ERR: [agent]", err)
		}

		urls := []string{}
		for _, url := range metrics {
			urls = append(urls, url.Url)
		}

		agent.RunChecks(urls)

	// case "check":
	// 	fmt.Println("Run status check of the configuration and resources")
	// 	data, err := model.GetMetrics("chuvicka", "http://troll.fejk.net/v1/v1/pomala_url")
	// 	if err != nil {
	// 		log.Panicln("ERR:", err)
	// 	}
	// 	for i, d := range data {
	// 		log.Println(i, d)
	// 	}

	default:
		fmt.Println("Read documentation if you needs help")
	}

}
