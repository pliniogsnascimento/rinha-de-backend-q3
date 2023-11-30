package database

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/ports/person"
	"gorm.io/gorm"
)

type PersonDTOList []PersonDTO

type PersonDTO struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string
	Nickname  string
	BirthDate string
	Stack     pq.StringArray `gorm:"type:text[]"`
}

func (dto *PersonDTO) toEntity() *person.Person {
	return &person.Person{
		ID:        dto.ID.String(),
		Name:      dto.Name,
		Nickname:  dto.Nickname,
		BirthDate: dto.BirthDate,
		Stack:     dto.Stack,
	}
}

func NewPersonDTO(entity *person.Person) *PersonDTO {
	dto := &PersonDTO{
		Name:      entity.Name,
		Nickname:  entity.Nickname,
		BirthDate: entity.BirthDate,
		Stack:     entity.Stack,
	}

	if entity.ID != "" {
		dto.ID = uuid.MustParse(entity.ID)
	}
	return dto
}

func (list *PersonDTOList) toEntity() *[]person.Person {
	personList := []person.Person{}
	for _, dto := range *list {
		personList = append(personList, *dto.toEntity())
	}

	return &personList
}
