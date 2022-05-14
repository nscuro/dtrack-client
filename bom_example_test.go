package dtrack_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/nscuro/dtrack-client"
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

	var (
		doneChan = make(chan struct{})
		errChan  = make(chan error)
		ticker   = time.NewTicker(1 * time.Second)
		timeout  = time.After(30 * time.Second)
	)

	go func() {
		defer func() {
			close(doneChan)
			close(errChan)
		}()

		for {
			select {
			case <-ticker.C:
				processing, err := client.BOM.IsBeingProcessed(context.TODO(), uploadToken)
				if err != nil {
					errChan <- err
					return
				}
				if !processing {
					doneChan <- struct{}{}
					return
				}
			case <-timeout:
				errChan <- fmt.Errorf("timeout exceeded")
				return
			}
		}
	}()

	select {
	case <-doneChan:
		fmt.Println("bom processing completed")
	case <-errChan:
		fmt.Printf("failed to wait for bom processing: %v\n", err)
	}
}
