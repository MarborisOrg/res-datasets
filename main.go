package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic("Error getting user home directory:")
	}
	homeDir := usr.HomeDir
	targetDir := filepath.Join(homeDir, ".marboris")

	if err = os.RemoveAll(targetDir); err != nil {
		panic("Error cleaning target directory:")
	}

	if err = os.MkdirAll(targetDir, 0o755); err != nil {
		panic("Error creating directory:")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		panic("Error getting current directory:")
	}
	resDir := currentDir + "res"

	err = filepath.WalkDir(resDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(resDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(targetDir, relPath)

		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0o755); err != nil {
				return fmt.Errorf("Error creating directory %s: %w", destPath, err)
			}
		} else {
			if err := copyFile(path, destPath); err != nil {
				return fmt.Errorf("Error copying file from %s to %s: %w", path, destPath, err)
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println("Error copying files:", err)
	}
}

func copyFile(source, destination string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer func() {
		dstFile.Close()
		if err != nil {
			os.Remove(destination)
		}
	}()

	_, err = io.CopyBuffer(dstFile, srcFile, make([]byte, 4096)) // Use a buffer for optimized copying
	if err != nil {
		return err
	}

	if !isWindows() {
		srcInfo, err := os.Stat(source)
		if err != nil {
			return err
		}
		if err := os.Chmod(destination, srcInfo.Mode()); err != nil {
			return err
		}
	}

	return nil
}

func isWindows() bool {
	return os.PathSeparator == '\\'
}
