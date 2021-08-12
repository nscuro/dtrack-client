package dtrack

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type License struct {
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	Text     string    `json:"text"`
	Template string    `json:"template"`
	Header   string    `json:"header"`
	Comment  string    `json:"comment"`

	LicenseID           string   `json:"licenseId"`
	OSIApproved         bool     `json:"isOsiApproved"`
	FSFLibre            bool     `json:"isFsfLibre"`
	DeprecatedLicenseID bool     `json:"isDeprecatedLicenseId"`
	SeeAlso             []string `json:"seeAlso"`
}

type LicensePage struct {
	Licenses   []License
	TotalCount int
}

type LicenseService struct {
	client *Client
}

func (l LicenseService) GetAll(ctx context.Context, po PageOptions) (*LicensePage, error) {
	req, err := l.client.newRequest(ctx, http.MethodGet, "/api/v1/license", withPageOptions(po))
	if err != nil {
		return nil, err
	}

	var licenses []License
	res, err := l.client.doRequest(req, &licenses)
	if err != nil {
		return nil, err
	}

	return &LicensePage{
		Licenses:   licenses,
		TotalCount: res.TotalCount,
	}, nil
}
