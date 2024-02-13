package presentation

import (
	"encoding/json"

	"github.com/insitek/resources/domain"
	"github.com/nats-io/nats.go/micro"
)

type Response struct {
	People domain.People `json:"people"`
	Error  string        `json:"error"`
}

type Request struct {
	OID string `json:"oid"`
}

type PersonAddRequest struct {
	Person domain.Person       `json:"person"`
	Org    domain.Organisation `json:"org"omitempty`
}

type PeopleHandler struct {
	s domain.PeopleService
}

func (h PeopleHandler) FindAll(req micro.Request) {
	org := Request{}
	data := req.Data()
	err := json.Unmarshal(data, &org)
	if err != nil {
		req.Error("400", "Invalid request", []byte(err.Error()))
		return
	}
	people, err := h.s.GetAll(org.OID)
	if err != nil {
		req.Error("400", "People not found", []byte(err.Error()))
		return
	}
	req.RespondJSON(Response{
		People: people,
	})
}

func (h PeopleHandler) Add(req micro.Request) {
	par := PersonAddRequest{}
	err := json.Unmarshal(req.Data(), &par)
	if err != nil {
		req.Error("400", "Invalid request", []byte(err.Error()))
		return
	}
	if err != nil {
		req.Error("400", "Person not found", []byte(err.Error()))
		return
	}
	err = h.s.Add(par.Person, par.Org)
	if err != nil {
		req.Error("400", "Person not added", []byte(err.Error()))
		return
	}
	req.RespondJSON("{\"success\": \"true\"}")
}

func NewPeopleHandler(s domain.PeopleService) PeopleHandler {
	return PeopleHandler{s: s}
}
