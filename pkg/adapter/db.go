package adapter

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
	"github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/person"
	"go.uber.org/zap"
)

type PersonRepo struct {
	conn   *pgx.Conn
	logger *zap.SugaredLogger
}

func NewDbConn(connStr string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), connStr)
}

func NewPersonRepo(conn *pgx.Conn, logger *zap.SugaredLogger) *PersonRepo {
	return &PersonRepo{conn: conn, logger: logger}
}

func (db *PersonRepo) FindByID(id string) *person.Person {
	person := person.Person{}
	var stack string
	var d time.Time

	err := db.conn.QueryRow(context.Background(), "select * from person where user_id=$1", id).
		Scan(&person.ID, &person.Name, &person.Nickname, &d, &stack)

	if err != nil {
		db.logger.Errorf("Query failed: %v\n", err)
		return nil
	}

	person.BirthDate = d.Format("2006-01-02")
	person.Stack = strings.Split(stack, ",")
	db.logger.Info(person)

	return &person
}

func (db *PersonRepo) FindAll() *[]person.Person {
	persons := []person.Person{}
	rows, err := db.conn.Query(context.Background(), "select user_id, user_name, user_nick, user_birth, user_stack from person")
	if err != nil {
		db.logger.Errorf("Query failed: %v\n", err)
		return nil
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			db.logger.Error(err)
			continue
		}

		person := person.Person{
			ID:        values[0].(string),
			Name:      values[1].(string),
			Nickname:  values[2].(string),
			BirthDate: values[3].(time.Time).Format("2006-01-02"),
			Stack:     strings.Split(values[4].(string), ","),
		}

		db.logger.Infow("Got person from DB",
			"id", person.ID,
			"name", person.Name,
			"nick", person.Nickname,
			"birth", person.BirthDate,
			"stack", person.Stack,
		)
		persons = append(persons, person)
	}

	return &persons
}

func (db *PersonRepo) Insert(person person.Person) (*person.Person, error) {
	person.ID = uuid.NewString()
	if cmd, err := db.conn.Exec(
		context.Background(),
		"INSERT INTO person(user_id, user_name, user_nick, user_birth, user_stack) values ($1, $2, $3, $4, $5)",
		person.ID,
		person.Name,
		person.Nickname,
		person.BirthDate,
		strings.Join(person.Stack, ","),
	); err != nil {
		db.logger.Error(cmd)
		return nil, err
	} else {
		db.logger.Info(cmd)
		return &person, nil
	}
}

func (db *PersonRepo) Count() (int16, error) {
	var count int16
	err := db.conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM person").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (db *PersonRepo) FindByTerm(term string) (*[]person.Person, error) {
	persons := []person.Person{}

	rows, err := db.conn.Query(context.Background(), "SELECT * FROM person WHERE user_stack ILIKE '%' || $1 || '%' OR user_name ILIKE '%' || $1 || '%' OR user_nick ILIKE '%' || $1 || '%'", term)
	if err != nil {
		db.logger.Errorf("Query failed: %v\n", err)
		return nil, err
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			db.logger.Error(err)
			continue
		}

		person := person.Person{
			ID:        values[0].(string),
			Name:      values[1].(string),
			Nickname:  values[2].(string),
			BirthDate: values[3].(time.Time).Format("2006-01-02"),
			Stack:     strings.Split(values[4].(string), ","),
		}

		db.logger.Infow("Got person from DB",
			"id", person.ID,
			"name", person.Name,
			"nick", person.Nickname,
			"birth", person.BirthDate,
			"stack", person.Stack,
		)
		persons = append(persons, person)
	}

	return &persons, nil
}

func (db *PersonRepo) Close() {
	db.conn.Close(context.Background())
}
