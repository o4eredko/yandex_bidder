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

func (r *groupRepo) Accounts(group *domain.Group) ([]*domain.Account, error) {
	accounts := make([]*domain.Account, 0)
	rows, err := r.db.
		Select("id", "name").
		From("accounts AS a").
		Where(dbx.HashExp{"group_id": group.ID}).
		Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		account := new(domain.Account)
		if err := rows.ScanStruct(account); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
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
	query := r.db.
		Select("g.id", "g.name", "schedule_start", "schedule_interval", "s.name AS strategy").
		From("groups AS g").
		LeftJoin("strategies AS s", dbx.NewExp("g.strategy_id = s.id")).
		Where(dbx.HashExp{"g.id": id})

	if err := query.One(group); err != nil {
		if err == sql.ErrNoRows {
			err = domain.ErrGroupNotFound
		}
		return nil, err
	}

	return group, nil
}

func (r *groupRepo) Update(group *domain.Group) error {
	fieldsToUpdate := dbx.Params{
		"schedule_start":    group.Start,
		"schedule_interval": group.Interval,
		"strategy":          group.Strategy,
	}
	whereClause := dbx.HashExp{"id": group.ID}

	query := r.db.Update("groups", fieldsToUpdate, whereClause)
	if _, err := query.Execute(); err != nil {
		return err
	}

	return nil
}
