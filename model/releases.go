package model

import "time"

type Release struct {
	Provider Provider
	Project  Project

	Name string
	Date time.Time
	URL  string
}
