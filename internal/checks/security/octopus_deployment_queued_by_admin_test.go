package security

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/releases"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusTerraformTestFramework/octoclient"
	"github.com/mcasperson/OctopusTerraformTestFramework/test"
	"github.com/mcasperson/OctopusTerraformTestFramework/wait"
	"path/filepath"
	"testing"
	"time"
)

func TestDeployedByAdmin(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		dir := filepath.Join("..", "..", "..", "test", "terraform")
		newSpaceId, err := testFramework.Act(t, container, dir, "15-deployedbyadmin", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		projectId, err := test.GetOutputVariable(t, filepath.Join(dir, "15-deployedbyadmin"), "project_id")

		if err != nil {
			return err
		}

		channels, err := newSpaceClient.Channels.Get(channels.Query{
			IDs:         nil,
			PartialName: "Default",
			Skip:        0,
			Take:        1,
		})

		if err != nil {
			return err
		}

		release, err := newSpaceClient.Releases.Add(&releases.Release{
			Assembled:                          time.Time{},
			BuildInformation:                   nil,
			ChannelID:                          channels.Items[0].ID,
			IgnoreChannelRules:                 false,
			LibraryVariableSetSnapshotIDs:      nil,
			ProjectDeploymentProcessSnapshotID: "",
			ProjectID:                          projectId,
			ProjectVariableSetSnapshotID:       "",
			ReleaseNotes:                       "",
			SelectedPackages:                   nil,
			SpaceID:                            "",
			Version:                            "0.0.1",
		})

		if err != nil {
			return err
		}

		environment, err := newSpaceClient.Environments.Get(environments.EnvironmentsQuery{
			IDs:         nil,
			PartialName: "Development",
			Skip:        0,
			Take:        1,
		})

		if err != nil {
			return err
		}

		deployment, err := newSpaceClient.Deployments.Add(&deployments.Deployment{
			Changes:                  nil,
			ChangesMarkdown:          "",
			ChannelID:                "",
			Comments:                 "",
			Created:                  nil,
			DeployedBy:               "",
			DeployedByID:             "",
			DeployedToMachineIDs:     nil,
			DeploymentProcessID:      "",
			EnvironmentID:            environment.Items[0].ID,
			ExcludedMachineIDs:       nil,
			FailureEncountered:       false,
			ForcePackageDownload:     false,
			ForcePackageRedeployment: false,
			FormValues:               nil,
			ManifestVariableSetID:    "",
			Name:                     "",
			ProjectID:                projectId,
			QueueTime:                nil,
			QueueTimeExpiry:          nil,
			ReleaseID:                release.ID,
			SkipActions:              nil,
			SpaceID:                  "",
			SpecificMachineIDs:       nil,
			TaskID:                   "",
			TenantID:                 "",
			TentacleRetentionPeriod:  nil,
			UseGuidedFailure:         false,
		})

		if err != nil {
			return err
		}

		err = wait.WaitForResource(func() error {
			_, err := newSpaceClient.Deployments.GetByID(deployment.ID)
			return err
		}, time.Minute)

		if err != nil {
			return err
		}

		check := NewOctopusDeploymentQueuedByAdminCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

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
