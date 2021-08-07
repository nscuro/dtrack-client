package main

import (
	"context"
	"flag"
	"fmt"
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
	findings := make([]dtrack.Finding, 0)

	var (
		pageNumber = 1
		pageSize   = 10
	)

	for {
		fr, err := client.GetFindings(ctx, pu, false, dtrack.PageOptions{
			PageNumber: pageNumber,
			PageSize:   pageSize,
		})
		if err != nil {
			log.Fatalf("failed to fetch findings: %v", err)
		}

		findings = append(findings, fr.Findings...)
		log.Printf("fetched %d/%d findings", len(findings), fr.TotalCount)

		if len(findings) >= fr.TotalCount {
			break
		}

		pageNumber++
	}

	if len(findings) == 0 {
		return
	}

	fmt.Println()
	for _, finding := range findings {
		fmt.Printf(" > %s [%s]\n", finding.Vulnerability.VulnID, finding.Vulnerability.Severity)
		fmt.Printf("   Component: %s\n", finding.Component.PackageURL)
		if finding.Analysis != nil {
			fmt.Printf("   Analysis: state=%s, suppressed=%t\n", finding.Analysis.State, finding.Analysis.Suppressed)
		}
		fmt.Printf("   Details: %s\n", fmt.Sprintf("%s/vulnerabilities/%s/%s", baseURL, finding.Vulnerability.Source, finding.Vulnerability.VulnID))
		fmt.Println()
	}
}
