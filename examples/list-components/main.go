package main

import (
	"context"
	"flag"
	"log"

	"github.com/google/uuid"
	"github.com/nscuro/dtrack-client"
)

func main() {
	var (
		baseURL     string
		apiKey      string
		projectUUID string
	)

	flag.StringVar(&baseURL, "url", "", "Dependency-Track URL")
	flag.StringVar(&apiKey, "apikey", "", "Dependency-Track API key")
	flag.StringVar(&projectUUID, "project", "", "Project UUID")
	flag.Parse()

	client, err := dtrack.NewClient(baseURL, dtrack.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}

	pu, err := uuid.Parse(projectUUID)
	if err != nil {
		log.Fatalf("failed to parse project uuid: %v", err)
	}

	ctx := context.Background()
	components := make([]dtrack.Component, 0)

	var (
		pageNumber = 1
		pageSize   = 10
	)

	for {
		cr, err := client.GetComponents(ctx, pu, dtrack.PageOptions{
			PageNumber: pageNumber,
			PageSize:   pageSize,
		})
		if err != nil {
			log.Fatalf("failed to fetch components: %v", err)
		}

		components = append(components, cr.Components...)
		log.Printf("fetched %d/%d components", len(components), cr.TotalCount)

		if len(components) >= cr.TotalCount {
			break
		}

		pageNumber++
	}
}
