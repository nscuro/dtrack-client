package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Component struct {
	UUID    uuid.UUID `json:"uuid"`
	Name    string    `json:"name"`
	Version string    `json:"version"`
	Group   string    `json:"group"`
}

type ComponentsResponse struct {
	Components []Component
	TotalCount int
}

func (c Client) GetComponents(ctx context.Context, puuid uuid.UUID, po PageOptions) (*ComponentsResponse, error) {
	req, err := c.newPagingRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/component/project/%s", puuid), nil, nil, po)
	if err != nil {
		return nil, err
	}

	var components []Component
	res, err := c.doRequest(req, &components)
	if err != nil {
		return nil, err
	}

	return &ComponentsResponse{
		Components: components,
		TotalCount: res.TotalCount,
	}, nil
}
