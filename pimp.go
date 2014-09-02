package main

import (
	"fmt"
	"log"
	"github.com/rwcarlsen/goexif/exif"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Pimp struct {
	TargetDir		string
	TimeZone		string
	RemoveImported	bool
}

type ImportResult	int

const (
	Success ImportResult = iota
	IsSameFile
	TargetExists
	ErrorCopying
	UnsupportedFormat
	NotReadable
)

var resultDesc = map[ImportResult]string {
	Success: 			"Success",
	IsSameFile: 		"Source/Destination are the same file",
	TargetExists: 		"Target file with same name exists",
	ErrorCopying: 		"Error copying data",
	UnsupportedFormat: 	"Unsupported input format",
	NotReadable: 		"Source file not readable",
}

func (pimp *Pimp) GetTargetPath(file os.FileInfo, dateTime string) (path string) {
	loc, err := time.LoadLocation(pimp.TimeZone)
	if err != nil {
		log.Fatal(err)
	}
	if t, err := time.ParseInLocation("2006:01:02 15:04:05", dateTime, loc); err == nil {
		return fmt.Sprintf("%s/%04d/%02d/%s", pimp.TargetDir, t.Year(), t.Month(), file.Name())
	}
	return fmt.Sprintf("%s/unkown/%s", pimp.TargetDir, file.Name())
}

func (pimp *Pimp) ImportPicture(path string, file os.FileInfo) ImportResult {
	ext := strings.TrimLeft(strings.ToLower(filepath.Ext(path)), ".")
	SUPPORTED := List {"jpg", "jpeg"}
	
	if (!file.Mode().IsRegular() || !SUPPORTED.Contains(ext) || file.Name()[0] == '.') {
		return UnsupportedFormat
	}
	
	var dateTime = ""
	f, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return NotReadable
	}

	x, err := exif.Decode(f)
	defer f.Close()
	if err == nil {
		if date, err := x.Get(exif.DateTimeOriginal); err == nil {
			dateTime = date.StringVal()
		} else if date, err := x.Get(exif.DateTime); err == nil {
			dateTime = date.StringVal()
		}
	}
	var targetPath = pimp.GetTargetPath(file, dateTime)

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		log.Fatalf("Error creating directory: %s", filepath.Dir(targetPath))
	}
	
	if target, err := os.Stat(targetPath); err == nil {
		if os.SameFile(file, target) {
			return IsSameFile
		}

		log.Printf("Different file with same name exists: %s <--> %s", path, targetPath)
		return TargetExists
	}
	
	log.Printf("- %s --> %s\n", path, targetPath)
	if err := CopyFile(path, targetPath); err != nil {
		log.Fatal(err)
	}
	
	return Success
}

func (pimp *Pimp) ImportPaths(paths []string) {
	count := make(map[ImportResult]uint32)
	for _, importDir := range paths {
		if !IsDirectory(importDir) {
			log.Printf("Directory no found: %s", importDir)
			continue
		}
		log.Printf("Scanning dir %s\n", importDir)

		filepath.Walk(importDir, func(path string, f os.FileInfo, err error) error {
			res := pimp.ImportPicture(path, f)
			if res == Success {
				if pimp.RemoveImported {
					os.Remove(path)
				}
			}
			count[res] += 1
			return nil
		})
	}
	
	log.Printf("%d files successfully imported.\n", count[Success])
	for k, v := range count {
		log.Printf("%s: %d\n", resultDesc[k], v)
	}
}
