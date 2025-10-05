package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

// Structures for parsing YAML
type fieldDef struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Pointer bool   `yaml:"pointer,omitempty"`
}

type typeDef struct {
	Name   string     `yaml:"name"`
	NodeID uint32     `yaml:"node_id"`
	Fields []fieldDef `yaml:"fields"`
}

type opcConfig struct {
	Namespace string    `yaml:"namespace"`
	Types     []typeDef `yaml:"types"`
}

func main() {
	// File paths
	yamlPath := filepath.FromSlash("tools/opc_codegen/opc_types.yaml")
	targetFile := filepath.FromSlash("pkg/opc_custom/models.go")
	mtimeFile := filepath.FromSlash("tools/opc_codegen/.opc_types.mtime")
	templatePath := filepath.FromSlash("tools/opc_codegen/codegen_custom_template.txt")

	log.Println("INFO: Starting code generation.")
	log.Printf("INFO: YAML file: %s", yamlPath)
	log.Printf("INFO: Template file: %s", templatePath)
	log.Printf("INFO: Output file: %s", targetFile)

	// Check if YAML exists
	yamlStat, err := os.Stat(yamlPath)
	if err != nil {
		log.Fatalf("ERROR: Could not read YAML file %s: %v", yamlPath, err)
	}
	yamlMTime := yamlStat.ModTime()

	// Check if generation is needed
	if prevMTimeBytes, err := os.ReadFile(mtimeFile); err == nil {
		if prevMTime, err := time.Parse(time.RFC3339, strings.TrimSpace(string(prevMTimeBytes))); err == nil {
			if !yamlMTime.After(prevMTime) {
				log.Println("INFO: YAML file has not changed — skipping generation.")
				return
			}
		} else {
			log.Printf("WARNING: Failed to parse previous mtime: %v", err)
		}
	}

	// Read YAML
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatalf("ERROR: Failed to read YAML file %s: %v", yamlPath, err)
	}

	var cfg opcConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("ERROR: Failed to parse YAML: %v", err)
	}
	log.Printf("INFO: Found %d types to generate.", len(cfg.Types))

	// Normalize ua.* types
	for i := range cfg.Types {
		for j := range cfg.Types[i].Fields {
			t := cfg.Types[i].Fields[j].Type
			if strings.HasPrefix(t, "ua.") {
				cfg.Types[i].Fields[j].Type = "ua." + strings.TrimPrefix(t, "ua.")
			}
		}
	}

	// Read template file
	tmplData, err := os.ReadFile(templatePath)
	if err != nil {
		log.Fatalf("ERROR: Failed to read template file %s: %v", templatePath, err)
	}

	// Create directory if it does not exist
	if err := os.MkdirAll(filepath.Dir(targetFile), 0755); err != nil {
		log.Fatalf("ERROR: Failed to create directory for %s: %v", targetFile, err)
	}

	// Create output file
	f, err := os.Create(targetFile)
	if err != nil {
		log.Fatalf("ERROR: Failed to create file %s: %v", targetFile, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("WARNING: Failed to close file %s: %v", targetFile, err)
		}
	}()

	// Parse template
	tmpl, err := template.New("gen").Parse(string(tmplData))
	if err != nil {
		log.Fatalf("ERROR: Failed to parse template: %v", err)
	}

	// Execute template
	if err := tmpl.Execute(f, cfg); err != nil {
		log.Fatalf("ERROR: Failed to generate code: %v", err)
	}

	// Update mtime file
	if err := os.WriteFile(mtimeFile, []byte(yamlMTime.Format(time.RFC3339)), 0644); err != nil {
		log.Printf("WARNING: Failed to update mtime file %s: %v", mtimeFile, err)
	}

	log.Printf("INFO: Code successfully generated: %s", targetFile)
}
