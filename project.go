package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Project struct {
	UUID    uuid.UUID `json:"uuid"`
	Name    string    `json:"name"`
	Version string    `json:"version"`
}

type ProjectsPage struct {
	Projects   []Project
	TotalCount int
}

type ProjectService struct {
	client *Client
}

func (p ProjectService) Get(ctx context.Context, u uuid.UUID) (*Project, error) {
	req, err := p.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/project/%s", u))
	if err != nil {
		return nil, err
	}

	var project Project
	_, err = p.client.doRequest(req, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (p ProjectService) Lookup(ctx context.Context, name, version string) (*Project, error) {
	params := map[string]string{
		"name":    name,
		"version": version,
	}

	req, err := p.client.newRequest(ctx, http.MethodGet, "/api/v1/project/lookup", withParams(params))
	if err != nil {
		return nil, err
	}

	var project Project
	_, err = p.client.doRequest(req, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (p ProjectService) GetAll(ctx context.Context, po PageOptions) (*ProjectsPage, error) {
	req, err := p.client.newRequest(ctx, http.MethodGet, "/api/v1/project", withPageOptions(po))
	if err != nil {
		return nil, err
	}

	var projects []Project
	res, err := p.client.doRequest(req, &projects)
	if err != nil {
		return nil, err
	}

	return &ProjectsPage{
		TotalCount: res.TotalCount,
		Projects:   projects,
	}, nil
}
