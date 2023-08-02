package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/briandowns/spinner"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks/factory"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/executor"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/reporters"
	"github.com/mcasperson/OctopusTerraformTestFramework/octoclient"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	url, space, apiKey, skipTests, verboseErrors := parseArgs()

	if url == "" {
		errorExit("You must specify the URL with the -url argument")
	}

	if apiKey == "" {
		errorExit("You must specify the API key with the -apiKey argument")
	}

	if !strings.HasPrefix(space, "Spaces-") {
		spaceId, err := lookupSpaceAsName(url, space, apiKey)

		if err != nil {
			errorExit("Failed to create the Octopus client")
		}

		space = spaceId
	}

	client, err := octoclient.CreateClient(url, space, apiKey)

	if err != nil {
		errorExit("Failed to create the Octopus client. Check that the url, api key, and space are correct.")
	}

	factory := factory.NewOctopusCheckFactory(client)
	checkCollection, err := factory.BuildAllChecks(skipTests)

	if err != nil {
		errorExit("Failed to create the checks")
	}

	executor := executor.NewOctopusCheckExecutor()
	results, err := executor.ExecuteChecks(checkCollection, func(check checks.OctopusCheck, err error) error {
		fmt.Fprintf(os.Stderr, "Failed to execute check "+check.Id())
		if verboseErrors {
			fmt.Fprintf(os.Stdout, "##octopus[stdout-verbose]")
		}
		fmt.Fprintf(os.Stderr, err.Error())
		if verboseErrors {
			fmt.Fprintf(os.Stdout, "##octopus[stdout-default]")
		}
		return nil
	})

	if err != nil {
		errorExit("Failed to run the checks")
	}

	reporter := reporters.NewOctopusPlainCheckReporter(checks.Warning)
	report, err := reporter.Generate(results)

	if err != nil {
		errorExit("Failed to generate the report")
	}

	s.Stop()
	fmt.Println(report)
}

func errorExit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func parseArgs() (string, string, string, string, bool) {
	var url string
	flag.StringVar(&url, "url", "", "The Octopus URL e.g. https://myinstance.octopus.app")

	var space string
	flag.StringVar(&space, "space", "", "The Octopus space name or ID")

	var apiKey string
	flag.StringVar(&apiKey, "apiKey", "", "The Octopus api key")

	var skipTests string
	flag.StringVar(&skipTests, "skipTests", "", "A comma separated list of tests to skip")

	var verboseErrors bool
	flag.BoolVar(&verboseErrors, "verboseErrors", false, "Print error details as verbose logs inOctopus")

	flag.Parse()

	if url == "" {
		url = os.Getenv("OCTOPUS_CLI_SERVER")
	}

	if apiKey == "" {
		apiKey = os.Getenv("OCTOPUS_CLI_API_KEY")
	}

	return url, space, apiKey, skipTests, verboseErrors
}

func lookupSpaceAsName(octopusUrl string, spaceName string, apiKey string) (string, error) {
	if len(strings.TrimSpace(spaceName)) == 0 {
		return "", errors.New("space can not be empty")
	}

	requestURL := fmt.Sprintf("%s/api/Spaces?take=1000&partialName=%s", octopusUrl, url.QueryEscape(spaceName))

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)

	if err != nil {
		return "", err
	}

	if apiKey != "" {
		req.Header.Set("X-Octopus-ApiKey", apiKey)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", nil
	}
	defer res.Body.Close()

	collection := resources.Resources[spaces.Space]{}
	err = json.NewDecoder(res.Body).Decode(&collection)

	if err != nil {
		return "", err
	}

	for _, space := range collection.Items {
		if space.Name == spaceName {
			return space.ID, nil
		}
	}

	return "", errors.New("did not find space with name " + spaceName)
}
