package pkg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/labstack/gommon/log"
)

var (
	basePath                 = "configs"
	gdriveConfigFileName     = "gdrive.yaml"
	databaseConfigFileName   = "databases.yaml"
	filesystemConfigFileName = "filesystems.yaml"
	mailPathConfigFileName   = "mail.yaml"
	gdrivePath               = fmt.Sprintf("%s%s%s", basePath, string(os.PathSeparator), gdriveConfigFileName)
	databasePath             = fmt.Sprintf("%s%s%s", basePath, string(os.PathSeparator), databaseConfigFileName)
	filesystemPath           = fmt.Sprintf("%s%s%s", basePath, string(os.PathSeparator), filesystemConfigFileName)
	mailPath                 = fmt.Sprintf("%s%s%s", basePath, string(os.PathSeparator), mailPathConfigFileName)
)

type (
	// Specifying possible connection to database server
	// later able to pick which connection can be used
	// to export .database
	Connection struct {
		Name     string `yaml:"name,omitempty"`
		Driver   string `yaml:"driver,omitempty"`
		Host     string `yaml:"host,omitempty"`
		Port     string `yaml:"port,omitempty"`
		User     string `yaml:"user,omitempty"`
		Password string `yaml:"password,omitempty"`
	}

	Config struct {
		Connection string   `yaml:"connection,omitempty"`
		List       []string `yaml:"list,omitempty"`
	}

	Database struct {
		Connections []Connection
		Databases   []Config
	}

	FileSystem struct {
		Path []string `yaml:"path,omitempty"`
	}

	// Google drive information that need
	// to store into google drive.
	DriveItems struct {
		Folder     string   `yaml:"folder,omitempty"`
		FileSystem bool     `yaml:"filesystem,omitempty"`
		Files      []string `yaml:"files,omitempty"`
		DriveId    string   `yaml:"driveid,omitempty"`
	}

	Gdrive struct {
		Config []DriveItems
	}

	// Define all configurations path.
	Settings struct {
		GdrivePath     string
		DatabasePath   string
		FilesystemPath string
		ConfigPath     string
		MailPath       string
	}

	Mail struct {
		From       string `yaml:"from,omitempty"`
		To         string `yaml:"to,omitempty"`
		Host       string `yaml:"host,omitempty"`
		Port       string `yaml:"port,omitempty"`
		Username   string `yaml:"username,omitempty"`
		Password   string `yaml:"password,omitempty"`
		Encryption string `yaml:"encryption,omitempty"`
	}
)

// New return default settings.
func New() *Settings {
	return &Settings{
		ConfigPath:     "",
		GdrivePath:     gdrivePath,
		DatabasePath:   databasePath,
		FilesystemPath: filesystemPath,
		MailPath:       mailPath,
	}
}

// ConstructPath replace user defined path for
// all existing configurations path
func (s *Settings) ConstructPath() {
	if s.ConfigPath != "" {
		s.FilesystemPath = fmt.Sprintf("%s%s%s", s.ConfigPath, string(os.PathSeparator), filesystemConfigFileName)
		s.DatabasePath = fmt.Sprintf("%s%s%s", s.ConfigPath, string(os.PathSeparator), databaseConfigFileName)
		s.GdrivePath = fmt.Sprintf("%s%s%s", s.ConfigPath, string(os.PathSeparator), gdriveConfigFileName)
	}
}

// Configurations defined method should implement
// by configuration files.
type Configurations interface {
	GetConfig([]byte, *os.File) Configurations
	GetPath(s *Settings) string
}

// GetPath get the mail configuration file path.
func (c Mail) GetPath(s *Settings) string {
	return s.MailPath
}

// GetConfig get the marshalling config of Mail configuration.
func (c Mail) GetConfig(buf []byte, reader *os.File) Configurations {
	defer reader.Close()
	var conf Mail
	yaml.Unmarshal(buf, &conf)
	return conf
}

// GetPath get the database configuration file path.
func (c Database) GetPath(s *Settings) string {
	return s.DatabasePath
}

// GetConfig get the marshalling config of Database configuration.
func (c Database) GetConfig(buf []byte, reader *os.File) Configurations {
	defer reader.Close()
	var conf Database
	yaml.Unmarshal(buf, &conf)
	return conf
}

// GetPath get the filesystem configuration file path.
func (c FileSystem) GetPath(s *Settings) string {
	return s.FilesystemPath
}

// GetConfig get the marshalling config of Filesystem configuration.
func (c FileSystem) GetConfig(buf []byte, reader *os.File) Configurations {
	defer reader.Close()
	var conf FileSystem
	yaml.Unmarshal(buf, &conf)
	return conf
}

// GetPath get the gdrive configuration file path.
func (c Gdrive) GetPath(s *Settings) string {
	return s.GdrivePath
}

// GetConfig get the marshalling config of Gdrive configuration.
func (c Gdrive) GetConfig(buf []byte, reader *os.File) Configurations {
	defer reader.Close()
	var conf Gdrive
	yaml.Unmarshal(buf, &conf)
	return conf
}

// GetConfig react as layer to get implemented configurations.
func (s *Settings) GetConfig(c Configurations) (Configurations, error) {
	if _, err := s.isFileExist(c); err != nil {
		return nil, err
	}
	buf, reader := s.readFile(c.GetPath(s))

	return c.GetConfig(buf, reader), nil
}

// isFileExist check the configuration files is exist.
func (s *Settings) isFileExist(c Configurations) (bool, error) {
	if !Exists(c.GetPath(s)) {
		return false, errors.New(fmt.Sprintf("File %v does not exist", c.GetPath(s)))
	}
	return true, nil
}

// readFile open file and return underlying bytes.
func (s *Settings) readFile(path string) ([]byte, *os.File) {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Can't open %v", err))
	}
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Can't read %v", err))
	}

	return buf, reader
}
