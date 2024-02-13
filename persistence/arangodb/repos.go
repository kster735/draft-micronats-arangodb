package arangodb

import (
	"github.com/insitek/resources/domain"
)

type PeopleRepository interface {
	Add(person domain.Person, org domain.Organisation) error
	Find(pid string) (domain.Person, error)
	GrantCredentials(person domain.Person, credentials domain.Credentials) error
	GetAll(oid string) (domain.People, error)
}
