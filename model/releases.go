package model

import "time"

type Release struct {
	Name string
	Date time.Time
	URL  string
}

type Tag struct {
	Name    string
	Date    time.Time
	URL     string
	Message string
}
