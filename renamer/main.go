package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type FileEntry struct {
	Info os.FileInfo
	Path string
}

func main() {
	fmt.Println("File renaming utility")

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "renamer: error getting pwd: %v", err)
		os.Exit(1)
	}

	fileRoot := flag.String("r", pwd, "root directory to begin search")

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

	for _, f := range walkResults {
		fmt.Println(f.Path, f.Info.Name())
	}
}
