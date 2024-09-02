package php

import (
	"encoding/json"
	"testing"
)

func TestOutDate_ParseOutDate(t *testing.T) {
	testJson := "{\n      \"name\": \"diontech/laravel-extended-scheduler\",\n      \"direct-dependency\": false,\n      \"homepage\": null,\n      \"source\": \"https://github.com/DionTech/laravel-extended-scheduler/tree/v1.2.0\",\n      \"version\": \"v1.2.0\",\n      \"release-age\": \"3 years old\",\n      \"release-date\": \"2021-07-01T18:19:20+00:00\",\n      \"latest\": \"v1.3.0\",\n      \"latest-status\": \"semver-safe-update\",\n      \"latest-release-date\": \"2022-02-25T13:21:52+00:00\",\n      \"description\": \"This package allows you to configure the scheduled tasks of the app via (database) model. It was developed to avoid handling theseconfigurations via a config file only, cause then we cannot share the same repo to n server instances when running different tasks is needed at each server.\",\n      \"abandoned\": false\n    }"
	outDate := OutDate{}

	err := json.Unmarshal([]byte(testJson), &outDate)
	if err != nil {
		t.Fatalf("json.Unmarshal() failed with %s", err)
	}

	outDateParsed, err := outDate.ParseOutDate()
	if err != nil {
		t.Fatalf("outDate.ParseOutDate() failed with %s", err)
	}

	if outDateParsed.Name != "diontech/laravel-extended-scheduler" {
		t.Errorf("expected Name to be 'diontech/laravel-extended-scheduler', got %v", outDateParsed.Name)
	}

	if outDateParsed.CurrentVersion.Major != 1 {
		t.Errorf("expected CurrentVersion.Major to be 1, got %v", outDateParsed.CurrentVersion.Major)
	}

	if outDateParsed.CurrentVersion.Minor != 2 {
		t.Errorf("expected CurrentVersion.Minor to be 2, got %v", outDateParsed.CurrentVersion.Minor)
	}

	if outDateParsed.CurrentVersion.Patch != 0 {
		t.Errorf("expected CurrentVersion.Patch to be 0, got %v", outDateParsed.CurrentVersion.Patch)
	}

	if outDateParsed.LatestVersion.Major != 1 {
		t.Errorf("expected LatestVersion.Major to be 1, got %v", outDateParsed.LatestVersion.Major)
	}

	if outDateParsed.LatestVersion.Minor != 3 {
		t.Errorf("expected LatestVersion.Minor to be 3, got %v", outDateParsed.LatestVersion.Minor)
	}

	if outDateParsed.LatestVersion.Patch != 0 {
		t.Errorf("expected LatestVersion.Patch to be 0, got %v", outDateParsed.LatestVersion.Patch)
	}

}
