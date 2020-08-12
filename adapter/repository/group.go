package repository

import (
	"github.com/rs/zerolog/log"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/sql"
)

type (
	groupRepository struct {
		store *sql.Store
	}

	GroupRepository interface {
	}
)

func NewGroupRepository(s *sql.Store) GroupRepository {
	return &groupRepository{
		store: s,
	}
}

func (r *groupRepository) GetAll() []*domain.Group {
	sql := "SELECT id, name, schedule_start, schedule_interval, strategy FROM groups"
	var groups []*domain.Group

	query := r.store.DB.NewQuery(sql)
	if err := query.All(groups); err != nil {
		log.Error().Err(err)
	}

	return groups
}
