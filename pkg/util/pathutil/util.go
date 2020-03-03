package pathutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	tmpSuffix      = ".tmp"
	ownerRW        = 0600
	userRWXGroupRX = 0750
)

// EnsureDir attempts to create given directory, panics if it fails to do so
func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, userRWXGroupRX)
	}

	return nil
}

// AtomicWriteFile creates a temp file in which to write data, then calls syscall.Rename to swap it and write it on
// filename for an atomic write. On failure temp file is removed and panics.
func AtomicWriteFile(filename string, data []byte) error {
	tempFilePath := filename + tmpSuffix

	if _, err := os.Stat(filename); err == nil {
		if err := os.Remove(filename); err != nil {
			return fmt.Errorf("remove %s: %w", filename, err)
		}
	}

	if _, err := os.Stat(tempFilePath); err == nil {
		if err := os.Remove(filename); err != nil {
			return fmt.Errorf("remove %s: %w", filename, err)
		}
	}

	if err := ioutil.WriteFile(tempFilePath, data, ownerRW); err != nil {
		return err
	}

	if err := os.Rename(tempFilePath, filename); err != nil {
		return err
	}

	return nil
}

// AtomicAppendToFile calls AtomicWriteFile but appends new data to destiny file
func AtomicAppendToFile(filename string, data []byte) error {
	oldFile, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		return err
	}

	return AtomicWriteFile(filename, append(oldFile, data...))
}
