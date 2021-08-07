package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Project struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ProjectsResponse struct {
	TotalCount int
	Projects   []Project
}

func (c Client) GetProject(ctx context.Context, u uuid.UUID) (*Project, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/project/%s", u))
	if err != nil {
		return nil, err
	}

	var project Project
	_, err = c.doRequest(req, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (c Client) LookupProject(ctx context.Context, name, version string) (*Project, error) {
	params := map[string]string{
		"name":    name,
		"version": version,
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/project/lookup", withParams(params))
	if err != nil {
		return nil, err
	}

	var project Project
	_, err = c.doRequest(req, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (c Client) GetProjects(ctx context.Context, po PageOptions) (*ProjectsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/project", withPageOptions(po))
	if err != nil {
		return nil, err
	}

	var projects []Project
	res, err := c.doRequest(req, &projects)
	if err != nil {
		return nil, err
	}

	return &ProjectsResponse{
		TotalCount: res.TotalCount,
		Projects:   projects,
	}, nil
}
