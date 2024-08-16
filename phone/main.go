package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/apokryptein/gophercises/phone/phonedb"
)

func main() {
	dataFile := flag.String("d", "", "data file containing phone numbers for write to database")
	normDb := flag.Bool("n", false, "normalize database phone numbers")
	resetDb := flag.Bool("r", false, "reset phone_numbers table")
	printData := flag.Bool("p", false, "print table data")
	flag.Parse()

	// DB URL format: postgres://<db-user>:<password>@<ip/host>:<port>/<database-name>
	dbUrl := os.Getenv("DATABASE_URL")

	connPool, err := phonedb.NewDB(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating new PGX Pool: %v", err)
		os.Exit(1)
	}
	defer connPool.Db.Close()

	// Populate database with phone numbers from file
	if isFlagPassed("d") {
		if err := connPool.PopulateDb(*dataFile); err != nil {
			os.Exit(1)
		}
		return
	}

	// Print table dataFile
	if *printData {
		if err := connPool.PrintData(); err != nil {
			fmt.Fprintf(os.Stderr, "phone: error printing table: %v", err)
			os.Exit(1)
		}
		return
	}

	// Reset phone_numbers table
	if *resetDb {
		if err := connPool.ResetDatabase(); err != nil {
			fmt.Fprintf(os.Stderr, "phone: error reseting database: %v", err)
			os.Exit(1)
		}
		fmt.Println("SUCCESS: database reset")
		return
	}

	// Normalize Phone Numbers per normDb flag
	if *normDb {
		// Read DB into []Contact
		contacts, err := connPool.ReadDb()
		if err != nil {
			os.Exit(1)
		}

		// Iterate through records to find numbers requiring normalization
		for _, contact := range contacts {
			normNum := normalizeNumber(contact.Phone_Number)

			// If normaliation is required
			if normNum != contact.Phone_Number {
				// Locate relevant record to find other records with the same normalized number
				lookup, err := connPool.LocateRecord(normNum)
				if err != nil {
					fmt.Fprintf(os.Stderr, "phone: error locating record: %v", err)
					os.Exit(1)
				}

				// If duplicate found, delete
				if lookup != nil {
					if err := connPool.DeleteRecord(contact.Id); err != nil {
						fmt.Fprintf(os.Stderr, "phone: error deleting record: %v", err)
						os.Exit(1)
					}
					fmt.Printf("DUPLICATE: Deleting record: %d:%s\n", lookup.Id, lookup.Phone_Number)
					continue
				}

				// If no duplicate found, update record
				contact.Phone_Number = normNum
				if err := connPool.UpdateRecord(&contact); err != nil {
					fmt.Fprintf(os.Stderr, "phone: error updating DB record: %v", err)
				}
				continue
			}
		}
	}
}

// normalizes phone number -> ##########
func normalizeNumber(n string) string {
	norm := regexp.MustCompile(`[^0-9]+`).ReplaceAllString(n, "")
	return norm
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
