package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Welcome to Link Parser")
	filename := flag.String("f", "", "file containing HTML")
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	html, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalf("link: problem reading html file: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(html))
}
