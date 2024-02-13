package main

import (
	"fmt"
	"log"
	"os"
	"time"

	driver "github.com/arangodb/go-driver"
	arangogttp "github.com/arangodb/go-driver/http"
	ps "github.com/insitek/resources/app/people"
	"github.com/insitek/resources/domain"
	apeople "github.com/insitek/resources/persistence/arangodb/people"
	"github.com/insitek/resources/persistence/inmemory/people"
	"github.com/insitek/resources/presentation"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

func main() {
	fmt.Println("Starting Resources microservice...")
	arangoUrl := os.Getenv("ARANGO_URL")
	arangoUser := os.Getenv("ARANGO_USER")
	arangoPass := os.Getenv("ARANGO_PASS")
	arangoDatabase := os.Getenv("ARANGO_DB")
	url := os.Getenv("NATS_URL")

	nc, _ := nats.Connect(url)
	defer nc.Close()
	// Establish connection with arangodb
	conn, err := arangogttp.NewConnection(arangogttp.ConnectionConfig{
		Endpoints: []string{arangoUrl},
	})
	if err != nil {
		log.Fatal("cannot connect to arangodb")
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(arangoUser, arangoPass),
	})
	if err != nil {
		log.Fatal("cannot login to arangodb")
	}

	database, err := client.Database(nil, arangoDatabase)
	if err != nil {
		log.Fatal("cannot get database %s, error: %v", arangoDatabase, err)
	}

	arepo, err := apeople.NewPeopleRepository(database)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new dummy PeopleRepository
	pr := people.PeopleRepository{}
	pr.People = domain.People{
		{
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			FirstName: "Jane",
			LastName:  "Doe",
		},
	}

	prepo, err := people.NewPeopleRepository(pr)
	if err != nil {
		log.Fatal(err)
	}
	pservice := ps.NewService(prepo, arepo)
	phandler := presentation.NewPeopleHandler(pservice)

	config := micro.Config{
		Name:        "Resources",
		Version:     "0.0.1",
		Description: "A microservice for managing resources",
	}

	microSvc, err := micro.AddService(nc, config)
	if err != nil {
		log.Fatal(err)
	}
	defer microSvc.Stop()

	resourcesGrp := microSvc.AddGroup("resources")

	err = resourcesGrp.AddEndpoint("people-getall", micro.HandlerFunc(phandler.FindAll), micro.WithEndpointSubject("people.getall"))
	if err != nil {
		log.Fatal(err)
	}

	err = resourcesGrp.AddEndpoint("AddPersonToPeople", micro.HandlerFunc(phandler.Add), micro.WithEndpointSubject("people.add"))
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(time.Minute)
	}
}
