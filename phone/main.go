package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/jackc/pgx/v5"
)

func main() {
	dataFile := flag.String("d", "", "file containing phone numbers for write to database")
	flag.Parse()

	if isFlagPassed("d") {
		if err := populateDb(*dataFile); err != nil {
			fmt.Fprintf(os.Stderr, "phone: error populating database: %v", err)
			os.Exit(1)
		}
	}

	testNum := "(123) 456-7893"
	fmt.Println(normalizeNumber(testNum))
}

// normalizes phone number -> ##########
func normalizeNumber(n string) string {
	norm := regexp.MustCompile(`[^0-9]+`).ReplaceAllString(n, "")
	return norm
}

func populateDb(filename string) error {
	// DB URL format: postgres://<db-user>:<password>@<ip/host>:<port>/<database-name>
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "phone: error opening database connection: %v", err)
		return err
	}
	defer conn.Close(context.Background())

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
	_, err = conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"phone_numbers"},
		[]string{"phone_number"},
		pgx.CopyFromRows(contacts),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "phone: error writing into database: %v", err)
		return err
	}

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
