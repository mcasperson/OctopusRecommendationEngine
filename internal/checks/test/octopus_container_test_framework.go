package test

import (
	"context"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/octoclient"
	lintwait "github.com/mcasperson/OctopusRecommendationEngine/internal/wait"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

/*
	This file contains a bunch of functions to support integration tests with a live Octopus instance hosted
	in a Docker container and managed by test containers.
*/

const ApiKey = "API-ABCDEFGHIJKLMNOPQURTUVWXYZ12345"

type OctopusContainer struct {
	testcontainers.Container
	URI string
}

type mysqlContainer struct {
	testcontainers.Container
	port string
	ip   string
}

type TestLogConsumer struct {
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	fmt.Println(string(l.Content))
}

type OctopusContainerTest struct {
}

func (o *OctopusContainerTest) enableContainerLogging(container testcontainers.Container, ctx context.Context) error {
	// Display the container logs
	err := container.StartLogProducer(ctx)
	if err != nil {
		return err
	}
	g := TestLogConsumer{}
	container.FollowOutput(&g)
	return nil
}

// getReaperSkipped will return true if running in a podman environment
func (o *OctopusContainerTest) getReaperSkipped() bool {
	if strings.Contains(os.Getenv("DOCKER_HOST"), "podman") {
		return true
	}

	return false
}

// getProvider returns the test containers provider
func (o *OctopusContainerTest) getProvider() testcontainers.ProviderType {
	if strings.Contains(os.Getenv("DOCKER_HOST"), "podman") {
		return testcontainers.ProviderPodman
	}

	return testcontainers.ProviderDocker
}

// setupNetwork creates an internal network for Octopus and MS SQL
func (o *OctopusContainerTest) setupNetwork(ctx context.Context) (testcontainers.Network, error) {
	return testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:           "octopusterraformtests",
			CheckDuplicate: false,
			SkipReaper:     o.getReaperSkipped(),
		},
		ProviderType: o.getProvider(),
	})
}

// setupDatabase creates a MSSQL container
func (o *OctopusContainerTest) setupDatabase(ctx context.Context) (*mysqlContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mcr.microsoft.com/mssql/server",
		ExposedPorts: []string{"1433/tcp"},
		Env: map[string]string{
			"ACCEPT_EULA": "Y",
			"SA_PASSWORD": "Password01!",
		},
		WaitingFor: wait.ForExec([]string{"/opt/mssql-tools/bin/sqlcmd", "-U", "sa", "-P", "Password01!", "-Q", "select 1"}).WithExitCodeMatcher(
			func(exitCode int) bool {
				return exitCode == 0
			}),
		SkipReaper: o.getReaperSkipped(),
		Networks: []string{
			"octopusterraformtests",
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "1433")
	if err != nil {
		return nil, err
	}

	return &mysqlContainer{
		Container: container,
		ip:        ip,
		port:      mappedPort.Port(),
	}, nil
}

