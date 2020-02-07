package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gemsorg/verification/pkg/authorization"
	"github.com/gemsorg/verification/pkg/database"
	"github.com/gemsorg/verification/pkg/datastore"
	"github.com/gemsorg/verification/pkg/externalsvc"
	"github.com/gemsorg/verification/pkg/registrysvc"
	"github.com/gemsorg/verification/pkg/service"
	"github.com/joho/godotenv"

	"github.com/gemsorg/verification/pkg/server"
)

func main() {
	environment := flag.String("env", "local", "use compose in compose-dev")
	flag.Parse()

	if *environment == "local" {
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
	authToken := authorizer.GetAuthToken()
	rsvc := registrysvc.New(authToken)
	external := externalsvc.New(authToken)
	svc := service.New(ds, authorizer, rsvc, external)
	s := server.New(svc)
	log.Println("info", fmt.Sprintf("Starting service on port 8186"))
	http.Handle("/", s)
	http.ListenAndServe(":8186", nil)
}
