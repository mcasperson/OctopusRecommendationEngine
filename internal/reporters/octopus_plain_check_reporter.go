package reporters

import (
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"strings"
)

type OctopusPlainCheckReporter struct {
	minSeverity int
}

func (o OctopusPlainCheckReporter) Generate(results []checks.OctopusCheckResult) (string, error) {
	if results == nil || len(results) == 0 {
		return "", nil
	}

	report := []string{}

	for _, r := range results {
		if r.Severity() >= o.minSeverity {
			report = append(report, r.Code()+": "+r.Description())
		}
	}

	return strings.Join(report[:], "\n"), nil
}