// setupOctopus creates an Octopus container
func (o *OctopusContainerTest) setupOctopus(ctx context.Context, connString string) (*OctopusContainer, error) {
	if os.Getenv("LICENSE") == "" {
		return nil, errors.New("the LICENSE environment variable must be set to a base 64 encoded Octopus license key")
	}

	if _, err := b64.StdEncoding.DecodeString(os.Getenv("LICENSE")); err != nil {
		return nil, errors.New("the LICENSE environment variable must be set to a base 64 encoded Octopus license key")
	}

	req := testcontainers.ContainerRequest{
		Image:        "octopusdeploy/octopusdeploy:latest",
		ExposedPorts: []string{"8080/tcp"},
		Env: map[string]string{
			"ACCEPT_EULA":                   "Y",
			"DB_CONNECTION_STRING":          connString,
			"ADMIN_API_KEY":                 ApiKey,
			"DISABLE_DIND":                  "Y",
			"ADMIN_USERNAME":                "admin",
			"ADMIN_PASSWORD":                "Password01!",
			"OCTOPUS_SERVER_BASE64_LICENSE": os.Getenv("LICENSE"),
		},
		Privileged: false,
		WaitingFor: wait.ForLog("Listening for HTTP requests on").WithStartupTimeout(30 * time.Minute),
		SkipReaper: o.getReaperSkipped(),
		Networks: []string{
			"octopusterraformtests",
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// Display the container logs
	o.enableContainerLogging(container, ctx)

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "8080")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &OctopusContainer{Container: container, URI: uri}, nil
}

// ArrangeTest is wrapper that initialises Octopus, runs a test, and cleans up the containers
func (o *OctopusContainerTest) ArrangeTest(t *testing.T, testFunc func(t *testing.T, container *OctopusContainer, client *client.Client) error) {
	err := retry.Do(
		func() error {

			if testing.Short() {
				t.Skip("skipping integration test")
			}

			ctx := context.Background()

			network, err := o.setupNetwork(ctx)
			if err != nil {
				return err
			}

			sqlServer, err := o.setupDatabase(ctx)
			if err != nil {
				return err
			}

			sqlIp, err := sqlServer.Container.ContainerIP(ctx)
			if err != nil {
				return err
			}

			t.Log("SQL Server IP: " + sqlIp)

			octopusContainer, err := o.setupOctopus(ctx, "Server="+sqlIp+",1433;Database=OctopusDeploy;User=sa;Password=Password01!")
			if err != nil {
				return err
			}

			// Clean up the container after the test is complete
			defer func() {
				// This fixes the "can not get logs from container which is dead or marked for removal" error
				// See https://github.com/testcontainers/testcontainers-go/issues/606
				octopusContainer.StopLogProducer()

				octoTerminateErr := octopusContainer.Terminate(ctx)
				sqlTerminateErr := sqlServer.Terminate(ctx)

				networkErr := network.Remove(ctx)

				if octoTerminateErr != nil || sqlTerminateErr != nil || networkErr != nil {
					t.Fatalf("failed to terminate container: %v %v", octoTerminateErr, sqlTerminateErr)
				}
			}()

			// give the server 5 minutes to start up
			err = lintwait.WaitForResource(func() error {
				resp, err := http.Get(octopusContainer.URI + "/api")
				if err != nil || resp.StatusCode != http.StatusOK {
					return errors.New("the api endpoint was not available")
				}
				return nil
			}, 5*time.Minute)

			if err != nil {
				return err
			}

			client, err := octoclient.CreateClient(octopusContainer.URI, "", ApiKey)
			if err != nil {
				return err
			}

			return testFunc(t, octopusContainer, client)
		},
		retry.Attempts(3),
	)

	if err != nil {
		t.Fatalf(err.Error())
	}
}

// cleanTerraformModule removes state and lock files to ensure we get a clean run each time
func (o *OctopusContainerTest) cleanTerraformModule(terraformProjectDir string) error {
	os.Remove(filepath.Join(terraformProjectDir, ".terraform.lock.hcl"))
	os.Remove(filepath.Join(terraformProjectDir, "terraform.tfstate"))
	os.Remove(filepath.Join(terraformProjectDir, ".terraform.tfstate.lock.info"))
	return nil
}

// terraformInit runs "terraform init"
func (o *OctopusContainerTest) terraformInit(t *testing.T, terraformProjectDir string) error {
	args := []string{"init", "-no-color"}
	cmnd := exec.Command("terraform", args...)
	cmnd.Dir = terraformProjectDir
	out, err := cmnd.Output()

	t.Log(string(out))

	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			t.Log(string(exitError.Stderr))
		} else {
			t.Log(err.Error())
		}

		return err
	}

	return nil
}

// terraformApply runs "terraform apply"
func (o *OctopusContainerTest) terraformApply(t *testing.T, terraformProjectDir string, server string, spaceId string, vars []string) error {
	newArgs := append([]string{
		"apply",
		"-auto-approve",
		"-no-color",
		"-var=octopus_server=" + server,
		"-var=octopus_apikey=" + ApiKey,
		"-var=octopus_space_id=" + spaceId,
	}, vars...)

	cmnd := exec.Command("terraform", newArgs...)
	cmnd.Dir = terraformProjectDir
	out, err := cmnd.Output()

	t.Log(string(out))

	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			t.Log(string(exitError.Stderr))
		} else {
			t.Log(err)
		}
		return err
	}

	return nil
}

