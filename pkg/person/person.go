package person

// TODO: Fix date format
type Person struct {
	ID        string   `json:"id"`
	Name      string   `json:"nome" binding:"required"`
	Nickname  string   `json:"apelido" binding:"required"`
	BirthDate string   `json:"nascimento" binding:"required"`
	Stack     []string `json:"stack"`
}

type PersonService interface {
	FindAll() *[]Person
	FindByTerm(term string) (*[]Person, error)
	Insert(person Person) (*Person, error)
	FindByID(id string) *Person
	Count() (int16, error)
}
