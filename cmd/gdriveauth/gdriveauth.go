package main

import (
	"github.com/metallurgical/gdrivesql/pkg"

	"flag"
)

var credentialsDirPath string

func init() {
	// ** Notes: Compulsory to specify all these if executables files
	// runs outside of gdrivesql directory.
	flag.StringVar(&credentialsDirPath, "c", "credentials", "Set custom absolute credentials folder path that store credentials.json and token.json")
	flag.Parse()
}

func main() {
	// Just to get token.json for the very 1st time
	var gd = &pkg.GoogleDrive{
		CredentialDirPath: credentialsDirPath,
	}
	gd.New()
}