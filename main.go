package main

import (
	"fmt"
	"log"
	"os"
	"theztd/chuvicka/agent"
	"theztd/chuvicka/auth"
	authModel "theztd/chuvicka/auth/model"
	"theztd/chuvicka/metrics"
	"theztd/chuvicka/server"

	"theztd/chuvicka/metrics/influx"

	"github.com/joho/godotenv"
)

func init() {
	env := os.Getenv("ENV")
	if len(env) < 1 {
		env = "devel"
	}
	log.Println("Running env", env)

	godotenv.Load(".env-" + env)

	influx.Url = os.Getenv("INFLUXDB_URL")
	influx.Token = os.Getenv("INFLUXDB_TOKEN")

	// docasne, jen pro testovani, pak se bude brat z account tabulky
	server.BucketName = "chuvicka"

	auth.DBHost = os.Getenv("AUTH_DB_HOST")
	auth.DBPort = os.Getenv("AUTH_DB_PORT")
	auth.DBUser = os.Getenv("AUTH_DB_USER")
	auth.DBPassword = os.Getenv("AUTH_DB_PASSWORD")
	auth.DBName = os.Getenv("AUTH_DB_NAME")

	auth.JWTHashSet(os.Getenv("JWT_HASH"))

	metrics.DBHost = os.Getenv("AUTH_DB_HOST")
	metrics.DBPort = os.Getenv("AUTH_DB_PORT")
	metrics.DBUser = os.Getenv("AUTH_DB_USER")
	metrics.DBPassword = os.Getenv("AUTH_DB_PASSWORD")
	metrics.DBName = os.Getenv("AUTH_DB_NAME")

	server.UI = os.Getenv("UI")
}

func main() {
	//loadConf()
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
		for _, url := range []string{"https://www.google.com", "http://localhost:8080/_healthz/ready.json", "https://troll.fejk.net/v1/sloooowww"} {
			newEp := metrics.Endpoint{}
			newEp.Url = url
			newEp.Request.Code = 200
			newEp.Request.Match = "*"
			newEp.Response.Code = 399
			newEp.Response.Match = "blahblah"

			log.Println(newEp)
			newEp.Add()
		}

	case "init-users":
		token1 := authModel.Token{Name: "org1 - karluv"}
		token1.Hash("supertajny-token")

		newUser1 := authModel.User{}
		newUser1.Email = "kaja@pokusnak.com"
		newUser1.Login = "karel"
		newUser1.HashPassword("tajne")
		//newUser1.Register()

		newUser2 := authModel.User{}
		newUser2.Email = "pepa@testik.com"
		newUser2.Login = "pepa"
		newUser2.HashPassword("heslo")
		//newUser2.Register()

		org1 := authModel.Organization{}
		org1.Name = "Organization1"
		org1.Description = "Malinka organizace z podhuri..."
		org1.Tokens = append(org1.Tokens, token1)
		org1.Users = append(org1.Users, newUser1, newUser2)
		org1.Save()

	/*
		Testing users
	*/
	case "test-users":
		// validate

		// var users []auth.User
		org, err := auth.GetOrg("Organization1")
		if err != nil {
			log.Println("ERR: Unable to find organization with the given name", err)
		}
		org.Pretty()

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
