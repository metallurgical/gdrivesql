package main

import (
	"bytes"
	"fmt"
	"github.com/metallurgical/gdrivesql/pkg"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
)

func main() {
	databases := pkg.GetDatabases()
	dbConfig := pkg.NewConnection()

	var wg sync.WaitGroup

	// How many goroutines need to be waited
	wg.Add(len(databases.Name))

	for _, name := range databases.Name {
		go dumping(name, dbConfig, &wg)
	}

	wg.Wait()

	fmt.Print("Main program exit!")
}

// dumping dump sql output from stdout into each of
// database's file name with .sql
func dumping (name string, c *pkg.Connection, wg *sync.WaitGroup) {
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
	err = ioutil.WriteFile(fmt.Sprintf("./temp/%s.sql", name), bytes, 0777)
	if err != nil {
		log.Fatalf("Cannot write to file: ", err)
	}
}