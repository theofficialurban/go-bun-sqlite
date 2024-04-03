package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type Database struct {
	SQL *sql.DB
	Bun *bun.DB
	DBF *dbfixture.Fixture
}

var db *Database = &Database{}

type User struct {
	ID    int64 `bun:",pk,autoincrement"`
	Name  string
	Email string
	Age   int
}

func CTNE[T any]() {
	d := db.Bun

	d.RegisterModel((*T)(nil))
	_, err := d.NewCreateTable().Model((*T)(nil)).IfNotExists().
		Exec(context.Background())
	if err != nil {
		panic(err)
	}

	var tp T
	t := reflect.TypeOf(tp)

	fmt.Printf("Created Table %s\n", t.Name())
}
func Insert[T any](data *T) {
	d := db.Bun
	_, err := d.NewInsert().Model(data).Exec(context.Background())
	if err != nil {
		panic(err)
	}
}
func GetAll[T any](relation string) *([]T) {
	d := db.Bun
	items := new([]T)
	if relation != "" {
		err := d.NewSelect().Model(items).Relation(relation).Scan(context.Background())
		if err != nil {
			panic(err)
		}
		return items
	} else {
		err := d.NewSelect().Model(items).Scan(context.Background())
		if err != nil {
			panic(err)
		}
		return items
	}
}
func init() {
	// Creates a SQLite database in `sql.sqlite` file
	sqldb, err := sql.Open(sqliteshim.DriverName(), "file:sql.sqlite?cache=shared")
	if err != nil {
		panic(err)
	}
	db.SQL = sqldb
	db.Bun = bun.NewDB(db.SQL, sqlitedialect.New())
	db.DBF = dbfixture.New(db.Bun)
	fmt.Println("Database Setup Complete!")

	CTNE[User]()
	db.LoadFixtures()

}

func (db *Database) LoadFixtures() {
	err := db.DBF.Load(context.Background(), os.DirFS("."), "fixture.yml")
	if err != nil {
		panic(err)
	}
}

func main() {
	if err := db.Bun.Ping(); err != nil {
		panic(err)
	}
}
