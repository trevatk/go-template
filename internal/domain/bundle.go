package domain

// Bundle service bundle
type Bundle struct {
	PersonService *PersonService
}

// NewBundle
func NewBundle(personService *PersonService) *Bundle {
	return &Bundle{
		PersonService: personService,
	}
}
