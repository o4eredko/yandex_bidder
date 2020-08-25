package domain

import (
	"time"
)

type (
	Campaign struct {
		ID int
	}

	Account struct {
		ID   int
		Name string
	}

	Strategy struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	Group struct {
		ID       int        `json:"id"`
		Name     string     `json:"name"`
		Strategy *string    `json:"strategy"`
		Start    *time.Time `json:"start" db:"schedule_start"`
		Interval *int       `json:"interval" db:"schedule_interval"`
	}

	Bid struct {
		CampaignID int
		Bid        int
	}
)

func (g *Group) CalculateMaxRetries() int {
	return *g.Interval / 2
}
