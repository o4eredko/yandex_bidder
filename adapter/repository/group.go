package repository

import (
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/sql"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
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

func (r *groupRepo) GetAll() ([]*domain.Group, error) {
	groups := make([]*domain.Group, 0)
	rows, err := r.store.DB.
		Select("id", "name", "schedule_start", "schedule_interval", "strategy").
		From("groups").
		Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var group *domain.Group
		if err := rows.ScanStruct(&group); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}
