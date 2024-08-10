package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	dataFile := flag.String("d", "", "file containing phone numbers for write to database")
	flag.Parse()

	populateDb(*dataFile)
}

func populateDb(filename string) error {
	fmt.Println(filename)

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
