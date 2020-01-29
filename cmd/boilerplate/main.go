package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gemsorg/boilerplate/pkg/authorization"
	"github.com/gemsorg/boilerplate/pkg/database"
	"github.com/gemsorg/boilerplate/pkg/datastore"
	"github.com/gemsorg/boilerplate/pkg/service"
	"github.com/joho/godotenv"

	"github.com/gemsorg/boilerplate/pkg/server"
)

func main() {
	environment := flag.String("env", "local", "use compose in compose-dev")
	flag.Parse()

	if *environment != "compose" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Connect to db
	db, err := database.Connect()
	if err != nil {
		log.Fatal("mysql connection error", err)
	}
	defer db.Close()
	ds := datastore.NewDatastore(db)
	authorizer := authorization.NewAuthorizer()
	svc := service.New(ds, authorizer)
	s := server.New(db, svc)
	log.Println("info", fmt.Sprintf("Starting service on port 3000"))
	http.Handle("/", s)
	http.ListenAndServe(":3000", nil)
}
