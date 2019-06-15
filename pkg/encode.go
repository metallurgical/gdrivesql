package pkg

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/labstack/gommon/log"

	"errors"
	"io/ioutil"
	"os"
)

var (
	basePath                 = "configs"
	gdriveConfigFileName     = "gdrive.yaml"
	databaseConfigFileName   = "databases.yaml"
	filesystemConfigFileName = "filesystems.yaml"
	gdrivePath               = fmt.Sprintf("%s%s%s", basePath, string(os.PathSeparator), gdriveConfigFileName)
	databasePath             = fmt.Sprintf("%s%s%s", basePath, string(os.PathSeparator), databaseConfigFileName)
	filesystemPath           = fmt.Sprintf("%s%s%s", basePath, string(os.PathSeparator), filesystemConfigFileName)
	//databasePath   = "configs/databases.yaml"
	//filesystemPath = "configs/filesystems.yaml"
)

type (
	// Specifying possible connection to database server
	// later able to pick which connection can be used
	// to export .database
	Connection struct {
		Name     string `yaml:name,omitempty`
		Driver   string `yaml:driver,omitempty`
		Host     string `yaml:host,omitempty`
		Port     string `yaml:port,omitempty`
		User     string `yaml:user,omitempty`
		Password string `yaml:password,omitempty`
	}

	Config struct {
		Connection string   `yaml:connection,omitempty`
		List       []string `yaml:list,omitempty`
	}

	Database struct {
		Connections []Connection
		Databases   []Config
	}

	FileSystem struct {
		Path []string `yaml:path,omitempty`
	}

	DriveItems struct {
		Folder     string   `yaml:folder,omitempty`
		FileSystem bool     `yaml:filesystem,omitempty`
		Files      []string `yaml:files,omitempty`
		DriveId    string   `yaml:driveid,omitempty`
	}

	Gdrive struct {
		Config []DriveItems
	}

	Settings struct {
		GdrivePath, DatabasePath, FilesystemPath string
		ConfigPath                               string
	}
)

// New return default settings.
func New() *Settings {
	return &Settings{
		ConfigPath:     "",
		GdrivePath:     gdrivePath,
		DatabasePath:   databasePath,
		FilesystemPath: filesystemPath,
	}
}

func (s *Settings) ConstructPath() {
	if s.ConfigPath != "" {
		s.FilesystemPath = fmt.Sprintf("%s%s%s", s.ConfigPath, string(os.PathSeparator), filesystemConfigFileName)
		s.DatabasePath = fmt.Sprintf("%s%s%s", s.ConfigPath, string(os.PathSeparator), databaseConfigFileName)
		s.GdrivePath = fmt.Sprintf("%s%s%s", s.ConfigPath, string(os.PathSeparator), gdriveConfigFileName)
	}
}

// GetGdrive get gdrive config to upload into Google Drive
func (s *Settings) GetGdrive() (*Gdrive, error) {
	if !Exists(s.GdrivePath) {
		return nil, errors.New("File gdrive.yaml does not exist")
	}
	reader, err := os.Open(s.GdrivePath)
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
func (s *Settings) GetDatabases() (*Database, error) {
	if !Exists(s.DatabasePath) {
		return nil, errors.New("File databases.yaml does not exist")
	}
	reader, err := os.Open(s.DatabasePath)
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
func (s *Settings) GetFileSystems() (*FileSystem, error) {
	if !Exists(s.FilesystemPath) {
		return nil, errors.New("File filesystems.yaml does not exist")
	}
	reader, err := os.Open(s.FilesystemPath)
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
