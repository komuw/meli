package main

import (
	"flag"
	"log"
	"os"
)

func Cli() {
	var showVersion bool
	var up bool

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
		"Builds, re/creates, starts, and attaches to containers for a service")

	flag.Parse()

	if showVersion {
		log.Println("Meli version:", Version)
		os.Exit(0)
	}
	if !up {
		os.Exit(0)
	}
}
