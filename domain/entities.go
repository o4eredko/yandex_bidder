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
		Name    string
		Handler BidHandler
	}

	Group struct {
		ID       int        `json:"id"`
		Name     string     `json:"name"`
		Strategy *Strategy  `json:"strategy"`
		Start    *time.Time `json:"start"`
		Interval *int       `json:"interval"`
	}

	Strategies struct {
		strategies map[string]BidHandler
	}
)
