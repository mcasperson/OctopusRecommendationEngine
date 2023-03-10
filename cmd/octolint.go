package main

import (
	"flag"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks/factory"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/executor"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/reporters"
	"github.com/mcasperson/OctopusTerraformTestFramework/octoclient"
	"os"
	"time"
)

func main() {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	url, space, apiKey := parseUrl()
	client, err := octoclient.CreateClient(url, space, apiKey)

	if err != nil {
		errorExit("Failed to create the Octopus client")
	}

	factory := factory.NewOctopusCheckFactory(client)
	checkCollection, err := factory.BuildAllChecks()

	if err != nil {
		errorExit("Failed to create the checks")
	}

	executor := executor.NewOctopusCheckExecutor()
	results, err := executor.ExecuteChecks(checkCollection, func(check checks.OctopusCheck, err error) error {
		fmt.Println("Failed to execute check " + check.Id() + ": " + err.Error())
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

func parseUrl() (string, string, string) {
	var url string
	flag.StringVar(&url, "url", "", "The Octopus URL e.g. https://myinstance.octopus.app")

	var space string
	flag.StringVar(&space, "space", "", "The Octopus space name or ID")

	var apiKey string
	flag.StringVar(&apiKey, "apiKey", "", "The Octopus api key")

	flag.Parse()

	return url, space, apiKey
}
