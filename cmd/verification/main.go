package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/expandorg/verification/pkg/authorization"
	"github.com/expandorg/verification/pkg/automatic"
	"github.com/expandorg/verification/pkg/database"
	"github.com/expandorg/verification/pkg/datastore"
	"github.com/expandorg/verification/pkg/externalsvc"
	"github.com/expandorg/verification/pkg/registrysvc"
	"github.com/expandorg/verification/pkg/responsesvc"
	"github.com/expandorg/verification/pkg/service"
	"github.com/joho/godotenv"

	"github.com/expandorg/verification/pkg/server"
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
	responseSVC := responsesvc.New(authToken)
	consensus := automatic.NewConsensus(ds, responseSVC)

	svc := service.New(ds, authorizer, rsvc, external, consensus)
	s := server.New(svc)
	log.Println("info", fmt.Sprintf("Starting service on port 8186"))
	http.Handle("/", s)
	http.ListenAndServe(":8186", nil)
}
