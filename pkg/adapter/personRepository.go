package adapter

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/person"
	"go.uber.org/zap"
)

type PersonRepo struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewPersonRepo(pool *pgxpool.Pool, logger *zap.SugaredLogger) *PersonRepo {
	return &PersonRepo{pool: pool, logger: logger}
}

func (db *PersonRepo) closeConnFromPool(conn *pgxpool.Conn) {
	conn.Release()
	// if err != nil {
	// 	db.logger.Errorf("Failed closing connection from pool: %v\n", err)
	// }
}

func (db *PersonRepo) FindByID(id string) *person.Person {
	conn, err := db.pool.Acquire(context.TODO())
	if err != nil {
		db.logger.Errorf("Failed to get connection from pool: %v\n", err)
		return nil
	}
	defer db.closeConnFromPool(conn)

	person := person.Person{}
	var stack string
	var d time.Time

	err = conn.QueryRow(context.Background(), "select * from person where user_id=$1", id).
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

func (db *PersonRepo) FindAll() (*[]person.Person, error) {
	conn, err := db.pool.Acquire(context.TODO())
	if err != nil {
		db.logger.Errorf("Failed to get connection from pool: %v\n", err)
		return nil, err
	}
	defer db.closeConnFromPool(conn)

	persons := []person.Person{}
	rows, err := conn.Query(context.Background(), "select user_id, user_name, user_nick, user_birth, user_stack from person")
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

func (db *PersonRepo) Insert(person person.Person) (*person.Person, error) {
	conn, err := db.pool.Acquire(context.TODO())
	if err != nil {
		db.logger.Errorf("Failed to get connection from pool: %v\n", err)
		return nil, err
	}
	defer db.closeConnFromPool(conn)

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		return nil, err
	}

	person.ID = uuid.NewString()
	if cmd, err := tx.Exec(
		context.Background(),
		"INSERT INTO person(user_id, user_name, user_nick, user_birth, user_stack) values ($1, $2, $3, $4, $5)",
		person.ID,
		person.Name,
		person.Nickname,
		person.BirthDate,
		strings.Join(person.Stack, ","),
	); err != nil {
		tx.Rollback(context.TODO())
		db.logger.Error(err)
		return nil, err
	} else {
		tx.Commit(context.TODO())
		db.logger.Info(cmd)
		return &person, nil
	}
}

func (db *PersonRepo) Count() (int16, error) {
	conn, err := db.pool.Acquire(context.TODO())
	if err != nil {
		db.logger.Errorf("Failed to get connection from pool: %v\n", err)
		return 0, err
	}
	defer db.closeConnFromPool(conn)

	var count int16
	err = conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM person").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (db *PersonRepo) FindByTerm(term string) (*[]person.Person, error) {
	conn, err := db.pool.Acquire(context.TODO())
	if err != nil {
		db.logger.Errorf("Failed to get connection from pool: %v\n", err)
		return nil, err
	}
	defer db.closeConnFromPool(conn)

	persons := []person.Person{}

	rows, err := conn.Query(context.Background(), "SELECT * FROM person WHERE user_stack ILIKE '%' || $1 || '%' OR user_name ILIKE '%' || $1 || '%' OR user_nick ILIKE '%' || $1 || '%'", term)
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

		db.logger.Debugw("Got person from DB",
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
	db.pool.Close()
}
