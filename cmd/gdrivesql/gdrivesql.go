package main

import (
	"github.com/metallurgical/gdrivesql/pkg"
	"github.com/mholt/archiver"

	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	Timezone            = "Asia/Kuala_Lumpur"
	DateFormat          = "2006-01-02@03-04-05PM"
	TempDir             = "./temp"
	FinalBackupFileName = "backup.tar.gz"
)

type (
	// Store temporary data for uploading purpose.
	fileArchive struct {
		googleDriveId   string
		archiveFileName string
	}

	// Main container to store slice of fileArchive.
	gdriveHolder struct {
		container []fileArchive
	}
)

var (
	tempGdriveHolder gdriveHolder
	flagDb, flagFile = true, true
	dbConfig         []pkg.Connection
)

func main() {
	// Checking existence  of databases.yaml, filesystems.yaml and gdrive.yaml
	// If all these files not exist, main process should abort the
	// execution.

	// Marshaling databases.yaml
	databases, err := pkg.GetDatabases()
	if err != nil {
		log.Printf("%v, skip database backup. \n", err.Error())
		flagDb = false
	} else {
		if len(databases.Databases) == 0 {
			log.Println("Database's not defined. Skip database backup.")
			flagDb = false
		}
	}

	// Marshalling filesystems.yaml
	filesystems, err := pkg.GetFileSystems()
	if err != nil {
		log.Printf("%v, skip filesystem backup. \n", err.Error())
		flagFile = false
	} else {
		if len(filesystems.Path) == 0 {
			log.Println("Filesystem's path not defined. Skip filesystem backup.")
			flagFile = false
		}
	}

	if !flagFile && !flagDb {
		log.Println("Both filesystems.yaml and databases.yaml does not exit. Nothing to backup. Abort process! \n")
		os.Exit(1)
	}

	// Marshalling gdrive.yaml
	gdrive, err := pkg.GetGdrive()
	if err != nil {
		log.Printf("%v, abort process! \n", err.Error())
		os.Exit(1)
	}

	dbConfig = databases.Connections

	var wgBase, wgBackup sync.WaitGroup // Use waitGroup to wait all goroutines

	if flagDb {
		wgBase.Add(len(databases.Databases))
		for _, config := range databases.Databases {
			//for _, db := range config.config {
			go dumping(config, &wgBase) // Dumping database
			//}
		}
	}

	if flagFile {
		wgBase.Add(len(filesystems.Path))
		for _, path := range filesystems.Path {
			go zipped(path, &wgBase) // Compressing filesystem
		}
	}

	// Invoke if either one(databases and filesystem) exist
	wgBase.Wait()
	wgBackup.Add(len(gdrive.Config))

	for _, config := range gdrive.Config {
		go backup(config, &wgBackup) // Compressing uploaded filesystem
	}

	wgBackup.Wait()
	wgBackup.Add(1) // Adding one since only 1 left

	go upload(&wgBackup) // Do backup into gDrive

	wgBackup.Wait()
	fmt.Println("Main program exit!")
}

// backup backup file to google drive. Filesystem and database
// will put together under same folder and compressed. After compressed
// all the archives will upload into google drive.
func backup(items pkg.DriveItems, wg *sync.WaitGroup) {
	defer wg.Done()

	files, err := ioutil.ReadDir(TempDir)
	if err != nil {
		log.Fatalf("Cannot read temp directory: ", err)
	}

	pathDir := fmt.Sprintf("%s/%s", TempDir, string(items.Folder))
	if !pkg.Exists(pathDir) {
		if err := os.Mkdir(pathDir, 0755); err != nil {
			log.Printf("Cant create directory: ", pathDir)
		}
	}

	for _, f := range files {
		firstName := strings.Split(f.Name(), "_")
		if pkg.Contains(items.Files, firstName[0]) {
			if !f.IsDir() {
				ext := strings.Split(firstName[len(firstName)-1], ".")
				switch ext[len(ext)-1] {
				case "sql":
					if err := pkg.Rename(pathDir, f); err != nil {
						log.Printf("Cannot rename file path: %v", pathDir)
					}
				case "gz":
					if items.FileSystem {
						if err := pkg.Rename(pathDir, f); err != nil {
							log.Printf("Cannot rename file path: %v", pathDir)
						}
					}
				}
			}
		}
	}

	// Lastly compress all the final folder that
	// need to upload
	tarFileName := compress(pathDir)
	tempGdriveHolder.container = append(tempGdriveHolder.container, fileArchive{
		googleDriveId:   items.DriveId,
		archiveFileName: tarFileName,
	})
}

