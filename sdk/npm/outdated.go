package npm

import (
	"fmt"
	"strconv"
	"strings"
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
	Name           string         `json:"name"`
	Current        string         `json:"current"`
	CurrentVersion *Version       `json:"currentVersion"`
	Wanted         string         `json:"wanted"`
	WantedVersion  *Version       `json:"wantedVersion"`
	Latest         string         `json:"latest"`
	LatestVersion  *Version       `json:"latestVersion"`
	Dependent      string         `json:"dependent"`
	Rating         *OutDateRating `json:"rating"`
}

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

type OutDateRating struct {
	StatusCode int
	Message    string
}

func NewOutDate(name, current, wanted, latest, dependent string) (*OutDate, error) {
	currentVersion, err := ConvertVersion(current)
	if err != nil {
		return nil, err
	}

	wantedVersion, err := ConvertVersion(wanted)
	if err != nil {
		return nil, err
	}

	latestVersion, err := ConvertVersion(latest)
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
		o.Rating.Message = fmt.Sprintf("%s %s | have %s", OutDateRatingMessageNewMajorVersionAvailable, o.Latest, o.Current)
	} else if o.NewMinorVersionAvailable() {
		o.Rating.StatusCode = OutDateRatingNewMinorVersionAvailable
		o.Rating.Message = fmt.Sprintf("%s %s | have %s", OutDateRatingMessageNewMinorVersionAvailable, o.Latest, o.Current)
	} else if o.NewPatchVersionAvailable() {
		o.Rating.StatusCode = OutDateRatingNewPatchVersionAvailable
		o.Rating.Message = fmt.Sprintf("%s %s | have %s", OutDateRatingMessageNewPatchVersionAvailable, o.Latest, o.Current)
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

func ConvertVersion(version string) (*Version, error) {
	//empty string -> means no version
	if version == "" {
		return &Version{
			Major: 0,
			Minor: 0,
			Patch: 0,
		}, nil
	}

	//remove all after a -, for example we have -rc2
	version = strings.Split(version, "-")[0]

	versionNumbers := strings.Split(version, ".")
	if len(versionNumbers) == 0 {
		return &Version{}, fmt.Errorf("invalid version syntax: %s", version)
	}

	if len(versionNumbers) > 3 {
		return &Version{}, fmt.Errorf("invalid version syntax: %s", version)
	}

	major := 0
	minor := 0
	patch := 0

	for i, v := range versionNumbers {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return &Version{}, fmt.Errorf("invalid version syntax: %v", v)
		}

		switch i {
		case 0:
			major = vInt
		case 1:
			minor = vInt
		case 2:
			patch = vInt
		}
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}
