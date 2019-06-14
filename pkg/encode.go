package pkg

import (
	"github.com/labstack/gommon/log"
	"github.com/go-yaml/yaml"

	"io/ioutil"
	"errors"
	"os"
)

const (
	gdrivePath = "configs/gdrive.yaml"
	databasePath = "configs/databases.yaml"
	filesystemPath = "configs/filesystems.yaml"
)

type (
	// Specifying possible connection to database server
	// later able to pick which connection can be used
	// to export .database
	Connection struct {
		Name string `yaml:name,omitempty`
		Driver string `yaml:driver,omitempty`
		Host string `yaml:host,omitempty`
		Port string `yaml:port,omitempty`
		User string `yaml:user,omitempty`
		Password string `yaml:password,omitempty`
	}

	Config struct {
		Connection string `yaml:connection,omitempty`
		List []string `yaml:list,omitempty`
	}

	Database struct {
		Connections []Connection
		Databases []Config
	}

	FileSystem struct {
		Path []string `yaml:path,omitempty`
	}

	DriveItems struct {
		Folder string `yaml:folder,omitempty`
		FileSystem bool `yaml:filesystem,omitempty`
		Files []string `yaml:files,omitempty`
		DriveId string `yaml:driveid,omitempty`
	}

	Gdrive struct {
		Config []DriveItems
	}
)

// GetGdrive get gdrive config to upload into Google Drive
func GetGdrive() (*Gdrive, error) {
	if !Exists(gdrivePath) {
		return nil, errors.New("File gdrive.yaml does not exist")
	}
	reader, err := os.Open(gdrivePath)
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
	return &gdrive, nil
}

// GetDatabases get all the databases's name to export into .sql
func GetDatabases() (*Database, error) {
	if !Exists(databasePath) {
		return nil, errors.New("File databases.yaml does not exist")
	}
	reader, err := os.Open(databasePath)
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

	return &databases, nil
}

// GetFileSystems get all filesystem's list to compress
func GetFileSystems() (*FileSystem, error) {
	if !Exists(filesystemPath) {
		return nil, errors.New("File filesystems.yaml does not exist")
	}
	reader, err := os.Open(filesystemPath)
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

	return &filesystems, nil
}