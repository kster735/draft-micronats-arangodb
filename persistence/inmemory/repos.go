package inmemory

import (
	"github.com/insitek/resources/domain"
)

type PeopleRepository interface {
	GetAll(oid string) (domain.People, error)
}
