package repository

import (
	"database/sql"

	dbx "github.com/go-ozzo/ozzo-dbx"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	groupRepo struct {
		db dbx.Builder
	}
)

func NewGroupRepo(db dbx.Builder) usecase.GroupRepo {
	return &groupRepo{
		db: db,
	}
}

func (r *groupRepo) GetAll() ([]*domain.Group, error) {
	groups := make([]*domain.Group, 0)
	rows, err := r.db.
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

func (r *groupRepo) GetByID(id int) (*domain.Group, error) {
	group := new(domain.Group)
	query := r.db.Select("id", "name", "schedule_start", "schedule_interval", "strategy").
		From("groups").
		Where(dbx.NewExp("id", dbx.Params{"id": id}))

	if err := query.One(group); err != nil {
		if err == sql.ErrNoRows {
			err = domain.ErrGroupNotFound
		}
		return nil, err
	}

	return group, nil
}

func (r *groupRepo) Update(group *domain.GroupUpdateIn) error {
	// params := dbx.Params{
	// 	"schedule_start":    group.Start,
	// 	"schedule_interval": group.Interval,
	// 	"strategy":          group.Strategy,
	// }
	// where := dbx.NewExp("id = {:id}", dbx.Params{"id": group.ID})
	//
	// query := r.store.DB.Update("groups", params, where)
	// if _, err := query.Execute(); err != nil {
	// 	return err
	// }

	// return r.store.DB.Model(group).Exclude("name").Update()

	return nil
}