// waitForSpace attempts to ensure the API and space is available before continuing
func (o *OctopusContainerTest) waitForSpace(server string, spaceId string) {
	// Error like:
	// Error: Octopus API error: Resource is not found or it doesn't exist in the current space context. Please contact your administrator for more information. []
	// are sometimes proceeded with:
	// "HTTP" "GET" to "localhost:32805""/api" "completed" with 503 in 00:00:00.0170358 (17ms) by "<anonymous>"
	// So wait until we get a valid response from the API endpoint before applying terraform
	lintwait.WaitForResource(func() error {
		response, err := http.Get(server + "/api")
		if err != nil {
			return err
		}
		if !(response.StatusCode >= 200 && response.StatusCode <= 299) {
			return errors.New("non 2xx status code returned")
		}
		return nil
	}, time.Minute)

	// Also wait for the space to be available
	lintwait.WaitForResource(func() error {
		response, err := http.Get(server + "/api/" + spaceId)
		if err != nil {
			return err
		}
		if !(response.StatusCode >= 200 && response.StatusCode <= 299) {
			return errors.New("non 2xx status code returned")
		}
		return nil
	}, time.Minute)
}

// initialiseOctopus uses Terraform to populate the test Octopus instance, making sure to clean up
// any files generated during previous Terraform executions to avoid conflicts and locking issues.
func (o *OctopusContainerTest) initialiseOctopus(t *testing.T, container *OctopusContainer, terraformDir string, spaceName string, initialiseVars []string, populateVars []string) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	t.Log("Working dir: " + path)

	// This test creates a new space and then populates the space.
	terraformProjectDirs := []string{}
	terraformProjectDirs = append(terraformProjectDirs, filepath.Join("..", "..", "..", "test", "terraform", "1-singlespace"))
	terraformProjectDirs = append(terraformProjectDirs, terraformDir)

	// First loop initialises the new space, second populates the space
	spaceId := "Spaces-1"
	for i, terraformProjectDir := range terraformProjectDirs {

		o.cleanTerraformModule(terraformProjectDir)

		err := o.terraformInit(t, terraformProjectDir)

		if err != nil {
			return err
		}

		// when initialising the new space, we need to define a new space name as a variable
		vars := []string{}
		if i == 0 {
			vars = append(initialiseVars, "-var=octopus_space_name="+spaceName)
		} else {
			vars = populateVars
		}

		o.waitForSpace(container.URI, spaceId)

		err = o.terraformApply(t, terraformProjectDir, container.URI, spaceId, vars)

		if err != nil {
			return err
		}

		// get the ID of any new space created, which will be used in the subsequent Terraform executions
		spaceId, err = GetOutputVariable(t, terraformProjectDir, "octopus_space_id")

		if err != nil {
			return err
		}
	}

	return nil
}

// getOutputVariable reads a Terraform output variable
func GetOutputVariable(t *testing.T, terraformDir string, outputVar string) (string, error) {
	cmnd := exec.Command(
		"terraform",
		"output",
		"-raw",
		outputVar)
	cmnd.Dir = terraformDir
	out, err := cmnd.Output()

	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			t.Log(string(exitError.Stderr))
		} else {
			t.Log(err)
		}
		return "", err
	}

	return string(out), nil
}

// Act initialises Octopus and MSSQL
func (o *OctopusContainerTest) Act(t *testing.T, container *OctopusContainer, terraformDir string, populateVars []string) (string, error) {
	t.Log("POPULATING TEST SPACE")

	spaceName := strings.ReplaceAll(fmt.Sprint(uuid.New()), "-", "")[:20]
	err := o.initialiseOctopus(t, container, terraformDir, spaceName, []string{}, populateVars)

	if err != nil {
		return "", err
	}

	return GetOutputVariable(t, filepath.Join("..", "..", "..", "test", "terraform", "1-singlespace"), "octopus_space_id")
}
