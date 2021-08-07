package main

import (
	"context"
	"encoding/base64"
	"flag"
	"log"
	"os"

	"github.com/nscuro/dtrack-client"
)

func main() {
	var (
		baseURL        string
		apiKey         string
		bomFilePath    string
		projectName    string
		projectVersion string
		autoCreate     bool
	)

	flag.StringVar(&baseURL, "url", "", "Dependency-Track URL")
	flag.StringVar(&apiKey, "apikey", "", "Dependency-Track API key")
	flag.StringVar(&bomFilePath, "bom", "", "BOM file path")
	flag.StringVar(&projectName, "name", "", "Project name")
	flag.StringVar(&projectVersion, "version", "", "Project version")
	flag.BoolVar(&autoCreate, "autocreate", false, "Create project if it doesn't exist")
	flag.Parse()

	client, err := dtrack.NewClient(baseURL, dtrack.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}

	bomContent, err := os.ReadFile(bomFilePath)
	if err != nil {
		log.Fatalf("failed to read bom file: %v", err)
	}

	token, err := client.UploadBOM(context.Background(), dtrack.BOMUploadRequest{
		ProjectName:    projectName,
		ProjectVersion: projectVersion,
		AutoCreate:     autoCreate,
		BOM:            base64.StdEncoding.EncodeToString(bomContent),
	})
	if err != nil {
		log.Fatalf("failed to upload bom: %v", err)
	}

	log.Printf("bom upload successful (token: %s)", token)
}
