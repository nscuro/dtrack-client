package dtrack_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/nscuro/dtrack-client"
)

// This example demonstrates how to fetch all findings for a given project.
// It also shows how to navigate through Dependency-Track's API paging.
func Example_getAllFindings() {
	client, _ := dtrack.NewClient("https://dtrack.example.com", dtrack.WithAPIKey("..."))
	projectUUID := uuid.MustParse("2d16089e-6d3a-437e-b334-f27eb2cbd7f4")

	var (
		findings   []dtrack.Finding
		pageNumber = 1
		pageSize   = 10
	)

	for {
		findingsPage, err := client.Finding.GetAll(context.TODO(), projectUUID, false, dtrack.PageOptions{
			PageNumber: pageNumber,
			PageSize:   pageSize,
		})
		if err != nil {
			panic(err)
		}

		findings = append(findings, findingsPage.Findings...)

		if len(findings) >= findingsPage.TotalCount {
			break
		}

		pageNumber++
	}
}
