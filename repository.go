package dtrack

import (
	"context"
	"net/http"
)

type RepositoryMetaComponent struct {
	LastCheck      int64  `json:"lastCheck"`
	LatestVersion  string `json:"latestVersion"`
	Name           string `json:"name"`
	Published      int64  `json:"published"`
	RepositoryType string `json:"repositoryType"`
}

type RepositoryService struct {
	client *Client
}

func (r RepositoryService) GetMetaComponent(ctx context.Context, purl string) (*RepositoryMetaComponent, error) {
	params := map[string]string{
		"purl": purl,
	}

	req, err := r.client.newRequest(ctx, http.MethodGet, "/api/v1/repository/latest", withParams(params))
	if err != nil {
		return nil, err
	}

	var meta RepositoryMetaComponent
	_, err = r.client.doRequest(req, &meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}
