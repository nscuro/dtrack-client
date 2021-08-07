package dtrack

import (
	"context"
	"net/http"
)

type RepositoryMetaComponent struct {
	LatestVersion string `json:"latestVersion"`
}

func (c Client) GetRepositoryMetaComponent(ctx context.Context, purl string) (*RepositoryMetaComponent, error) {
	params := map[string]string{
		"purl": purl,
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/repository/latest", withParams(params))
	if err != nil {
		return nil, err
	}

	var meta RepositoryMetaComponent
	_, err = c.doRequest(req, &meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}
