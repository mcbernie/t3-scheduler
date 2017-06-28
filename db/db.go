package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // driver for mysql usage
)

//Database simple Struct for database operations
type Database struct {
	Configuration string
}

//Create new database operations struct
func Create(config string) Database {
	return Database{Configuration: config}
}

//Query function type
type Query func(db *sql.DB) interface{}

//DatabaseOperation Method for make operations on database
func (d *Database) DatabaseOperation(fn Query) interface{} {
	db, err := sql.Open("mysql", d.Configuration)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

	ret := fn(db)
	//fmt.Println("return:", ret)
	return ret
}
