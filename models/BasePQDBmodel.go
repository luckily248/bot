package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type BasePQDBmodel struct {
	DB *sql.DB
}

func (this *BasePQDBmodel) init() (err error) {
	db, err := sql.Open("postgres", "user=postgres password=81099371 dbname=postgres sslmode=disable")
	if err != nil {
		return
	}
	this.DB = db
	return
}
