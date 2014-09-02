package main

import (
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"os/user"
	"path"
)

var opts struct {
	NoOp bool `short:"n" long:"noop" description:"Do not actually copy or move files." default:"false"`
	Move bool `short:"m" long:"move" description:"Removes imported pictures." default:"false"`
}

func main() {
	me, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to determine current user: %s\n", err)
	}

	pimp := &Pimp{path.Join(me.HomeDir, "Pictures"), "Local", false}

	importPaths, err := flags.Parse(&opts)
	if err != nil {
		log.Printf("Error parsing command-line arguments: %s", err)
		os.Exit(1)
	}

	/*if len(importPaths) < 1 {
		importPaths = append(importPaths, os.Getwd())
	}*/

	if opts.NoOp {
		log.Println("No-copy mode enabled.")
	}
	if opts.Move {
		log.Println("Imported files will be deleted")
		pimp.RemoveImported = true
	}

	if !IsDirectory(pimp.TargetDir) {
		log.Fatalf("Target dir not found: %s", pimp.TargetDir)
	}

	pimp.ImportPaths(importPaths)
}
