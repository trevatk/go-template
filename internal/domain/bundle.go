// Package domain application layer of service
package domain

// Bundle service bundle
type Bundle struct {
	PersonService *PersonService
}

// NewBundle create new service bundle
func NewBundle(personService *PersonService) *Bundle {
	return &Bundle{
		PersonService: personService,
	}
}
