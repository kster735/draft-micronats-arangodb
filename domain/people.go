package domain

type Person struct {
	PID       string `json:"pid"omitempty`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type People []Person

type PeopleService interface {
	GetAll(oid string) (People, error)
	Add(person Person, org Organisation) error
}
