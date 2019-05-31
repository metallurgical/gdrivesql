package pkg

import (
	"github.com/labstack/gommon/log"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
)

type Database struct {
	Name []string `yaml:name,omitempty`
}

type FileSystem struct {
	Path []string `yaml:path,omitempty`
}

type DriveItems struct {
	Path string `yaml:path,omitempty`
	FileSystem bool `yaml:filesystem,omitempty`
	Files []string `yaml:files,omitempty`
}

type Gdrive struct {
	Config []DriveItems
}

// GetGdrive get gdrive config to upload into Google Drive
func GetGdrive() *Gdrive {
	reader, err := os.Open("configs/gdrive.yaml")
	if err != nil {
		log.Fatalf("Can't open configs/gdrive.yaml", err)
	}
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("Can't read configs/gdrive.yaml", err)
	}
	defer reader.Close()

	var gdrive Gdrive
	yaml.Unmarshal(buf, &gdrive)
	return &gdrive
}

// GetDatabases get all the databases's name to export into .sql
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

// GetFileSystems get all filesystem's list to compress
func GetFileSystems() *FileSystem {
	reader, err := os.Open("configs/filesystems.yaml")
	if err != nil {
		log.Fatalf("Can't open configs/filesystems.yaml", err)
	}
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("Can't read configs/filesystems.yaml", err)
	}
	defer reader.Close()

	var filesystems FileSystem
	yaml.Unmarshal(buf, &filesystems)
	return &filesystems
}