package models

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type BasePQDBmodel struct {
	DB *sql.DB
}

func (this *BasePQDBmodel) init() (err error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return
	}
	this.DB = db
	return
}
