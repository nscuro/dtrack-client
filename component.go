package dtrack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Component struct {
	UUID       uuid.UUID `json:"uuid"`
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	Group      string    `json:"group"`
	PackageURL string    `json:"purl"`
}

type ComponentsPage struct {
	Components []Component
	TotalCount int
}

type ComponentService struct {
	client *Client
}

func (c ComponentService) GetAll(ctx context.Context, project uuid.UUID, po PageOptions) (*ComponentsPage, error) {
	req, err := c.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/component/project/%s", project), withPageOptions(po))
	if err != nil {
		return nil, err
	}

	var components []Component
	res, err := c.client.doRequest(req, &components)
	if err != nil {
		return nil, err
	}

	return &ComponentsPage{
		Components: components,
		TotalCount: res.TotalCount,
	}, nil
}
