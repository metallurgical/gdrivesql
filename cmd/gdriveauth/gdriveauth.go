package main

import "github.com/metallurgical/gdrivesql/pkg"

func main() {
	// Just to get token.json for the very 1st time
	(&pkg.GoogleDrive{}).New()
}