package octoclient

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"net/url"
	"strings"
)

// CreateClient creates a Octopus octoclient to the given url
func CreateClient(uri string, spaceId string, apiKey string) (*client.Client, error) {
	url, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	return client.NewClient(nil, url, apiKey, spaceId)
}

func ErrorIsPermissionError(err error) bool {
	apiError, ok := err.(*core.APIError)
	if ok {
		return strings.Index(apiError.ErrorMessage, "You do not have permission") != -1
	}
	return true
}

func ReturnPermissionResultOrError(id string, err error) (checks.OctopusCheckResult, error) {
	if ErrorIsPermissionError(err) {
		return checks.NewOctopusCheckResultImpl(
			"You do not have permission to run the check: "+err.Error(),
			id,
			"",
			checks.Permission,
			checks.Organization), nil
	}
	return nil, err
}
