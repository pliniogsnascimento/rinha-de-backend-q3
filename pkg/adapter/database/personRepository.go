package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/ports/person"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PersonRepo struct {
	gormDb *gorm.DB
	logger *zap.SugaredLogger
}

func NewPersonRepo(gormDb *gorm.DB, logger *zap.SugaredLogger) *PersonRepo {
	return &PersonRepo{gormDb: gormDb, logger: logger}
}

func (db *PersonRepo) closeConnFromPool(conn *pgxpool.Conn) {
	conn.Release()
}

func (db *PersonRepo) FindByID(id string) *person.Person {
	var dto PersonDTO
	if err := db.gormDb.Where("id = ?", id).First(&dto).Error; err != nil {
		db.logger.Errorln(err)
		return nil
	}
	return dto.toEntity()
}

func (db *PersonRepo) FindAll() (*[]person.Person, error) {
	var personDtoList PersonDTOList
	if err := db.gormDb.Find(&personDtoList).Error; err != nil {
		db.logger.Errorln(err)
		return nil, err
	}
	return personDtoList.toEntity(), nil
}

func (db *PersonRepo) Insert(person person.Person) (*person.Person, error) {
	dto := NewPersonDTO(&person)
	if err := db.gormDb.Create(dto).Error; err != nil {
		return nil, err
	}
	return dto.toEntity(), nil
}

func (db *PersonRepo) Count() (int16, error) {
	var count int64
	db.gormDb.Model(&PersonDTO{}).Count(&count)
	return int16(count), nil
}

func (db *PersonRepo) FindByTerm(term string) (*[]person.Person, error) {
	var personDtoList PersonDTOList

	err := db.gormDb.Where("array_to_string(stack, ',') ILIKE '%' || ? || '%' OR name ILIKE '%' || ? || '%' OR nickname ILIKE '%' || ? || '%'",
		term, term, term).
		Find(&personDtoList).Error

	if err != nil {
		db.logger.Errorf("Failed to getting by term: %v\n", err)
		return nil, err
	}
	return personDtoList.toEntity(), nil
}
