package pkg

import (
	"github.com/labstack/gommon/log"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
)

type Database struct {
	Name []string `yaml:a,omitempty`
}

type FileSystem struct {
	Path []string `yaml:a,omitempty`
}

// GetDatabases get all the databases from config yaml file
func GetDatabases() *Database {
	reader, err := os.Open("configs/databases.yaml")
	if err != nil {
		log.Fatalf("Can't open configs/databases.yaml", err)
	}
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("Can't read configs/databases.yaml", err)
	}
	defer reader.Close()

	var databases Database
	yaml.Unmarshal(buf, &databases)
	return &databases
}