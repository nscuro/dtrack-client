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

	project, err := client.Project.Get(context.Background(), uuid.MustParse(projectUUID))
	if err != nil {
		log.Fatalf("failed to fetch project: %v", err)
	}

	findings, err := getFindingsForProject(client, project)
	if err != nil {
		log.Fatalf("failed to fetch findings: %v", err)
	}

	if len(findings) == 0 {
		return
	}

	fmt.Println()
	for _, finding := range findings {
		analysis := finding.Analysis
		vulnerability := finding.Vulnerability

		fmt.Printf(" > %s [%s]\n", vulnerability.VulnID, vulnerability.Severity)
		fmt.Printf("   Component: %s\n", finding.Component.PackageURL)
		fmt.Printf("   Analysis: state=%s, suppressed=%t\n", analysis.State, analysis.Suppressed)
		fmt.Printf("   Details: %s\n", fmt.Sprintf("%s/vulnerabilities/%s/%s", baseURL, vulnerability.Source, vulnerability.VulnID))
		fmt.Println()
	}
}

func getFindingsForProject(client *dtrack.Client, project *dtrack.Project) ([]dtrack.Finding, error) {
	log.Printf("fetching findings for project %s %s", project.Name, project.Version)

	var (
		findings   []dtrack.Finding
		pageNumber = 1
		pageSize   = 10
	)

	for {
		fr, err := client.Finding.GetAll(context.Background(), project.UUID, false, dtrack.PageOptions{
			PageNumber: pageNumber,
			PageSize:   pageSize,
		})
		if err != nil {
			return nil, err
		}

		findings = append(findings, fr.Findings...)
		log.Printf("fetched %d/%d findings", len(findings), fr.TotalCount)

		if len(findings) >= fr.TotalCount {
			break
		}

		pageNumber++
	}

	return findings, nil
}
