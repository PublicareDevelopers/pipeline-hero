package npm

import (
	"fmt"
	"strconv"
	"strings"
)

type OutDate struct {
	Current        string   `json:"current"`
	CurrentVersion *Version `json:"currentVersion"`
	Wanted         string   `json:"wanted"`
	WantedVersion  *Version `json:"wantedVersion"`
	Latest         string   `json:"latest"`
	LatestVersion  *Version `json:"latestVersion"`
	Dependent      string   `json:"dependent"`
}

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

func NewOutDate(current, wanted, latest, dependent string) (*OutDate, error) {
	return &OutDate{
		Current:   current,
		Wanted:    wanted,
		Latest:    latest,
		Dependent: dependent,
	}, nil
}

func (o *OutDate) NewMajorVersionAvailable() bool {
	return false
}

func ConvertVersion(version string) (*Version, error) {
	versionNumbers := strings.Split(version, ".")
	if len(versionNumbers) == 0 {
		return &Version{}, fmt.Errorf("invalid version syntax")
	}

	if len(versionNumbers) > 3 {
		return &Version{}, fmt.Errorf("invalid version syntax")
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
