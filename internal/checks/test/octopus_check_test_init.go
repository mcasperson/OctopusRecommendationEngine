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

func enableContainerLogging(container testcontainers.Container, ctx context.Context) error {
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
func getReaperSkipped() bool {
	if strings.Contains(os.Getenv("DOCKER_HOST"), "podman") {
		return true
	}

	return false
}

// getProvider returns the test containers provider
func getProvider() testcontainers.ProviderType {
	if strings.Contains(os.Getenv("DOCKER_HOST"), "podman") {
		return testcontainers.ProviderPodman
	}

	return testcontainers.ProviderDocker
}

// setupNetwork creates an internal network for Octopus and MS SQL
func setupNetwork(ctx context.Context) (testcontainers.Network, error) {
	return testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:           "octopusterraformtests",
			CheckDuplicate: false,
			SkipReaper:     getReaperSkipped(),
		},
		ProviderType: getProvider(),
	})
}

// setupDatabase creates a MSSQL container
func setupDatabase(ctx context.Context) (*mysqlContainer, error) {
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
		SkipReaper: getReaperSkipped(),
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
func setupOctopus(ctx context.Context, connString string) (*OctopusContainer, error) {
	if os.Getenv("LICENSE") == "" {
		return nil, errors.New("the LICENSE environment variable must be set to a base 64 encoded Octopus license key")
	}

	if _, err := b64.StdEncoding.DecodeString(os.Getenv("LICENSE")); err != nil {
		return nil, errors.New("the LICENSE environment variable must be set to a base 64 encoded Octopus license key")
	}

	req := testcontainers.ContainerRequest{
		// Be aware that later versions of Octopus killed Github Actions.
		// I think maybe they used more memory? 2022.2 works fine though.
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
		SkipReaper: getReaperSkipped(),
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
	enableContainerLogging(container, ctx)

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
func ArrangeTest(t *testing.T, testFunc func(t *testing.T, container *OctopusContainer, client *client.Client) error) {
	err := retry.Do(
		func() error {

			if testing.Short() {
				t.Skip("skipping integration test")
			}

			ctx := context.Background()

			network, err := setupNetwork(ctx)
			if err != nil {
				return err
			}

			sqlServer, err := setupDatabase(ctx)
			if err != nil {
				return err
			}

			sqlIp, err := sqlServer.Container.ContainerIP(ctx)
			if err != nil {
				return err
			}

			t.Log("SQL Server IP: " + sqlIp)

			octopusContainer, err := setupOctopus(ctx, "Server="+sqlIp+",1433;Database=OctopusDeploy;User=sa;Password=Password01!")
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
			success := false
			for start := time.Now(); ; {
				if time.Since(start) > 5*time.Minute {
					break
				}

				resp, err := http.Get(octopusContainer.URI + "/api")
				if err == nil && resp.StatusCode == http.StatusOK {
					success = true
					t.Log("Successfully contacted the Octopus API")
					break
				}

				time.Sleep(10 * time.Second)
			}

			if !success {
				t.Fatalf("Failed to access the Octopus API")
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

// initialiseOctopus uses Terraform to populate the test Octopus instance, making sure to clean up
// any files generated during previous Terraform executions to avoid conflicts and locking issues.
func initialiseOctopus(t *testing.T, container *OctopusContainer, terraformDir string, spaceName string, initialiseVars []string, populateVars []string) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	t.Log("Working dir: " + path)

	// This test creates a new space and then populates the space.
	terraformProjectDirs := []string{}
	terraformProjectDirs = append(terraformProjectDirs, filepath.Join("..", "..", "..", "test", "terraform", "1-singlespace"))
	terraformProjectDirs = append(terraformProjectDirs, filepath.Join(terraformDir))

	// First loop initialises the new space, second populates the space
	spaceId := "Spaces-1"
	for i, terraformProjectDir := range terraformProjectDirs {

		os.Remove(filepath.Join(terraformProjectDir, ".terraform.lock.hcl"))
		os.Remove(filepath.Join(terraformProjectDir, "terraform.tfstate"))

		args := []string{"init", "-no-color"}
		cmnd := exec.Command("terraform", args...)
		cmnd.Dir = terraformProjectDir
		out, err := cmnd.Output()

		if err != nil {
			exitError, ok := err.(*exec.ExitError)
			if ok {
				t.Log(string(exitError.Stderr))
			} else {
				t.Log(err.Error())
			}

			return err
		}

		t.Log(string(out))

		// when initialising the new space, we need to define a new space name as a variable
		vars := []string{}
		if i == 0 {
			vars = append(initialiseVars, "-var=octopus_space_name="+spaceName)
		} else {
			vars = populateVars
		}

		newArgs := append([]string{
			"apply",
			"-auto-approve",
			"-no-color",
			"-var=octopus_server=" + container.URI,
			"-var=octopus_apikey=" + ApiKey,
			"-var=octopus_space_id=" + spaceId,
		}, vars...)

		cmnd = exec.Command("terraform", newArgs...)
		cmnd.Dir = terraformProjectDir
		out, err = cmnd.Output()

		if err != nil {
			exitError, ok := err.(*exec.ExitError)
			if ok {
				t.Log(string(exitError.Stderr))
			} else {
				t.Log(err)
			}
			return err
		}

		t.Log(string(out))

		// get the ID of any new space created, which will be used in the subsequent Terraform executions
		spaceId, err = GetOutputVariable(t, terraformProjectDir, "octopus_space_id")

		if err != nil {
			exitError, ok := err.(*exec.ExitError)
			if ok {
				t.Log(string(exitError.Stderr))
			} else {
				t.Log(err)
			}
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
func Act(t *testing.T, container *OctopusContainer, terraformDir string, populateVars []string) (string, error) {
	t.Log("POPULATING TEST SPACE")

	spaceName := strings.ReplaceAll(fmt.Sprint(uuid.New()), "-", "")[:20]
	err := initialiseOctopus(t, container, terraformDir, spaceName, []string{}, populateVars)

	if err != nil {
		return "", err
	}

	return GetOutputVariable(t, filepath.Join("..", "..", "..", "test", "terraform", "1-singlespace"), "octopus_space_id")
}
