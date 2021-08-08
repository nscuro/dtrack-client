package dtrack

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestAboutService_Get(t *testing.T) {
	client, err := NewClient("http://localhost")
	require.NoError(t, err)

	httpmock.ActivateNonDefault(client.httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://localhost/api/version",
		httpmock.NewStringResponder(http.StatusOK, `{
	"timestamp": "2020-09-29T22:02:23Z",
	"version": "4.0.0-SNAPSHOT",
	"uuid": "c35ce882-398d-46ed-8c36-6148bf73f941",
	"systemUuid": "f2eee6f6-a161-418a-baf5-b57a2e30de82",
	"application": "Dependency-Track",
	"framework": {
		"timestamp": "2020-07-20T15:56:44Z",
		"version": "1.8.0-SNAPSHOT",
		"uuid": "beee8786-ca9c-473a-b7a5-efcc95e8c469",
		"name": "Alpine"
	}
}`))

	about, err := client.About.Get(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, about)

	require.Equal(t, "2020-09-29T22:02:23Z", about.Timestamp)
	require.Equal(t, "4.0.0-SNAPSHOT", about.Version)
	require.Equal(t, "c35ce882-398d-46ed-8c36-6148bf73f941", about.UUID.String())
	require.Equal(t, "f2eee6f6-a161-418a-baf5-b57a2e30de82", about.SystemUUID.String())
	require.Equal(t, "Dependency-Track", about.Application)

	require.Equal(t, "2020-07-20T15:56:44Z", about.Framework.Timestamp)
	require.Equal(t, "1.8.0-SNAPSHOT", about.Framework.Version)
	require.Equal(t, "beee8786-ca9c-473a-b7a5-efcc95e8c469", about.Framework.UUID.String())
	require.Equal(t, "Alpine", about.Framework.Name)
}
