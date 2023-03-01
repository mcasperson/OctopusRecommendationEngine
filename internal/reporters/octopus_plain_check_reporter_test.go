package reporters

import (
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"testing"
)

func TestNoChecks(t *testing.T) {
	results, err := OctopusPlainCheckReporter{}.Generate(nil)

	if err != nil {
		t.Fatal("Should not have returned an error")
	}

	if results != "" {
		t.Fatal("Should not have returned any results")
	}
}

func TestFailChecks(t *testing.T) {
	failedResult := checks.NewOctopusCheckResultImpl("This check always fails", "OctoRecAlwaysFail", "", checks.Error, "")
	results, err := OctopusPlainCheckReporter{minSeverity: checks.Error}.Generate([]checks.OctopusCheckResult{failedResult})

	if err != nil {
		t.Fatal("Should not have returned an error")
	}

	if results != "[OctoRecAlwaysFail] This check always fails" {
		t.Fatal("Should have returned 1 results")
	}
}

func TestFailAndPassChecks(t *testing.T) {
	failedResult := checks.NewOctopusCheckResultImpl("This check always fails", "OctoRecAlwaysFail", "", checks.Error, "")
	passResult := checks.NewOctopusCheckResultImpl("This check always passes", "OctoRecAlwaysPass", "", checks.Ok, "")
	results, err := OctopusPlainCheckReporter{minSeverity: checks.Error}.Generate([]checks.OctopusCheckResult{failedResult, passResult})

	if err != nil {
		t.Fatal("Should not have returned an error")
	}

	if results != "[OctoRecAlwaysFail] This check always fails" {
		t.Fatal("Should have returned 1 results")
	}
}

func TestFailAndPassWithOkChecks(t *testing.T) {
	failedResult := checks.NewOctopusCheckResultImpl("This check always fails", "OctoRecAlwaysFail", "", checks.Error, "")
	passResult := checks.NewOctopusCheckResultImpl("This check always passes", "OctoRecAlwaysPass", "", checks.Ok, "")
	results, err := OctopusPlainCheckReporter{minSeverity: checks.Ok}.Generate([]checks.OctopusCheckResult{failedResult, passResult})

	if err != nil {
		t.Fatal("Should not have returned an error")
	}

	if results != "[OctoRecAlwaysFail] This check always fails\n\n[OctoRecAlwaysPass] This check always passes" {
		t.Fatal("Should have returned 1 results")
	}
}
