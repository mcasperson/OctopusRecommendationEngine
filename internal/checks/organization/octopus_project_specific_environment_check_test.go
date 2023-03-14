package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusTerraformTestFramework/octoclient"
	"github.com/mcasperson/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"
)

func TestProjectSpecificEnvironmentsFromDefaultLifecycle(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform"), "22-singleprojectenvironments", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusProjectSpecificEnvironmentCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute()

		if err != nil {
			return err
		}

		// Assert
		if result.Severity() != checks.Warning {
			return errors.New("Check should have failed")
		}

		return nil
	})
}

func TestProjectSpecificEnvironmentsFromCustomLifecycle(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform"), "23-singleprojectenvironmentswithcustomlifecycles", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusProjectSpecificEnvironmentCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute()

		if err != nil {
			return err
		}

		// Assert
		if result.Severity() != checks.Warning {
			return errors.New("Check should have failed")
		}

		return nil
	})
}
