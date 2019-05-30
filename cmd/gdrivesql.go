package main

import (
	_"bufio"
	"fmt"
	"github.com/metallurgical/gdrivesql/pkg"
	"io/ioutil"
	_ "io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	databases := pkg.GetDatabases()
	dbConfig := pkg.NewConnection()

	channel := make(chan int, len(databases.Name))
	for _, name := range databases.Name {
		go dumping(name, dbConfig, channel)
	}

	k := <- channel

	fmt.Print(k)
}

// dumping dump sql output from stdout into each of
// database's file name with .sql
func dumping (name string, dbConfig *pkg.Connection, done chan<- int) {
	args := []string{
		fmt.Sprintf("--port=%s", dbConfig.Port),
		fmt.Sprintf("--host=%s", dbConfig.Host),
		fmt.Sprintf("--password=%s", ""),
	}
	args = append(args, name)
	log.Printf("exec mysqldump with %v", args)
	cmd := exec.Command("mysqldump", args...)
	stdout, err := cmd.StdoutPipe()
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
	err = ioutil.WriteFile(fmt.Sprintf("./temp/%s.sql", name), bytes, 0777)
	if err != nil {
		log.Fatalf("Cannot write to file: ", err)
	}

	done <- 1
}