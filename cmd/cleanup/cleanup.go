package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	tempDir = "./temp"
)

func main() {
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		log.Printf("Error reading directory: %v", err)
	}

	for _, file := range files {
		name := fmt.Sprintf("./temp/%s", string(file.Name()))
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