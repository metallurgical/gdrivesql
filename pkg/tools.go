package pkg

import (
	"fmt"
	"os"
)

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Contains check string exist in slice
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Rename rename existing file and replace with new path
// same like move.
func Rename(oldPath string, path string, f os.FileInfo) error {
	return os.Rename(
		fmt.Sprintf("%s/%s", oldPath, string(f.Name())),
		fmt.Sprintf("%s/%s", path, f.Name()),
	)
}