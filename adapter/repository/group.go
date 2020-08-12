package repository

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/rs/zerolog/log"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/sql"
)

type (
	groupRepo struct {
		store *sql.Store
	}
)

func NewGroupRepository(s *sql.Store) usecase.GroupRepo {
	return &groupRepo{
		store: s,
	}
}

func (r *groupRepo) GetAll() []*domain.Group {
	var groups []*domain.Group
	query := r.store.DB.
		Select("id, name, schedule_start, schedule_interval, strategy").
		From("groups")

	if err := query.All(groups); err != nil {
		log.Error().Err(err)
	}

	return groups
}
