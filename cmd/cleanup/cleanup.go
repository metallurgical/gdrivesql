package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"flag"
)

const tempDir = "./temp"

var tempPath string

func init() {
	flag.StringVar(&tempPath, "t", tempDir, "Set custom relative temporary folder path that store backup files and folder")
	flag.Parse()
}

func main() {
	if tempPath == "" {
		log.Fatalf("Temp path cant be empty")
	}
	files, err := ioutil.ReadDir(tempPath)
	if err != nil {
		log.Printf("Error reading directory: %v", err)
	}

	for _, file := range files {
		name := fmt.Sprintf("%s/%s", tempPath, string(file.Name()))
		if file.IsDir() {
			if err := os.Remove(name); err != nil {
				log.Println("Cannot remove directory: %s", file.Name())
			}
		} else {
			if file.Name() != ".gitignore" {
				if err := os.Remove(name); err != nil {
					log.Println("Cannot remove directory: %s", file.Name())
				}
			}
		}
	}
}