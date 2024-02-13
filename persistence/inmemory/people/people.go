package people

import "github.com/insitek/resources/domain"

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type PeopleRepository struct {
	People domain.People `json:"people"`
}

func (pr *PeopleRepository) GetAll(oid string) (domain.People, error) {
	return pr.People, nil
}

func NewPeopleRepository(pr PeopleRepository) (PeopleRepository, error) {
	return pr, nil
}
