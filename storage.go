package promotions

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"

	_ "github.com/lib/pq"
)

const tbl = "t_promotions"

type DB struct {
	conn *sql.DB
}

func NewDBConn() DB {
	return DB{
		conn: connect(),
	}
}

// Postgres connect
func connect() *sql.DB {
	connectionStr := "user=user password=user dbname=promotions sslmode=disable" // same creds as in docker compose
	conn, err := sql.Open("postgres", connectionStr)
	if err != nil {
		panic(err)
	}
	return conn
}

func (d DB) getById(ctx context.Context, id int) (empty Row, err error) {
	row, err := d.conn.QueryContext(ctx, "select * from "+tbl+" where id = $1", id)
	defer row.Close()

	var res Row
	for row.Next() {
		if err = row.Scan(&res.ID, &res.Key, &res.Price, &res.ExpirationDate); err != nil {
			return empty, err
		}
	}
	return res, nil
}

func (d DB) batchInsert(ctx context.Context, file *os.File) (err error) {
	identifier := fmt.Sprintf("%s(key,price,expiration_date)", tbl)
	if err = copyFrom(ctx, d.conn, identifier, file); err != nil {
		return err
	}
	return nil
}

// Postgres COPY function
// this function allows moving data between tables and standard file-system files or buffers
// commits csv file into the table
func copyFrom(ctx context.Context, db *sql.DB, table string, r io.Reader) error {
	query := fmt.Sprintf("copy %s from stdin with (format csv)", table)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if _, err = stmt.ExecContext(ctx, sc.Text()); err != nil {
			return err
		}
	}
	if err = sc.Err(); err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx); err != nil {
		return err
	}
	return tx.Commit()
}

// Clear the database table and reset auto-increment (serial) counter
func (d DB) clearDB(ctx context.Context) (err error) {
	_, err = d.conn.ExecContext(
		ctx,
		"truncate table "+tbl+"; alter sequence t_promotions_id_seq restart 1;",
	)
	return err
}
