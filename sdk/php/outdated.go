package php

import "time"

type OutDate struct {
	Name              string    `json:"name"`
	DirectDependency  bool      `json:"direct-dependency"`
	Homepage          string    `json:"homepage"`
	Source            string    `json:"source"`
	Version           string    `json:"version"`
	ReleaseAge        string    `json:"release-age"`
	ReleaseDate       time.Time `json:"release-date"`
	Latest            string    `json:"latest"`
	LatestStatus      string    `json:"latest-status"`
	LatestReleaseDate time.Time `json:"latest-release-date"`
	Description       string    `json:"description"`
	Abandoned         bool      `json:"abandoned"`
}
