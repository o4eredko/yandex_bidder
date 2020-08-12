package repository

import (
	dbx "github.com/go-ozzo/ozzo-dbx"

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
	var groups []*domain.Group
	query := r.store.DB.
		Select("id, name, schedule_start, schedule_interval, strategy").
		From("groups")

	if err := query.All(groups); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *groupRepo) GetByID(id int) (*domain.Group, error) {
	group := new(domain.Group)
	query := r.store.DB.NewQuery(
		`SELECT id, name, schedule_start, schedule_interval, strategy
 		 FROM groups
		 WHERE id = {:id}`,
	).Bind(dbx.Params{"id": id})

	if err := query.One(group); err != nil {

	}
}

// func (r *groupRepo) Update(group *domain.GroupUpdateIn) error {
// params := dbx.Params{
// 	"schedule_start":    group.Start,
// 	"schedule_interval": group.Interval,
// 	"strategy":          group.Strategy,
// }
// where := dbx.NewExp("id = {:id}", dbx.Params{"id": group.ID})

// query := r.store.DB.Update("groups", params, where)
// if _, err := query.Execute(); err != nil {
// 	return err
// }

// return r.store.DB.Model(group).Exclude("name").Update()

// return nil
// }
