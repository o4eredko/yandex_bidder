package repository

import (
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

// func (r *groupRepository) GetAll() []*domain.Group {
// 	query := "SELECT id, name, schedule_start, schedule_interval FROM groups"
// 	var groups []*domain.Group
// }
