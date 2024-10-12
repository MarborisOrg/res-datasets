package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"

	"marboris/util"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic("Error getting user home directory:")
	}
	homeDir := usr.HomeDir
	targetDir := filepath.Join(homeDir, ".marboris")

	// Remove the old directory if it exists
	if err = os.RemoveAll(targetDir); err != nil {
		panic("Error cleaning target directory:")
	}

	// Create the new directory
	if err = os.MkdirAll(targetDir, 0o755); err != nil {
		panic("Error creating directory:")
	}

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		panic("Error getting current directory:")
	}
	resDir := filepath.Join(currentDir, "res")

	// Copy the directory using system-specific commands
	err = copyDir(resDir, targetDir)
	if err != nil {
		fmt.Println("Error copying directory:", err)
	} else {
		util.RunChecker()
		fmt.Println("Datasets saved successfully! We are ready to go")
	}
}

// copyDir uses system-specific commands to copy directories
func copyDir(src, dst string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Use xcopy command on Windows
		cmd = exec.Command("xcopy", src, dst, "/E", "/I", "/Y")
	} else {
		// Use cp command on Unix-based systems
		cmd = exec.Command("cp", "-r", src, dst)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("copy error: %s", string(output))
	}

	return nil
}
