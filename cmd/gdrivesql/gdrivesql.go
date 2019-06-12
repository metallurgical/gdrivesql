package main

import (
	"github.com/metallurgical/gdrivesql/pkg"
	"github.com/mholt/archiver"
	"google.golang.org/api/drive/v3"

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
	Timezone   = "Asia/Kuala_Lumpur"
	DateFormat = "2006-01-02@03-04-05PM"
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

var tempGdriveHolder gdriveHolder

func main() {
	databases := pkg.GetDatabases()
	filesystems := pkg.GetFileSystems()
	gdrive := pkg.GetGdrive()
	dbConfig := pkg.NewConnection()

	var wg sync.WaitGroup
	// How many goroutines need to be waited
	wg.Add(len(databases.Name))
	for _, name := range databases.Name {
		go dumping(name, dbConfig, &wg)
	}
	wg.Add(len(filesystems.Path))
	for _, path := range filesystems.Path {
		go zipped(path, &wg)
	}
	wg.Wait()

	var w sync.WaitGroup
	w.Add(len(gdrive.Config))
	for _, config := range gdrive.Config {
		go backup(config, &w)
	}
	w.Wait()

	w.Add(1)
	go upload(&w)
	w.Wait()

	fmt.Println("Main program exit!")
}

// backup backup file to google drive. Filesystem and database
// will put together under same folder and compressed. After compressed
// all the archives will upload into google drive.
func backup(items pkg.DriveItems, wg *sync.WaitGroup) {
	defer wg.Done()

	files, err := ioutil.ReadDir("./temp")
	if err != nil {
		log.Fatalf("Cannot read temp directory: ", err)
	}

	pathDir := fmt.Sprintf("./temp/%s", string(items.Folder))
	if !Exists(pathDir) {
		if err := os.Mkdir(pathDir, 0755); err != nil {
			log.Printf("Cant create directory: ", pathDir)
		}
	}

	for _, f := range files {
		firstName := strings.Split(f.Name(), "_")
		//log.Printf("Firstname: %v. Items: %v", firstName[0], items.Files)
		if contains(items.Files, firstName[0]) {
			if !f.IsDir() {
				ext := strings.Split(firstName[len(firstName)-1], ".")
				switch ext[len(ext)-1] {
				case "sql":
					if err := rename(pathDir, f); err != nil {
						log.Printf("Cannot rename file path: %v", pathDir)
					}
				case "gz":
					if items.FileSystem {
						if err := rename(pathDir, f); err != nil {
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
	//fmt.Print(tempGdriveHolder.container)
	defer wg.Done()

	srv := (&pkg.GoogleDrive{}).New()
	loc, err := time.LoadLocation(Timezone)
	if err != nil {
		log.Fatalf("Error loaded timezone: ", err)
	}

	folderName := time.Now().In(loc).Format(DateFormat)

	for _, gdrive := range tempGdriveHolder.container {
		dir, err := createDir(srv, gdrive.googleDriveId, folderName)
		if err != nil {
			log.Println("Could not create dir: " + err.Error())
		}
		// Step 1. Open the file
		f, err := os.Open(fmt.Sprintf("./temp/%s", gdrive.archiveFileName))
		if err != nil {
			panic(fmt.Sprintf("cannot open file: %v", err))
		}
		defer f.Close()

		// User backup.tar.gz instead of gdrive.archiveFileName
		// to avoid confusion
		createFile(srv, "backup.tar.gz", f, dir.Id)
	}
}

// dumping dump sql output from stdout into each of
// database's file name with .sql
func dumping(name string, c *pkg.Connection, wg *sync.WaitGroup) {
	defer wg.Done()
	args := []string{
		fmt.Sprintf("--port=%s", c.Port),
		fmt.Sprintf("--host=%s", c.Host),
		fmt.Sprintf("--user=%s", c.User),
		fmt.Sprintf("--password=%s", ""),
	}
	args = append(args, name)
	log.Printf("exec mysqldump with %v", args)
	cmd := exec.Command("mysqldump", args...)
	stdout, err := cmd.StdoutPipe()
	var out bytes.Buffer
	cmd.Stderr = &out
	if err != nil {
		log.Fatalf("Error to execute mysqlump command: ", err)
		os.Exit(1)
	}
	//fmt.Printf("%q\n", out.String()) // to log the real error

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
	filename := fmt.Sprintf("./temp/%s_%s.sql", name, currentDate)
	err = ioutil.WriteFile(filename, bytes, 0777)
	if err != nil {
		log.Fatalf("Cannot write to file: ", err)
	}
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
	filename := fmt.Sprintf("./temp/%s", filenameWithoutPath)
	// archive format is determined by file extension
	err = archiver.Archive(files, filename)
	if err != nil {
		log.Fatal(err)
	}
	return filenameWithoutPath
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// contains check string exist in slice
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// rename rename existing file and replace with new path
// same like move.
func rename(path string, f os.FileInfo) error {
	return os.Rename(
		fmt.Sprintf("./temp/%s", string(f.Name())),
		fmt.Sprintf("%s/%s", path, f.Name()),
	)
}

// createDir create directory under particular
// Parent ID inside google drive.
func createDir(srv *drive.Service, parentId string, folderName string) (*drive.File, error) {
	d := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}
	dir, err := srv.Files.Create(d).Do()
	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}
	return dir, nil
}

// createFile create file(upload) into google drive.
func createFile(srv *drive.Service, name string, fileToUpload *os.File, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: "application/tar+gzip",
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := srv.Files.Create(f).Media(fileToUpload).Do()
	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}
	return file, nil
}
