package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgdb struct {
	db *pgxpool.Pool
}

type Contact struct {
	Phone_Number string // named this way to match database column
	Id           int
}

func main() {
	// TODO: add flag to list database rows & related logic
	// TODO: add flag to normalize numbers & related logic
	dataFile := flag.String("d", "", "data file containing phone numbers for write to database")
	normDb := flag.Bool("n", false, "normalize database phone numbers")
	flag.Parse()

	// DB URL format: postgres://<db-user>:<password>@<ip/host>:<port>/<database-name>
	dbUrl := os.Getenv("DATABASE_URL")

	connPool, err := NewDB(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating new PGX Pool: %v", err)
		os.Exit(1)
	}
	defer connPool.db.Close()

	if isFlagPassed("d") {
		if err := connPool.PopulateDb(*dataFile); err != nil {
			os.Exit(1)
		}
		return
	}

	if *normDb {
		contacts, err := connPool.readDb()
		if err != nil {
			os.Exit(1)
		}

		for _, contact := range contacts {
			normNum := normalizeNumber(contact.Phone_Number)
			if normNum == contact.Phone_Number {
				continue
			}

			contact.Phone_Number = normNum
			if err := connPool.UpdateRecord(&contact); err != nil {
				fmt.Fprintf(os.Stderr, "phone: error updating DB record: %v", err)
			}
		}
	}
}

func NewDB(ctx context.Context, connStr string) (*pgdb, error) {
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &pgdb{db}, nil
}

// normalizes phone number -> ##########
func normalizeNumber(n string) string {
	norm := regexp.MustCompile(`[^0-9]+`).ReplaceAllString(n, "")
	return norm
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
	_, err = pg.db.CopyFrom(
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

func (pg *pgdb) readDb() ([]Contact, error) {
	query := `SELECT id, phone_number FROM phone_numbers`

	rows, err := pg.db.Query(context.Background(), query)
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

func (pg *pgdb) UpdateRecord(contact *Contact) error {
	var err error

	query := `UPDATE phone_numbers SET phone_number = @phone_number WHERE id = @id`
	args := pgx.NamedArgs{
		"id":           contact.Id,
		"phone_number": contact.Phone_Number,
	}
	_, err = pg.db.Exec(context.Background(), query, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "phone: error updating database row: %v\n", err)
		return err
	}

	fmt.Printf("SUCCESS: %s updated.\n", contact.Phone_Number)
	return nil
}

// Checks if flag was passed
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
