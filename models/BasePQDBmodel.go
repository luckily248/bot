package models

import (
	"database/sql"

	"os"

	_ "github.com/lib/pq"
)

type BasePGDBmodel struct {
	DB *sql.DB
}

func (this *BasePGDBmodel) init() (err error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return
	}
	this.DB = db
	return
}
