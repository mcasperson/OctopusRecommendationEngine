package security

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusTerraformTestFramework/octoclient"
	"github.com/mcasperson/OctopusTerraformTestFramework/test"
	"os"
	"path/filepath"
	"testing"
)

func TestDuplicatedGitCreds(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		dir := filepath.Join("..", "..", "..", "test", "terraform")
		newSpaceId, err := testFramework.Act(t, container, dir, "26-duplicatedgitusernames", []string{
			"-var=testgitrepo=" + os.Getenv("TEST_GIT_REPO"),
			"-var=testgitusername=" + os.Getenv("TEST_GIT_USERNAME"),
			"-var=testgitpassword=" + os.Getenv("TEST_GIT_PASSWORD"),
		})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusDuplicatedGitCredentialsCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

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
