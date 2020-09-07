package domain

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
)

type (
	GroupUpdateIn struct {
		ID       int
		Start    *time.Time `json:"start"`
		Interval int        `json:"interval"`
		Strategy string     `json:"strategy"`
	}

	AccountWithCampaigns struct {
		*entities.Account
		Campaigns []*entities.Campaign
	}

	AccountBids struct {
		AccountName string          `json:"name"`
		Bids        []*entities.Bid `json:"bids"`
	}

	GroupToUpdateBids struct {
		Name       string         `json:"name"`
		Accounts   []*AccountBids `json:"accounts"`
		MaxRetries int            `json:"max_retries"`
	}

	GroupToggleIn struct {
		Action string `json:"action"`
	}
)

func (g *GroupUpdateIn) Validate() error {
	return validation.ValidateStruct(
		g,
		validation.Field(
			&g.ID,
			validation.Required,
		),
		validation.Field(
			&g.Start,
			validation.Required,
			validation.Min(time.Now().UTC()),
		),
		validation.Field(
			&g.Interval,
			validation.Required,
		),
		validation.Field(
			&g.Strategy,
			validation.Required,
		),
	)
}

func (g *GroupToggleIn) Validate() error {
	return validation.ValidateStruct(
		g,
		validation.Field(
			&g.Action,
			validation.Required,
			validation.In("start", "pause"),
		),
	)
}
