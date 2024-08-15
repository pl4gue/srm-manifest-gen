package main

import (
	"io/fs"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"

	srm "github.com/pl4gue/srm-manifest-gen/internal"
)

func check(msg string, err error, keyvals ...interface{}) {
	if err == huh.ErrUserAborted {
		return
	}

	if err != nil {
		log.Fatal(msg, "err", err, keyvals)
	}
}

func setup_log() {
	os.Mkdir("log", fs.FileMode(os.O_RDWR))
	file, err := os.Create("log/log.txt")

	check("Encountered an error creating log file.", err)

	log.SetOutput(file)
}

func main() {
	app := srm.App{}

	setup_log()

	log.SetPrefix("Prompting user.")
	err := app.PromptRoot()
	check("Encountered a fatal error while getting root directory.", err)

	log.SetPrefix("Populating executables.")
	err = app.Populate()
	check("Encountered a fatal error while populating executables.", err)

	log.SetPrefix("Choosing executables.")
	err = app.ChooseExecutables()
	check("Encountered a fatal error while selecting and confirming executables.", err)

	log.SetPrefix("Parsing executables as JSON.")
	json_bytes, err := app.JSONify()
	check("Encountered a fatal error while parsing executables to JSON.", err)

	log.SetPrefix("Writing JSON to file.")
	err = srm.WriteToFile(json_bytes, "manifest.json")
	check("Encountered a fatal error while writing JSON bytes to file", err, "json_bytes", json_bytes)

	log.SetPrefix("")
	log.Debug("Program ended without errors.")
}
