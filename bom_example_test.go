package dtrack_test

import (
	"context"
	"encoding/base64"
	"github.com/nscuro/dtrack-client"
	"os"
	"sync"
	"time"
)

// This example demonstrates how to upload a Bill of Materials and wait for its processing to complete.
func Example_uploadBOM() {
	client, _ := dtrack.NewClient("https://dtrack.example.com", dtrack.WithAPIKey("..."))

	bomContent, err := os.ReadFile("bom.xml")
	if err != nil {
		panic(err)
	}

	uploadToken, err := client.BOM.Upload(context.TODO(), dtrack.BOMUploadRequest{
		ProjectName:    "acme-app",
		ProjectVersion: "1.0.0",
		AutoCreate:     true,
		BOM:            base64.StdEncoding.EncodeToString(bomContent),
	})
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		timeout := time.After(30 * time.Second)

	loop:
		for {
			select {
			case <-ticker.C:
				processing, err := client.BOM.IsBeingProcessed(context.TODO(), uploadToken)
				if err != nil {
					panic(err)
				}
				if !processing {
					break loop
				}
			case <-timeout:
				panic("timeout exceeded")
			}
		}
	}()

	wg.Wait()
}
