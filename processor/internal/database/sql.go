package database_sql

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type Database struct {
	ctx context.Context
	db  *pgx.Conn
}

func New(uri string) (*Database, error) {

	ctx := context.Background()
	db, err := pgx.Connect(ctx, uri) // *sql.DB
	if err != nil {
		return nil, err
	}

	return &Database{ctx: ctx, db: db}, nil
}

func (db *Database) Migrate() error {

	query := `create table if not exists exam_results (
    id SERIAL,
    name VARCHAR(255),
    sirname VARCHAR(255),
    score double precision
)`

	_, err := db.db.Exec(db.ctx, query)
	if err != nil {
		return err
	}

	query = "CREATE INDEX IF NOT EXISTS exam_results_name_idx on exam_results using btree (name)"

	_, err = db.db.Exec(db.ctx, query)

	if err != nil {
		return err
	}

	query = "CREATE INDEX IF NOT EXISTS exam_results_sirname_idx on exam_results using btree (sirname)"

	_, err = db.db.Exec(db.ctx, query)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) SaveRecord(name string, sirname string, score int64) error {
	query := `insert into exam_results(name, sirname, score) values($1, $2, $3) RETURNING id`

	_, err := db.db.Exec(db.ctx, query, name, sirname, score)
	return err
}
