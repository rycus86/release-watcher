package model

import "time"

type Release struct {
	Provider Provider
	Project  GenericProject

	Name string
	Date time.Time
	URL  string
}
