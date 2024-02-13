package people

import (
	"context"
	"log"

	"github.com/arangodb/go-driver"
	"github.com/insitek/resources/domain"
)

type PeopleRepository struct {
	d driver.Database
}

func (pr PeopleRepository) Add(person domain.Person, org domain.Organisation) error {
	pcollection, err := pr.d.Collection(context.Background(), "people")
	if err != nil {
		pcollection, err = pr.d.CreateCollection(context.Background(), "people", nil)
		if err != nil {
			log.Printf("could not find or create collection people, error: %v", err)
			return err
		}
	}

	pcollection.CreateDocument(context.Background(), person)

	return nil
}

func (pr PeopleRepository) Find(pid string) (domain.Person, error) {
	return domain.Person{}, nil
}

func (pr PeopleRepository) GrantCredentials(person domain.Person, credentials domain.Credentials) error {
	return nil
}

func (pr PeopleRepository) GetAll(oid string) (domain.People, error) {
	aql := `
	WITH people, organisations
	let orgid = CONCAT( 'organisations/', @oid)
		FOR v, e, p IN 1..1 INBOUND orgid person_in_organisation
  			RETURN { pid: p.vertices[1]._key, first_name: p.vertices[1].first_name, last_name: p.vertices[1].last_name}
	`
	cursor, err := pr.d.Query(context.Background(), aql, map[string]interface{}{"oid": oid})
	if err != nil {
		log.Printf("could not execute query, error: %v", err)
		return domain.People{}, err
	}
	defer cursor.Close()
	people := domain.People{}
	for {
		if !cursor.HasMore() {
			break
		}
		p := domain.Person{}
		_, err := cursor.ReadDocument(context.Background(), &p)
		if err != nil {
			return domain.People{}, err
		}
		people = append(people, domain.Person{
			PID:       p.PID,
			FirstName: p.FirstName,
			LastName:  p.LastName,
		})
	}
	return people, nil
}

func NewPeopleRepository(d driver.Database) (PeopleRepository, error) {
	return PeopleRepository{d}, nil
}
