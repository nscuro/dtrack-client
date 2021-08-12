package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Project struct {
	UUID        uuid.UUID `json:"uuid"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	Group       string    `json:"group"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	Classifier  string    `json:"classifier"`

	CPE       string `json:"cpe"`
	PURL      string `json:"purl"`
	SWIDTagID string `json:"swidTagId"`

	DirectDependencies string            `json:"directDependencies"`
	Properties         []ProjectProperty `json:"properties"`
	Tags               []Tag             `json:"tags"`
	Active             bool              `json:"active"`
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

type ProjectCloneRequest struct {
	UUID                uuid.UUID `json:"project"`
	Version             string    `json:"version"`
	IncludeAuditHistory bool      `json:"includeAuditHistory"`
	IncludeComponents   bool      `json:"includeComponents"`
	IncludeProperties   bool      `json:"includeProperties"`
	IncludeServices     bool      `json:"includeServices"`
	IncludeTags         bool      `json:"includeTags"`
}

func (p ProjectService) Clone(ctx context.Context, cloneReq ProjectCloneRequest) error {
	req, err := p.client.newRequest(ctx, http.MethodPut, "/api/v1/project/clone", withBody(cloneReq))
	if err != nil {
		return err
	}

	_, err = p.client.doRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}
