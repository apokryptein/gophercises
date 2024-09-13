package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type FileEntry struct {
	Info os.FileInfo
	Path string
}

func main() {
	fmt.Println("File renaming utility")

	// Get current directory to act as default
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "renamer: error getting pwd: %v", err)
		os.Exit(1)
	}

	// Define and parse CL flags
	fileRoot := flag.String("r", pwd, "root directory to begin search")
	regexPattern := flag.String("p", ".*_[0-9]{3}\\.[a-zA-Z]{3}$", "regex file matching pattern")
	flag.Parse()

	// Walk files/directories from fileRoot
	// Append files to []FileEntry
	var walkResults []FileEntry
	err = filepath.WalkDir(*fileRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		walkResults = append(walkResults, FileEntry{Path: path, Info: info})
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "renamer: error walking directory: %v", err)
		os.Exit(1)
	}

	// Functionality check for now:
	// Iterate through results and check using regex
	for _, f := range walkResults {
		// fmt.Println(f.Path, f.Info.Name())
		if checkFilename(f.Info.Name(), *regexPattern) {
			fmt.Println("Match: ", f.Info.Name())
		}
	}
}

// Check filename string against supplied regex
func checkFilename(filename string, rx string) bool {
	re := regexp.MustCompile(rx)
	match := re.MatchString(filename)
	return match
}
