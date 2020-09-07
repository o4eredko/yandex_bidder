package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	Group struct {
		ID          int        `json:"id"`
		Name        string     `json:"name"`
		Strategy    *string    `json:"strategy"`
		Start       *time.Time `json:"start" db:"schedule_start"`
		Interval    *int       `json:"interval" db:"schedule_interval"`
		IsScheduled bool       `json:"is_scheduled"`
	}
)

func (g *Group) MaxRetries() int {
	if g.Interval == nil {
		return 0
	}
	return *g.Interval / 2
}

func (g *Group) ValidateSchedule() error {
	return validation.ValidateStruct(
		g,
		validation.Field(
			&g.ID,
			validation.Required,
		),
		validation.Field(
			&g.Interval,
			validation.Required,
		),
		validation.Field(
			&g.Start,
			validation.Required,
		),
		validation.Field(
			&g.Strategy,
			validation.Required,
		),
	)
}
