package cli

import (
	"flag"
	"log"
	"os"
)

func Cli() (bool, string) {
	// TODO; use a more sensible cli lib.
	var showVersion bool
	var up bool
	var d bool
	var dockerComposeFile string = "docker-compose.yml"
	var followLogs = true

	flag.BoolVar(
		&showVersion,
		"version",
		false,
		"Show version information.")
	flag.BoolVar(
		&showVersion,
		"v",
		false,
		"Show version information.")
	flag.BoolVar(
		&up,
		"up",
		false,
		"Builds, re/creates, starts, and attaches to containers for a service.")
	flag.BoolVar(
		&d,
		"d",
		false,
		"Run containers in the background")
	flag.StringVar(
		&dockerComposeFile,
		"f",
		"docker-compose.yml",
		"path to docker-compose.yml file. By default, meli checks the current directory.")

	flag.Parse()

	if showVersion {
		log.Println("Meli version:", "To be released soon..")
		os.Exit(0)
	}
	if !up {
		log.Println("to use Meli, run: \n\n\t meli -up")
		os.Exit(0)
	}
	if d {
		followLogs = false
	}

	return followLogs, dockerComposeFile
}
