package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type FileEntry struct {
	Info os.FileInfo
	Path string
}

func main() {
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

	// TODO: put WalkDir logic into separate function

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

		// Check if filename matches desired regex prior to storing in []FileEntry
		if checkFilename(info.Name(), *regexPattern) {
			walkResults = append(walkResults, FileEntry{Path: path, Info: info})
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "renamer: error walking directory: %v", err)
		os.Exit(1)
	}

	// Iterate through results and check using regex
	for _, f := range walkResults {
		err := renameFile(f, len(walkResults))
		if err != nil {
			fmt.Fprintf(os.Stderr, "renamer: error renaming file: %v", err)
		}
	}
}

// Check filename string against supplied regex
func checkFilename(filename string, rx string) bool {
	re := regexp.MustCompile(rx)
	match := re.MatchString(filename)
	return match
}

// Parses filename and renames file on OS
func renameFile(file FileEntry, n int) error {
	// Gather relevant indexes to parse filename
	i := strings.LastIndex(file.Info.Name(), "_")
	iExt := strings.LastIndex(file.Info.Name(), ".")

	// Grab relevant pieces: basename, number, extension
	fName := file.Info.Name()[:i]
	fExt := file.Info.Name()[iExt:]
	fNum, err := strconv.Atoi(file.Info.Name()[i+1 : iExt])
	if err != nil {
		return err
	}

	// Generate new filename: 'basename (# of total).ext'
	// TODO: make ending more dynamic
	newFilename := fmt.Sprintf("%s (%d of %d)%s", fName, fNum, n, fExt)
	fmt.Printf("Renaming '%s' to '%s'\n", file.Info.Name(), newFilename)

	// Grab directory path
	iPath := strings.LastIndex(file.Path, "/")
	fPath := file.Path[:iPath+1]

	// Rename file
	err = os.Rename(fPath+file.Info.Name(), fPath+newFilename)
	return err
}
