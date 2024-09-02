package php

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/helper"
	"time"
)

const (
	OutDateRatingNewMajorVersionAvailable         = 500
	OutDateRatingIndirectNewMajorVersionAvailable = 430
	OutDateRatingNewMinorVersionAvailable         = 400
	OutDateRatingNewPatchVersionAvailable         = 300
	OutDateRatingNoNewVersionAvailable            = 200
)

const (
	OutDateRatingMessageNewMajorVersionAvailable        = "New major version available"
	OutputRatingMessageIndirectNewMajorVersionAvailable = "Indirect new major version available"
	OutDateRatingMessageNewMinorVersionAvailable        = "New minor version available"
	OutDateRatingMessageNewPatchVersionAvailable        = "New patch version available"
)

type OutDate struct {
	Name              string          `json:"name"`
	DirectDependency  bool            `json:"direct-dependency"`
	Homepage          string          `json:"homepage"`
	Source            string          `json:"source"`
	CurrentVersion    *helper.Version `json:"currentVersion"`
	Version           string          `json:"version"`
	ReleaseAge        string          `json:"release-age"`
	ReleaseDate       time.Time       `json:"release-date"`
	Latest            string          `json:"latest"`
	LatestVersion     *helper.Version `json:"latestVersion"`
	LatestStatus      string          `json:"latest-status"`
	LatestReleaseDate time.Time       `json:"latest-release-date"`
	Description       string          `json:"description"`
	Abandoned         bool            `json:"abandoned"`
	Rating            *OutDateRating  `json:"rating"`
}

type OutDateRating struct {
	StatusCode int
	Message    string
}

func (o *OutDate) ParseOutDate() (*OutDate, error) {
	o.Rating = &OutDateRating{
		StatusCode: 0,
		Message:    "",
	}

	currentVersion, err := helper.ConvertVersion(o.Version)
	if err != nil {
		return o, err
	}

	o.CurrentVersion = currentVersion

	latestVersion, err := helper.ConvertVersion(o.Latest)
	if err != nil {
		return o, err
	}

	o.LatestVersion = latestVersion

	return o, nil
}

func (o *OutDate) Rate() {
	if o.NewMajorVersionAvailable() {
		if o.DirectDependency {
			o.Rating.StatusCode = OutDateRatingNewMajorVersionAvailable
			o.Rating.Message = fmt.Sprintf("%s %s | have %s", OutDateRatingMessageNewMajorVersionAvailable, o.Latest, o.Version)
		} else {
			o.Rating.StatusCode = OutDateRatingIndirectNewMajorVersionAvailable
			o.Rating.Message = fmt.Sprintf("%s %s | have %s", OutputRatingMessageIndirectNewMajorVersionAvailable, o.Latest, o.Version)
		}
	} else if o.NewMinorVersionAvailable() {
		o.Rating.StatusCode = OutDateRatingNewMinorVersionAvailable
		o.Rating.Message = fmt.Sprintf("%s %s | have %s", OutDateRatingMessageNewMinorVersionAvailable, o.Latest, o.Version)
	} else if o.NewPatchVersionAvailable() {
		o.Rating.StatusCode = OutDateRatingNewPatchVersionAvailable
		o.Rating.Message = fmt.Sprintf("%s %s | have %s", OutDateRatingMessageNewPatchVersionAvailable, o.Latest, o.Version)
	} else {
		o.Rating.StatusCode = OutDateRatingNoNewVersionAvailable
	}
}

func (o *OutDate) NewMajorVersionAvailable() bool {
	return o.CurrentVersion.Major < o.LatestVersion.Major
}

func (o *OutDate) NewMinorVersionAvailable() bool {
	return o.CurrentVersion.Major == o.LatestVersion.Major &&
		o.CurrentVersion.Minor < o.LatestVersion.Minor
}

func (o *OutDate) NewPatchVersionAvailable() bool {
	return o.CurrentVersion.Major == o.LatestVersion.Major &&
		o.CurrentVersion.Minor == o.LatestVersion.Minor &&
		o.CurrentVersion.Patch < o.LatestVersion.Patch
}
