package organization

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks/test"
	"path/filepath"
	"testing"
)

func TestNormalEnvironmentCount(t *testing.T) {
	test.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := test.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform", "2-fewenvironments"), []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := test.CreateClient(container.URI, newSpaceId)

		if err != nil {
			return err
		}

		check := NewOctopusEnvironmentCountCheck(newSpaceClient)

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

func TestExcessiveEnvironmentCount(t *testing.T) {
	test.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := test.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform", "3-toomanyenvironments"), []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := test.CreateClient(container.URI, newSpaceId)

		if err != nil {
			return err
		}

		check := NewOctopusEnvironmentCountCheck(newSpaceClient)

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
