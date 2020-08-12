package domain

import (
	"time"
)

type (
	Stats struct {
		Clicks int
		// ...
	}

	Campaign struct {
		ID    int
		Stats Stats
	}

	Schedule struct {
		Start    time.Time
		Interval int
	}

	BidHandler func(int, int, int) (int, error)

	Strategy struct {
		Name     string
		Schedule Schedule
		Handler  BidHandler
	}

	Group struct {
		Name     string
		Strategy Strategy
	}

	Strategies struct {
		strategies map[string]BidHandler
	}
)
