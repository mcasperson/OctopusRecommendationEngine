package security

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"regexp"
	"strings"
	"time"
)

type APIKeyKey struct {
	Hint *string
}

// APIKey is used because the go client has an invalid APIKey value that prevents the usual functions for querying users keys
type APIKey struct {
	APIKey  APIKeyKey  `json:"ApiKey,omitempty"`
	Expires *time.Time `json:"Expires,omitempty"`
}

// OctopusPerpetualApiKeysCheck reports on any perpetual api keys
type OctopusPerpetualApiKeysCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
}

func NewOctopusPerpetualApiKeysCheck(client *client.Client, errorHandler checks.OctopusClientErrorHandler) OctopusPerpetualApiKeysCheck {
	return OctopusPerpetualApiKeysCheck{client: client, errorHandler: errorHandler}
}

func (o OctopusPerpetualApiKeysCheck) Id() string {
	return "OctoRecPerpetualApiKeys"
}

func (o OctopusPerpetualApiKeysCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	users, err := o.client.Users.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}

	linksTemplate := regexp.MustCompile(`\{.+\}`)
	perpetualApiKeys := []string{}
	for _, u := range users {

		apiKeysLink := linksTemplate.ReplaceAllString(u.Links["ApiKeys"], "")
		keys, err := newclient.Get[resources.Resources[APIKey]](o.client.HttpSession(), apiKeysLink)

		if err != nil {
			if !o.errorHandler.ShouldContinue(err) {
				return nil, err
			}
			continue
		}

		for _, k := range keys.Items {
			if k.Expires == nil && k.APIKey.Hint != nil {
				perpetualApiKeys = append(perpetualApiKeys, *k.APIKey.Hint+"... ("+u.Username+")")
			}
		}
	}

	if len(perpetualApiKeys) != 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following API keys do not expire: "+strings.Join(perpetualApiKeys, ", "),
			o.Id(),
			"",
			checks.Warning,
			checks.Security), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"No perpetual API keys found",
		o.Id(),
		"",
		checks.Ok,
		checks.Security), nil
}
