package helper

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
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
