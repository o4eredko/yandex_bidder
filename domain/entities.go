package domain

import (
	"time"
)

type (
	Campaign struct {
		ID    int
		Stats Stats
	}

	Schedule struct {
		Start    time.Time
		Interval int
	}

	BidHandler func(int, int, int) int

	Strategy struct {
		Name    string
		Handler BidHandler
	}

	Group struct {
		ID       int        `json:"id"`
		Name     string     `json:"name"`
		Strategy *string    `json:"strategy"`
		Start    *time.Time `json:"start" db:"schedule_start"`
		Interval *int       `json:"interval" db:"schedule_interval"`
	}

	Strategies struct {
		strategies map[string]BidHandler
	}

	Stats struct {
		AccountName string `db:"account_name"`
		CampaignId  int    `db:"campaign_id"`
		Impressions int
		Cost        int
		Clicks      int
	}
)
