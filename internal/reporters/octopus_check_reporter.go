package reporters

import "github.com/mcasperson/OctopusRecommendationEngine/internal/checks"

type OctopusCheckReporter interface {
	Generate(results []checks.OctopusCheckResult) (string, error)
}
