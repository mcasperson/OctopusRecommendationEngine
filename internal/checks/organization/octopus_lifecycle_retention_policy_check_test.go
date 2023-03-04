package organization

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks/test"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/octoclient"
	"path/filepath"
	"testing"
)

func TestLifecyclesMeetRecommendations(t *testing.T) {
	test.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := test.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform", "16-lifecyclesmeetrecommendations"), []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusLifecycleRetentionPolicyCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

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

func TestLifecycleKeepsReleasesForever(t *testing.T) {
	test.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := test.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform", "17-lifecyclekeepsreleasesforever"), []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusLifecycleRetentionPolicyCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

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

func TestLifecycleKeepsFilesForever(t *testing.T) {
	test.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := test.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform", "18-lifecyclekeepsfilesforever"), []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusLifecycleRetentionPolicyCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

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

func TestLifecyclePhaseKeepsReleasesForever(t *testing.T) {
	test.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := test.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform", "19-lifecyclephasekeepsreleasesforever"), []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusLifecycleRetentionPolicyCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

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

func TestLifecyclePhaseKeepsFilesForever(t *testing.T) {
	test.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := test.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform", "20-lifecyclephasekeepsfilesforever"), []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusLifecycleRetentionPolicyCheck(newSpaceClient, checks.OctopusClientPermissiveErrorHandler{})

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
