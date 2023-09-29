package code

import (
	"strings"
	"testing"
)

func TestAnalyser_SetCoverProfile(t *testing.T) {
	cp := "?   \tdatacatalogapi/sdk/awslog\t[no test files]\n?   \tdatacatalogapi/sdk/csv\t[no test files]\n?   \tdatacatalogapi/sdk/datetime\t[no test files]\n?   \tdatacatalogapi/sdk/debugbar\t[no test files]\n?   \tdatacatalogapi/sdk/email\t[no test files]\n?   \tdatacatalogapi/sdk/eventpusher\t[no test files]\n?   \tdatacatalogapi/sdk/fileupload\t[no test files]\nok  \tdatacatalogapi/sdk/accesscontroltests\t6.356s\tcoverage: 0.8% of statements in ./sdk/...\n?   \tdatacatalogapi/sdk/permissions\t[no test files]\nok  \tdatacatalogapi/sdk/activities\t14.238s\tcoverage: 1.4% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/automations\t2.462s\tcoverage: 11.4% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/awsutil\t0.507s\tcoverage: 13.2% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/business\t0.059s\tcoverage: 0.3% of statements in ./sdk/... [no tests to run]\nok  \tdatacatalogapi/sdk/cache\t10.739s\tcoverage: 24.0% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/cluster\t160.247s\tcoverage: 37.5% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/clustercalculator\t0.024s\tcoverage: 0.0% of statements in ./sdk/... [no tests to run]\nok  \tdatacatalogapi/sdk/diff\t0.007s\tcoverage: 33.3% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/dsltests\t0.764s\tcoverage: 1.4% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/export\t8.097s\tcoverage: 15.9% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/fieldtypes\t0.693s\tcoverage: 4.1% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/gfdashboard\t21.251s\tcoverage: 73.6% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/identity\t12.152s\tcoverage: 43.5% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/jsapi\t0.046s\tcoverage: 5.4% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/msd\t29.701s\tcoverage: 10.7% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/notifications\t1.323s\tcoverage: 22.5% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/notificationshub\t26.611s\tcoverage: 41.3% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/openaiutil\t0.003s\tcoverage: 0.0% of statements in ./sdk/... [no tests to run]\nok  \tdatacatalogapi/sdk/pseudomizer\t0.004s\tcoverage: 67.4% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/realtimenotifications\t0.137s\tcoverage: 4.3% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/schemadefinition\t2.452s\tcoverage: 13.9% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/sitetracker\t3.931s\tcoverage: 3.1% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/statusnotifications\t79.496s\tcoverage: 14.9% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/task\t49.584s\tcoverage: 21.3% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/templateengine\t0.453s\tcoverage: 8.6% of statements in ./sdk/...\nok  \tdatacatalogapi/sdk/webom\t2.689s\tcoverage: 6.9% of statements in ./sdk/..."
	lines := strings.Split(cp, "\n")

	a := &Analyser{}

	a.SetCoverProfile(cp)
	if a.CoverProfile == "" {
		t.Errorf("Error: coverProfile is empty\n")
	}

	if len(a.profile) == 0 {
		t.Errorf("Error: profile is empty\n")
	}

	if len(a.profile) != len(lines) {
		t.Errorf("Error: profile is not correct\n")
	}

	//test the result of datacatalogapi/sdk/activities
	testProfile := a.profile[9]
	if testProfile.Folder != "datacatalogapi/sdk/activities" {
		t.Errorf("Error: folder is wrong: %s\n", testProfile.Folder)
	}

	if testProfile.Duration != 14.238 {
		t.Errorf("Error: duration is wrong: %f\n", testProfile.Duration)
	}

	if testProfile.Coverage != 1.4 {
		t.Errorf("Error: coverage is wrong: %f\n", testProfile.Coverage)
	}

	//test the first one with no coverage
	testProfile = a.profile[0]
	if testProfile.Folder != "datacatalogapi/sdk/awslog" {
		t.Errorf("Error: folder is wrong: %s\n", testProfile.Folder)
	}

	if testProfile.Duration != 0 {
		t.Errorf("Error: duration is wrong: %f\n", testProfile.Duration)
	}

	if testProfile.Coverage != 0 {
		t.Errorf("Error: coverage is wrong: %f\n", testProfile.Coverage)
	}
}
