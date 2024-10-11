package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}
	homeDir := usr.HomeDir

	targetDir := filepath.Join(homeDir, ".marboris", "res")

	if err = os.RemoveAll(targetDir); err != nil {
		fmt.Println("Error cleaning target directory:", err)
		return
	}

	if err = os.MkdirAll(targetDir, 0o755); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	excludedFiles := []string{
		"main.go",
		"go.mod",
		"go.sum",
		"README",
		".git",
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	err = filepath.WalkDir(currentDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded files and directories
		for _, excluded := range excludedFiles {
			if strings.HasSuffix(path, excluded) || strings.Contains(path, string(os.PathSeparator)+excluded) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		relPath, err := filepath.Rel(currentDir, path)
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
