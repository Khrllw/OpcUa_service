package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type buildTarget struct {
	GOOS       string
	GOARCH     string
	OutputDir  string
	OutputName string
}

// Обновленный путь к исходникам
const sourcePath = "./cmd/app/app.go"

func main() {
	targets := []buildTarget{
		{
			GOOS:       "windows",
			GOARCH:     "amd64",
			OutputDir:  "./build",
			OutputName: "windows_opc_ua.exe",
		},
		{
			GOOS:       "linux",
			GOARCH:     "amd64",
			OutputDir:  "./build",
			OutputName: "linux_opc_ua",
		},
		{
			GOOS:       "darwin",
			GOARCH:     "amd64",
			OutputDir:  "./build",
			OutputName: "macos_opc_ua",
		},
	}

	log.Println("Starting build process...")

	for _, target := range targets {
		log.Printf("Building for %s/%s...", target.GOOS, target.GOARCH)

		if err := os.MkdirAll(target.OutputDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create directory %s: %v", target.OutputDir, err)
		}

		outputPath := filepath.Join(target.OutputDir, target.OutputName)
		cmd := exec.Command("go", "build", "-o", outputPath, sourcePath)

		cmd.Env = append(os.Environ(),
			fmt.Sprintf("GOOS=%s", target.GOOS),
			fmt.Sprintf("GOARCH=%s", target.GOARCH),
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("ERROR building for %s/%s.", target.GOOS, target.GOARCH)
			log.Fatalf("Command failed with error: %v\nOutput:\n%s", err, string(output))
		}

		log.Printf("Successfully built: %s", outputPath)
	}

	log.Println("All builds completed successfully!")
}
