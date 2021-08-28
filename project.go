package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Project struct {
	UUID        uuid.UUID `json:"uuid,omitempty"`
	Author      string    `json:"author,omitempty"`
	Publisher   string    `json:"publisher,omitempty"`
	Group       string    `json:"group,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Version     string    `json:"version,omitempty"`
	Classifier  string    `json:"classifier,omitempty"`

	CPE       string `json:"cpe,omitempty"`
	PURL      string `json:"purl,omitempty"`
	SWIDTagID string `json:"swidTagId,omitempty"`

	DirectDependencies string            `json:"directDependencies,omitempty"`
	Properties         []ProjectProperty `json:"properties,omitempty"`
	Tags               []Tag             `json:"tags,omitempty"`
	Active             bool              `json:"active"`
}

type ProjectsPage struct {
	Projects   []Project
	TotalCount int
}

type ProjectService struct {
	client *Client
}

func (p ProjectService) Get(ctx context.Context, projectUUID uuid.UUID) (*Project, error) {
	req, err := p.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/project/%s", projectUUID))
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

func (p ProjectService) Create(ctx context.Context, project Project) (*Project, error) {
	req, err := p.client.newRequest(ctx, http.MethodPut, "/api/v1/project", withBody(project))
	if err != nil {
		return nil, err
	}

	var createdProject Project
	_, err = p.client.doRequest(req, &createdProject)
	if err != nil {
		return nil, err
	}

	return &createdProject, nil
}

func (p ProjectService) Update(ctx context.Context, project Project) (*Project, error) {
	req, err := p.client.newRequest(ctx, http.MethodPost, "/api/v1/project", withBody(project))
	if err != nil {
		return nil, err
	}

	var createdProject Project
	_, err = p.client.doRequest(req, &createdProject)
	if err != nil {
		return nil, err
	}

	return &createdProject, nil
}

func (p ProjectService) Delete(ctx context.Context, projectUUID uuid.UUID) error {
	req, err := p.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/project/%s", projectUUID))
	if err != nil {
		return err
	}

	_, err = p.client.doRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
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

type ProjectCloneRequest struct {
	ProjectUUID         uuid.UUID `json:"project"`
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
