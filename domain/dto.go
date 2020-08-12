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
)

func (g *GroupUpdateIn) Validate() error {
	return validation.ValidateStruct(
		g,
		validation.Field(
			&g.ScheduleStart,
			validation.Required,
			validation.Date("1970-12-31T01:23:45Z"),
			validation.Min(time.Now()),
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
