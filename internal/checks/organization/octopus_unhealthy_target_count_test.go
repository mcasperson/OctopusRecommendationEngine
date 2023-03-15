package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusTerraformTestFramework/octoclient"
	"github.com/mcasperson/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"
	"time"
)

func TestUnhealthyTargets(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform"), "27-unhealthytargets", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		// loop for a bit until the target is unhealthy
		for i := 0; i < 6; i++ {
			machines, err := newSpaceClient.Machines.GetAll()

			if err != nil {
				return err
			}

			if len(machines) > 0 && machines[0].HealthStatus != "Healthy" && machines[0].HealthStatus != "Unknown" {
				break
			}

			time.Sleep(time.Second * 10)
		}

		check := NewOctopusUnhealthyTargetCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

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
