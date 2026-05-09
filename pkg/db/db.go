package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var DB *sqlx.DB

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(255) NOT NULL,
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX idx_date ON scheduler (date);
`

func Init(dbFile string) error {
	if env := os.Getenv("TODO_DBFILE"); env != "" {
		dbFile = env
	}

	_, err := os.Stat(dbFile)
	install := err != nil

	db, err := sqlx.Connect("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		db.MustExec(schema)
	}

	DB = db
	return nil
}
