package cli

import (
	"flag"
	"log"
	"os"
)

func Cli() (bool, bool, string) {
	// TODO; use a more sensible cli lib.
	var showVersion bool
	var up bool
	var d bool
	var dockerComposeFile = "docker-compose.yml"
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
		"path to docker-compose.yml file.")

	flag.Parse()

	if showVersion {
		return true, followLogs, ""
	}
	if !up {
		log.Println("to use Meli, run: \n\n\t meli -up")
		os.Exit(0)
	}
	if d {
		followLogs = false
	}

	return false, followLogs, dockerComposeFile
}
