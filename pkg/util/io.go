package util

import (
	"os"
)

func PurgeFile(file string) {
	os.Remove(file)
}

func AppendFile(file string, content []byte) error {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write(content); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
