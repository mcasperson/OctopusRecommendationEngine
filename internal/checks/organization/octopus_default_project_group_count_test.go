package organization

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks/test"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/octoclient"
	"path/filepath"
	"testing"
)

func TestNormalProjectCount(t *testing.T) {
	test.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := test.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform", "4-smallprojectcount"), []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusDefaultProjectGroupCountCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute()

		if err != nil {
			return err
		}

		// Assert
		if result.Severity() != checks.Ok {
			t.Fatal("Check should have passed")
		}

		return nil
	})
}

func TestExcessiveProjectCount(t *testing.T) {
	test.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := test.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform", "5-largeprojectcount"), []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusDefaultProjectGroupCountCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute()

		if err != nil {
			return err
		}

		// Assert
		if result.Severity() != checks.Warning {
			t.Fatal("Check should have produced a warning")
		}

		return nil
	})
}

func TestExcessiveProjectCountWithPermissionsError(t *testing.T) {
	test.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		dir := filepath.Join("..", "..", "..", "test", "terraform", "14-limitedserviceaccount")
		newSpaceId, err := test.Act(t, container, dir, []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		serviceAccountId, err := test.GetOutputVariable(t, dir, "service_account_id")

		if err != nil {
			return err
		}

		apiKey, err := newSpaceClient.APIKeys.Create(users.NewAPIKey("Limited API Key", serviceAccountId))

		if err != nil {
			return err
		}

		limitedClient, err := octoclient.CreateClient(container.URI, newSpaceId, apiKey.APIKey)

		if err != nil {
			return err
		}

		check := NewOctopusDefaultProjectGroupCountCheck(limitedClient, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute()

		if err != nil {
			return err
		}

		// Assert
		if result.Severity() != checks.Permission {
			t.Fatal("Check should have produced a permission warning")
		}

		return nil
	})
}
