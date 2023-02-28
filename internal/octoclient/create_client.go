package octoclient

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"net/url"
)

// CreateClient creates a Octopus octoclient to the given url
func CreateClient(uri string, spaceId string, apiKey string) (*client.Client, error) {
	url, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	return client.NewClient(nil, url, apiKey, spaceId)
}
