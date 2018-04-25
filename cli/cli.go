/*
Package cli provides the command line interface for the meli application.
*/
package cli

import (
	"flag"
	"fmt"
	"os"
)

// Cli parses input from stdin
func Cli() (bool, bool, bool, string) {
	// TODO; use a more sensible cli lib.
	var showVersion bool
	var up bool
	var d bool
	var build bool
	var dockerComposeFile = "docker-compose.yml"
	var followLogs = true
	var rebuild = false

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
	flag.BoolVar(
		&build,
		"build",
		false,
		"Rebuild services")
	flag.StringVar(
		&dockerComposeFile,
		"f",
		"docker-compose.yml",
		"path to docker-compose.yml file.")

	flag.Parse()

	if showVersion {
		return true, followLogs, rebuild, ""
	}
	if !up {
		fmt.Println("to use Meli, run: \n\n\t meli -up")
		os.Exit(0)
	}
	if d {
		followLogs = false
	}
	if build {
		rebuild = true
	}

	return false, followLogs, rebuild, dockerComposeFile
}
