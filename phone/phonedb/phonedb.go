package phonedb

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jedib0t/go-pretty/table"
)

type pgdb struct {
	Db *pgxpool.Pool
}

type Contact struct {
	Phone_Number string // named this way to match database column
	Id           int
}

func NewDB(ctx context.Context, connStr string) (*pgdb, error) {
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &pgdb{db}, nil
}

func (pg *pgdb) PopulateDb(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "phone: error opening file for read: %v", err)
		return err
	}
	defer f.Close()

	contacts := [][]any{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		contacts = append(contacts, []any{s.Text()})
	}

	// Table Name: phone_numbers
	// Table Columns: id SERIAL PRIMARY KEY, phone_number TEXT
	_, err = pg.Db.CopyFrom(
		context.Background(),
		pgx.Identifier{"phone_numbers"},
		[]string{"phone_number"},
		pgx.CopyFromRows(contacts),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "phone: error writing into database: %v\n", err)
		return err
	}

	fmt.Printf("SUCCESS: Data from %s written to database.\n", filename)
	return nil
}

func (pg *pgdb) PrintData() error {
	contacts, err := pg.ReadDb()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Phone Number"})
	for _, contact := range contacts {
		t.AppendRow(table.Row{contact.Id, contact.Phone_Number})
	}
	t.Render()

	return nil
}

func (pg *pgdb) ResetDatabase() error {
	batch := &pgx.Batch{}
	batch.Queue("TRUNCATE phone_numbers")
	batch.Queue("ALTER SEQUENCE phone_numbers_id_seq RESTART")

	br := pg.Db.SendBatch(context.Background(), batch)

	_, err := br.Exec()
	if err != nil {
		return err
	}

	return br.Close()
}

func (pg *pgdb) LocateRecord(number string) (*Contact, error) {
	// var c Contact

	query := `SELECT * FROM phone_numbers WHERE phone_number = @phone_number`
	args := pgx.NamedArgs{
		"phone_number": number,
	}

	rows, err := pg.Db.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Contact])
	if err != nil {
		return nil, err
	}

	if len(contacts) == 0 {
		return nil, nil
	}

	return &contacts[0], nil
}

func (pg *pgdb) DeleteRecord(id int) error {
	query := `DELETE FROM phone_numbers WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := pg.Db.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}
	return nil
}

func (pg *pgdb) UpdateRecord(contact *Contact) error {
	var err error

	query := `UPDATE phone_numbers SET phone_number = @phone_number WHERE id = @id`
	args := pgx.NamedArgs{
		"id":           contact.Id,
		"phone_number": contact.Phone_Number,
	}
	_, err = pg.Db.Exec(context.Background(), query, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "phone: error updating database row: %v\n", err)
		return err
	}

	fmt.Printf("SUCCESS: %s updated.\n", contact.Phone_Number)
	return nil
}

func (pg *pgdb) ReadDb() ([]Contact, error) {
	query := `SELECT id, phone_number FROM phone_numbers`

	rows, err := pg.Db.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "phone: error reading database: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	contacts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Contact])
	if err != nil {
		fmt.Fprintf(os.Stderr, "phone: error parsing rows to struct: %v\n", err)
		return nil, err
	}

	return contacts, nil
}
