package dtrack

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type ProjectProperty struct {
	Group       string `json:"groupName"`
	Name        string `json:"propertyName"`
	Value       string `json:"propertyValue"`
	Type        string `json:"propertyType"`
	Description string `json:"description"`
}

type ProjectPropertiesPage struct {
	Properties []ProjectProperty
	TotalCount int
}

type ProjectPropertyService struct {
	client *Client
}

func (p ProjectPropertyService) GetAll(ctx context.Context, projectUUID uuid.UUID, po PageOptions) (*ProjectPropertiesPage, error) {
	req, err := p.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/project/%s/property", projectUUID), withPageOptions(po))
	if err != nil {
		return nil, err
	}

	var properties []ProjectProperty
	res, err := p.client.doRequest(req, &properties)
	if err != nil {
		return nil, err
	}

	return &ProjectPropertiesPage{
		Properties: properties,
		TotalCount: res.TotalCount,
	}, nil
}

func (p ProjectPropertyService) Create(ctx context.Context, projectUUID uuid.UUID, property ProjectProperty) (*ProjectProperty, error) {
	req, err := p.client.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/project/%s/property", projectUUID), withBody(property))
	if err != nil {
		return nil, err
	}

	var createdProperty ProjectProperty
	_, err = p.client.doRequest(req, &createdProperty)
	if err != nil {
		return nil, err
	}

	return &createdProperty, nil
}

func (p ProjectPropertyService) Update(ctx context.Context, projectUUID uuid.UUID, property ProjectProperty) (*ProjectProperty, error) {
	req, err := p.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/project/%s/property", projectUUID), withBody(property))
	if err != nil {
		return nil, err
	}

	var updatedProperty ProjectProperty
	_, err = p.client.doRequest(req, &updatedProperty)
	if err != nil {
		return nil, err
	}

	return &updatedProperty, nil
}

func (p ProjectPropertyService) Delete(ctx context.Context, projectUUID uuid.UUID, groupName, propertyName string) error {
	property := ProjectProperty{
		Group: groupName,
		Name:  propertyName,
	}

	req, err := p.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/project/%s/property", projectUUID), withBody(property))
	if err != nil {
		return err
	}

	_, err = p.client.doRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}
