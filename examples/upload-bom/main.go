package main

import (
	"context"
	"encoding/base64"
	"flag"
	"log"
	"os"
	"time"

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
		wait           bool
	)

	flag.StringVar(&baseURL, "url", "", "Dependency-Track URL")
	flag.StringVar(&apiKey, "apikey", "", "Dependency-Track API key")
	flag.StringVar(&bomFilePath, "bom", "", "BOM file path")
	flag.StringVar(&projectName, "name", "", "Project name")
	flag.StringVar(&projectVersion, "version", "", "Project version")
	flag.BoolVar(&autoCreate, "autocreate", false, "Create project if it doesn't exist")
	flag.BoolVar(&wait, "wait", false, "Wait for BOM processing to complete")
	flag.Parse()

	client, err := dtrack.NewClient(baseURL, dtrack.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}

	bomContent, err := os.ReadFile(bomFilePath)
	if err != nil {
		log.Fatalf("failed to read bom file: %v", err)
	}

	token, err := client.BOM.Upload(context.Background(), dtrack.BOMUploadRequest{
		ProjectName:    projectName,
		ProjectVersion: projectVersion,
		AutoCreate:     autoCreate,
		BOM:            base64.StdEncoding.EncodeToString(bomContent),
	})
	if err != nil {
		log.Fatalf("failed to upload bom: %v", err)
	}

	log.Printf("bom upload successful (token: %s)", token)

	if !wait {
		return
	}

	log.Println("waiting for bom processing to complete")
	doneChan := make(chan struct{})

	go func(ticker *time.Ticker, timeout <-chan time.Time) {
	loop:
		for {
			select {
			case <-ticker.C:
				processing, err := client.BOM.IsBeingProcessed(context.Background(), token)
				if err != nil {
					log.Fatalf("failed to check bom processing status: %v", err)
				}
				if !processing {
					break loop
				}
				log.Println("still waiting")
			case <-timeout:
				log.Fatalln("timeout exceeded")
			}
		}
		doneChan <- struct{}{}
	}(time.NewTicker(1*time.Second), time.After(10*time.Second))

	<-doneChan
}
