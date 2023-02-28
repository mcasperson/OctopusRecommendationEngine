package organization

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks/test"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/octoclient"
	"path/filepath"
	"testing"
)

func TestPopulatedDeploymentProcess(t *testing.T) {
	test.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := test.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform", "7-populateddeploymentprocess"), []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusEmptyProjectCheck(newSpaceClient)

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

func TestEmptyDeploymentProcess(t *testing.T) {
	test.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := test.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform", "6-emptydeploymentprocess"), []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusEmptyProjectCheck(newSpaceClient)

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