// upload upload file into google drive.
func upload(wg *sync.WaitGroup) {
	defer wg.Done()

	srv := (&pkg.GoogleDrive{}).New()
	loc, err := time.LoadLocation(Timezone)
	if err != nil {
		log.Fatalf("Error loaded timezone: ", err)
	}

	folderName := time.Now().In(loc).Format(DateFormat)

	for _, gdrive := range tempGdriveHolder.container {
		dir, err := pkg.CreateDir(srv, gdrive.googleDriveId, folderName)
		if err != nil {
			log.Println("Could not create dir: " + err.Error())
		}
		// Step 1. Open the file
		f, err := os.Open(fmt.Sprintf("%s/%s", TempDir, gdrive.archiveFileName))
		if err != nil {
			panic(fmt.Sprintf("cannot open file: %v", err))
		}
		defer f.Close()

		// User backup.tar.gz instead of gdrive.archiveFileName
		// to avoid confusion
		pkg.CreateFile(srv, FinalBackupFileName, f, dir.Id)
	}
}

// dumping dump sql output from stdout into each of
// database's file name with .sql
func dumping(config pkg.Config, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, db := range config.List {

		var name = db

		cmd, args := dumpDriver(name, config.Connection)
		log.Printf("exec command with %v", args)

		stdout, err := cmd.StdoutPipe()
		var out bytes.Buffer
		cmd.Stderr = &out
		if err != nil {
			log.Fatalf("Error to execute mysqlump command: ", err)
			os.Exit(1)
		}

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		bytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			log.Fatalf("Read error: ", err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
		loc, err := time.LoadLocation(Timezone)
		if err != nil {
			log.Fatalf("Error loaded timezone: ", err)
		}
		currentDate := time.Now().In(loc).Format(DateFormat)
		filename := fmt.Sprintf("%s/%s_%s.sql", TempDir, name, currentDate)
		err = ioutil.WriteFile(filename, bytes, 0777)
		if err != nil {
			log.Fatalf("Cannot write to file: ", err)
		}
	}
}

// dumpDriver decide which `sql` command that need
// to execute. Either mysql or postgresSQL.
func dumpDriver(dbName string, connectionName string) (*exec.Cmd, []string) {
	var args []string
	var cmd *exec.Cmd

	for _, c := range dbConfig {

		if connectionName == c.Name {
			log.Println("CONNECTION NAME: %s - %s", connectionName, c.Driver)
			// Mysql configurations
			if c.Driver == "mysql" {
				args = []string{
					fmt.Sprintf("--port=%s", c.Port),
					fmt.Sprintf("--host=%s", c.Host),
					fmt.Sprintf("--user=%s", c.User),
					fmt.Sprintf("--password=%s", c.Password),
				}
				args = append(args, dbName)

				cmd = exec.Command("mysqldump", args...)

			} else if c.Driver == "postgres" { // PostgresSQL configurations
				args = []string{
					fmt.Sprintf("--port=%s", c.Port),
					fmt.Sprintf("--host=%s", c.Host),
					fmt.Sprintf("--username=%s", c.User),
					//fmt.Sprintf("--password=%s", c.Password),
					fmt.Sprintf("--dbname=%s", dbName),
				}

				cmd = exec.Command("pg_dump", args...)
				//cmd = exec.Command(fmt.Sprintf("pg_dump --dbname=postgresql://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, dbName))
			}
		}
	}
	return cmd, args
}

// zipped compressed any filesystem
func zipped(path string, wg *sync.WaitGroup) {
	log.Printf("Compressing path:  %v", string(path))
	defer wg.Done()
	compress(path)
}

// compress archive the directory and file
func compress(path string) string {
	files := []string{path}
	loc, err := time.LoadLocation(Timezone)
	if err != nil {
		log.Fatalf("Error loaded timezone: ", err)
	}
	currentDate := time.Now().In(loc).Format(DateFormat)

	slicePath := strings.Split(string(path), "/")
	tempFileName := slicePath[len(slicePath)-1]
	filenameWithoutPath := fmt.Sprintf("%s_%s.tar.gz", tempFileName, currentDate)
	filename := fmt.Sprintf("%s/%s", TempDir, filenameWithoutPath)
	// archive format is determined by file extension
	err = archiver.Archive(files, filename)
	if err != nil {
		log.Fatal(err)
	}
	return filenameWithoutPath
}
