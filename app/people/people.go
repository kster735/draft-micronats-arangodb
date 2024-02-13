package people

import (
	"github.com/insitek/resources/domain"
	apeople "github.com/insitek/resources/persistence/arangodb/people"
	"github.com/insitek/resources/persistence/inmemory/people"
)

type Service struct {
	repo  people.PeopleRepository
	arepo apeople.PeopleRepository
}

// GetAll implements domain.PeopleService.
func (s Service) GetAll(oid string) (domain.People, error) {
	return s.arepo.GetAll(oid)
}

func (s Service) Add(person domain.Person, org domain.Organisation) error {
	return s.arepo.Add(person, org)
}

func NewService(repo people.PeopleRepository, arepo apeople.PeopleRepository) Service {
	return Service{repo: repo, arepo: arepo}
}
