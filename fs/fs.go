package fs

import (
	cMD5 "crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FS struct{}

func (fs *FS) Exists(path string) (exists bool, err error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (fs *FS) Write(path string, contents io.Reader) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, contents); err != nil {
		return fmt.Errorf("failed to copy contents to file: %s", err)
	}
	return nil
}

func (fs *FS) CreateDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %s", path, err)
	}

	return nil
}

func (fs *FS) DeleteFilesExcept(path string, filenames []string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to list files: %s", err)
	}

	for _, file := range files {
		if !fs.fileInSet(file.Name(), filenames) {
			err := fs.RemoveFile(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (fs *FS) RemoveFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to remove file %s: %s", path, err)
	}

	return nil
}

func (fs *FS) MD5(path string) (md5 string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open %s: %s", path, err)
	}
	defer file.Close()

	hash := cMD5.New()

	if _, err = io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to read %s: %s", path, err)
	}

	return fmt.Sprintf("%x", hash.Sum([]byte{})), nil
}

func (fs *FS) Length(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s: %s", path, err)
	}
	defer file.Close()

	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

func (fs *FS) Move(source string, destination string) error {
	if err := os.Rename(source, destination); err != nil {
		return fmt.Errorf("failed to move %s to %s: %s", source, destination, err)
	}

	return nil
}

func (fs *FS) fileInSet(filenameToFind string, filenames []string) bool {
	for _, filename := range filenames {
		if filenameToFind == filename {
			return true
		}
	}
	return false
}
