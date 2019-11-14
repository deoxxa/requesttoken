package sqlstore

import (
	"database/sql"
	"testing"

	_ "github.com/proullon/ramsql/driver"

	"fknsrs.biz/p/requesttoken/internal/storetest"
)

func getDB(name string) *sql.DB {
	db, err := sql.Open("ramsql", name)
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(DefaultSchema); err != nil {
		panic(err)
	}

	return db
}

func TestCreate(t *testing.T) {
	db := getDB("TestCreate")
	defer db.Close()

	s := Store{DB: db}
	storetest.TestCreate(&s, t)
}

func TestConsume(t *testing.T) {
	db := getDB("TestConsume")
	defer db.Close()

	s := Store{DB: db}
	storetest.TestConsume(&s, t)
}

func TestConsumeExpired(t *testing.T) {
	db := getDB("TestConsumeExpired")
	defer db.Close()

	s := Store{DB: db}
	storetest.TestConsumeExpired(&s, t)
}

func TestConsumeWrongSession(t *testing.T) {
	db := getDB("TestConsumeWrongSession")
	defer db.Close()

	s := Store{DB: db}
	storetest.TestConsumeWrongSession(&s, t)
}

func TestConsumeMissingToken(t *testing.T) {
	db := getDB("TestConsumeMissingToken")
	defer db.Close()

	s := Store{DB: db}
	storetest.TestConsumeMissingToken(&s, t)
}
