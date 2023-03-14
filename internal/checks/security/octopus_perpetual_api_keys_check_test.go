package security

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusTerraformTestFramework/octoclient"
	"github.com/mcasperson/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"
)

func TestPerpetualApiKeys(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		dir := filepath.Join("..", "..", "..", "test", "terraform")
		newSpaceId, err := testFramework.Act(t, container, dir, "1-singlespace", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		allUsers, err := newSpaceClient.Users.GetAll()

		if err != nil {
			return err
		}

		newSpaceClient.APIKeys.Create(&users.APIKey{
			APIKey:  "",
			Created: nil,
			Purpose: "",
			UserID:  allUsers[0].ID,
			Expires: nil,
		})

		check := NewOctopusPerpetualApiKeysCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute()

		if err != nil {
			return err
		}

		// Assert
		if result.Severity() != checks.Warning {
			return errors.New("Check should have returned a warning")
		}

		return nil
	})
}
