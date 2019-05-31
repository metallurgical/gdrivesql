package main

import (
	"github.com/metallurgical/gdrivesql/pkg"
	"github.com/mholt/archiver"

	"fmt"
	"io/ioutil"
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	Timezone = "Asia/Kuala_Lumpur"
	DateFormat = "2006-01-02@03-04-05PM"
)

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
		go compress(path, &wg)
	}
	wg.Wait()

	var w sync.WaitGroup
	w.Add(len(gdrive.Config))
	for _, config := range gdrive.Config {
		go backup(config, &w)
	}
	w.Wait()

	fmt.Print("Main program exit!")
}

func backup(items pkg.DriveItems, wg *sync.WaitGroup) {
	log.Printf("Backup all items:  %v", items)
	wg.Done()

	files, err := ioutil.ReadDir("./temp")
	if err != nil {
		log.Fatalf("Cannot read temp directory: ", err)
	}

	for _, f := range files {
		log.Printf("File name is %s : ", f.Name())
	}
	log.Print("huhu")
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

func compress(path string, wg *sync.WaitGroup) {
	log.Printf("Compressing path:  %v", string(path))
	defer wg.Done()

	files := []string{path}
	loc, err := time.LoadLocation(Timezone)
	if err != nil {
		log.Fatalf("Error loaded timezone: ", err)
	}
	currentDate := time.Now().In(loc).Format(DateFormat)

	slicePath := strings.Split(string(path), "/")
	tempFileName := slicePath[len(slicePath) - 1]
	filename := fmt.Sprintf("./temp/%s_%s.tar.gz", tempFileName, currentDate)
	// archive format is determined by file extension
	err = archiver.Archive(files, filename)
	if err != nil {
		log.Fatal(err)
	}
}