package npm

import (
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/helper"
)

const (
	OutDateRatingNewMajorVersionAvailable = 500
	OutDateRatingNewMinorVersionAvailable = 400
	OutDateRatingNewPatchVersionAvailable = 300
	OutDateRatingNoNewVersionAvailable    = 200
)

const (
	OutDateRatingMessageNewMajorVersionAvailable = "New major version available"
	OutDateRatingMessageNewMinorVersionAvailable = "New minor version available"
	OutDateRatingMessageNewPatchVersionAvailable = "New patch version available"
)

type OutDate struct {
	Name           string          `json:"name"`
	Current        string          `json:"current"`
	CurrentVersion *helper.Version `json:"currentVersion"`
	Wanted         string          `json:"wanted"`
	WantedVersion  *helper.Version `json:"wantedVersion"`
	Latest         string          `json:"latest"`
	LatestVersion  *helper.Version `json:"latestVersion"`
	Dependent      string          `json:"dependent"`
	Rating         *OutDateRating  `json:"rating"`
}

type OutDateRating struct {
	StatusCode int
	Message    string
}

func NewOutDate(name, current, wanted, latest, dependent string) (*OutDate, error) {
	currentVersion, err := helper.ConvertVersion(current)
	if err != nil {
		return nil, err
	}

	wantedVersion, err := helper.ConvertVersion(wanted)
	if err != nil {
		return nil, err
	}

	latestVersion, err := helper.ConvertVersion(latest)
	if err != nil {
		return nil, err
	}

	return &OutDate{
		Name:           name,
		Current:        current,
		CurrentVersion: currentVersion,
		Wanted:         wanted,
		WantedVersion:  wantedVersion,
		Latest:         latest,
		LatestVersion:  latestVersion,
		Dependent:      dependent,
		Rating:         &OutDateRating{StatusCode: 0, Message: ""},
	}, nil
}

func (o *OutDate) Rate() {
	if o.NewMajorVersionAvailable() {
		o.Rating.StatusCode = OutDateRatingNewMajorVersionAvailable
		o.Rating.Message = fmt.Sprintf("%s %s | have %s", OutDateRatingMessageNewMajorVersionAvailable, o.Latest, o.Wanted)
	} else if o.NewMinorVersionAvailable() {
		o.Rating.StatusCode = OutDateRatingNewMinorVersionAvailable
		o.Rating.Message = fmt.Sprintf("%s %s | have %s", OutDateRatingMessageNewMinorVersionAvailable, o.Latest, o.Wanted)
	} else if o.NewPatchVersionAvailable() {
		o.Rating.StatusCode = OutDateRatingNewPatchVersionAvailable
		o.Rating.Message = fmt.Sprintf("%s %s | have %s", OutDateRatingMessageNewPatchVersionAvailable, o.Latest, o.Wanted)
	} else {
		o.Rating.StatusCode = OutDateRatingNoNewVersionAvailable
	}
}

func (o *OutDate) NewMajorVersionAvailable() bool {
	return o.WantedVersion.Major < o.LatestVersion.Major
}

func (o *OutDate) NewMinorVersionAvailable() bool {
	return o.WantedVersion.Major == o.LatestVersion.Major &&
		o.WantedVersion.Minor < o.LatestVersion.Minor
}

func (o *OutDate) NewPatchVersionAvailable() bool {
	return o.WantedVersion.Major == o.LatestVersion.Major &&
		o.WantedVersion.Minor == o.LatestVersion.Minor &&
		o.WantedVersion.Patch < o.LatestVersion.Patch
}
