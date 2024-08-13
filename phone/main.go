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
	// TODO: implement connection pools to avoid creating frequent new connections
	dataFile := flag.String("d", "", "file containing phone numbers for write to database")
	flag.Parse()

	// DB URL format: postgres://<db-user>:<password>@<ip/host>:<port>/<database-name>
	dbUrl := os.Getenv("DATABASE_URL")

	db, _ := NewDB(context.Background(), dbUrl)

	if isFlagPassed("d") {
		if err := db.populateDb(*dataFile); err != nil {
			os.Exit(1)
		}
		return
	}

	contacts, err := readDb(dbUrl)
	if err != nil {
		os.Exit(1)
	}

	if err = updateDb(dbUrl, contacts); err != nil {
		os.Exit(1)
	}
}

func NewDB(ctx context.Context, connStr string) (*pgdb, error) {
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating new PGX Pool: %v", err)
		os.Exit(1)
	}
	return &pgdb{db}, nil
}

// normalizes phone number -> ##########
func normalizeNumber(n string) string {
	norm := regexp.MustCompile(`[^0-9]+`).ReplaceAllString(n, "")
	return norm
}

func (pg *pgdb) populateDb(filename string) error {
	//conn, err := pgx.Connect(context.Background(), dbUrl)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "phone: error opening database connection: %v", err)
	//	return err
	//}
	//defer conn.Close(context.Background())

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

func readDb(dbUrl string) ([]Contact, error) {
	// TODO: put db connection code in separate function or possibly package for reuse
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "phone: error opening database connection: %v", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	query := `SELECT id, phone_number FROM phone_numbers`

	rows, err := conn.Query(context.Background(), query)
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

func updateDb(dbUrl string, contacts []Contact) error {
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "phone: error opening database connection: %v\n", err)
		return err
	}
	defer conn.Close(context.Background())

	for _, contact := range contacts {
		normNum := normalizeNumber(contact.Phone_Number)
		if normNum == contact.Phone_Number {
			continue
		}

		query := `UPDATE phone_numbers SET phone_number = @phone_number WHERE id = @id`
		args := pgx.NamedArgs{
			"id":           contact.Id,
			"phone_number": normNum,
		}
		_, err = conn.Exec(context.Background(), query, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "phone: error updating database row: %v\n", err)
			return err
		}
	}

	fmt.Println("SUCCESS: Phone numbers in database have been normalized.")
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
