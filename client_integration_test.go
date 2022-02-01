package dtrack_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/nscuro/dtrack-client"
)

func TestIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}

	suite.Run(t, &IntegrationTestSuite{})
}

type IntegrationTestSuite struct {
	suite.Suite
	container testcontainers.Container
	client    *dtrack.Client
}

func (s *IntegrationTestSuite) SetupSuite() {
	container, baseURL, err := s.setupContainer()
	s.Require().NoError(err)

	client, err := s.setupClient(baseURL)
	s.Require().NoError(err)

	s.container = container
	s.client = client
}

func (s *IntegrationTestSuite) TearDownSuite() {
	if s.container != nil {
		s.Require().NoError(s.container.Terminate(context.TODO()))
	}
}

func (s IntegrationTestSuite) TestAboutGet() {
	about, err := s.client.About.Get(context.TODO())
	s.Require().NoError(err)

	s.Require().Equal("Dependency-Track", about.Application)
	s.Require().Equal("Alpine", about.Framework.Name)
}

func (s IntegrationTestSuite) TestAPIError() {
	_, err := s.client.Project.Create(context.TODO(), dtrack.Project{})
	s.Require().Error(err)

	var apiErr *dtrack.APIError
	s.Require().ErrorAs(err, &apiErr)
	s.Require().Equal(http.StatusBadRequest, apiErr.StatusCode)
	s.Require().NotEmpty(apiErr.Message)
}

func (s IntegrationTestSuite) TestProjectCreateDelete() {
	project, err := s.client.Project.Create(context.TODO(), dtrack.Project{
		Name:    "TestProjectCreateDelete",
		Version: "1.0.0",
	})
	s.Require().NoError(err)
	s.Require().NotNil(project)
	s.Require().NotEmpty(project.UUID.String())
	s.Require().Equal("TestProjectCreateDelete", project.Name)
	s.Require().Equal("1.0.0", project.Version)

	err = s.client.Project.Delete(context.TODO(), project.UUID)
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) TestProjectPropertyCreateUpdateDelete() {
	project, err := s.client.Project.Create(context.TODO(), dtrack.Project{
		Name:    "TestProjectPropertyCreateUpdateDelete",
		Version: "1.0.0",
	})
	s.Require().NoError(err)

	property, err := s.client.ProjectProperty.Create(context.TODO(), project.UUID, dtrack.ProjectProperty{
		Group: "TestProjectPropertyCreateUpdateDelete",
		Name:  "TestName",
		Value: "TestValue",
		Type:  "STRING",
	})
	s.Require().NoError(err)
	s.Require().NotNil(property)

	property.Value = "TestValueUpdated"

	property, err = s.client.ProjectProperty.Update(context.TODO(), project.UUID, *property)
	s.Require().NoError(err)
	s.Require().NotNil(property)
	s.Require().Equal("TestValueUpdated", property.Value)

	err = s.client.ProjectProperty.Delete(context.TODO(), project.UUID, property.Group, property.Name)
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) setupContainer() (testcontainers.Container, string, error) {
	ctx, _ := context.WithTimeout(context.TODO(), 2*time.Minute) //nolint:govet

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "dependencytrack/apiserver:4.3.6",
			ExposedPorts: []string{"8080/tcp"},
			WaitingFor:   wait.ForLog("Dependency-Track is ready"),
			AutoRemove:   true,
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "8080")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get container port: %w", err)
	}

	return container, fmt.Sprintf("http://%s:%s", host, port), nil
}

func (s IntegrationTestSuite) setupClient(baseURL string) (*dtrack.Client, error) {
	client, err := dtrack.NewClient(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client: %w", err)
	}

	err = client.User.ForceChangePassword(context.TODO(), "admin", "admin", "nimda")
	if err != nil {
		return nil, fmt.Errorf("failed to change password: %w", err)
	}

	token, err := client.User.Login(context.TODO(), "admin", "nimda")
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	client, err = dtrack.NewClient(baseURL, dtrack.WithBearerToken(token))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize authenticated client: %w", err)
	}

	return client, nil
}
