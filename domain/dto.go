package domain

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	GroupUpdateIn struct {
		ScheduleStart    *time.Time `json:"start"`
		ScheduleInterval int        `json:"interval"`
		StrategyName     string     `json:"strategy"`
	}

	AccountWithCampaigns struct {
		*Account
		Campaigns []*Campaign
	}

	BidsOut struct {
		AccountName string `json:"account_name"`
		Bids        []*Bid
		MaxRetries  int `json:"max_retries"`
	}
)

func (g *GroupUpdateIn) Validate() error {
	return validation.ValidateStruct(
		g,
		validation.Field(
			&g.ScheduleStart,
			validation.Required,
			validation.Min(time.Now().UTC()),
		),
		validation.Field(
			&g.ScheduleInterval,
			validation.Required,
		),
		validation.Field(
			&g.StrategyName,
			validation.Required,
		),
	)
}
