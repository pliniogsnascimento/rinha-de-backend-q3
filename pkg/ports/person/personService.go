package person

type PersonService interface {
	FindAll() (*[]Person, error)
	FindByTerm(term string) (*[]Person, error)
	Insert(person Person) (*Person, error)
	FindByID(id string) *Person
	Count() (int16, error)
}
